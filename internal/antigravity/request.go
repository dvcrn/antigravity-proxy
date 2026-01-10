package antigravity

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/dvcrn/gemini-code-assist-proxy/internal/logger"
	"github.com/google/uuid"
)

func prepareAntigravityRequest(req *GenerateContentRequest) {
	if req == nil {
		return
	}

	req.UserAgent = RequestUserAgent
	req.RequestType = RequestTypeAgent
	if req.RequestID == "" {
		req.RequestID = "agent-" + uuid.NewString()
	}

	if req.Request.SessionID == "" {
		req.Request.SessionID = deriveSessionID(req.Request.Contents)
	}

	if missing := fillMissingParameters(req.Request.Tools); missing > 0 {
		logger.Get().Warn().
			Int("missing_parameters", missing).
			Str("missing_names", missingParameterNames(req.Request.Tools, 6)).
			Msg("Defaulted missing parameters in request tools")
	}

	if missing := ensureFunctionCallIDs(req.Request.Contents); missing > 0 {
		logger.Get().Warn().
			Int("missing_ids", missing).
			Msg("Defaulted missing functionCall IDs in request contents")
	}

	if missing := ensureFunctionResponseIDs(req.Request.Contents); missing > 0 {
		logger.Get().Warn().
			Int("missing_ids", missing).
			Msg("Defaulted missing functionResponse IDs in request contents")
	}

	req.Request.SystemInstruction = buildAntigravitySystemInstruction(req.Request.SystemInstruction)
}

func buildAntigravitySystemInstruction(existing *SystemInstruction) *SystemInstruction {
	parts := []ContentPart{
		{Text: SystemInstructionText},
		{Text: "Please ignore the following [ignore]" + SystemInstructionText + "[/ignore]"},
	}

	if existing != nil {
		for _, part := range existing.Parts {
			if part.Text != "" {
				parts = append(parts, ContentPart{Text: part.Text})
			}
		}
	}

	return &SystemInstruction{
		Role:  "user",
		Parts: parts,
	}
}

func deriveSessionID(contents []Content) string {
	for _, msg := range contents {
		if strings.ToLower(msg.Role) != "user" {
			continue
		}

		var parts []string
		for _, part := range msg.Parts {
			if part.Text != "" {
				parts = append(parts, part.Text)
			}
		}

		if len(parts) == 0 {
			continue
		}

		sum := sha256.Sum256([]byte(strings.Join(parts, "\n")))
		return hex.EncodeToString(sum[:])[:32]
	}

	return uuid.NewString()
}

func ensureFunctionCallIDs(contents []Content) int {
	missing := 0
	for contentIndex := range contents {
		for partIndex := range contents[contentIndex].Parts {
			part := contents[contentIndex].Parts[partIndex]
			if part.FunctionCall == nil {
				continue
			}
			if strings.TrimSpace(part.FunctionCall.ID) == "" {
				contents[contentIndex].Parts[partIndex].FunctionCall.ID = "toolu_" + uuid.NewString()
				missing++
			}
		}
	}
	return missing
}

func ensureFunctionResponseIDs(contents []Content) int {
	missing := 0
	var pending []string
	perName := map[string][]string{}

	for contentIndex := range contents {
		for partIndex := range contents[contentIndex].Parts {
			part := &contents[contentIndex].Parts[partIndex]
			if part.FunctionCall != nil {
				id := strings.TrimSpace(part.FunctionCall.ID)
				if id == "" {
					continue
				}
				pending = append(pending, id)
				if name := strings.TrimSpace(part.FunctionCall.Name); name != "" {
					perName[name] = append(perName[name], id)
				}
				continue
			}

			if part.FunctionResponse == nil {
				continue
			}

			if strings.TrimSpace(part.FunctionResponse.ID) != "" {
				// Consume any queued IDs so we don't reuse them later.
				pending = removeFirstMatch(pending, part.FunctionResponse.ID)
				if name := strings.TrimSpace(part.FunctionResponse.Name); name != "" {
					perName[name] = removeFirstMatch(perName[name], part.FunctionResponse.ID)
				}
				continue
			}

			assigned := ""
			if name := strings.TrimSpace(part.FunctionResponse.Name); name != "" {
				if ids := perName[name]; len(ids) > 0 {
					assigned = ids[0]
					perName[name] = ids[1:]
					pending = removeFirstMatch(pending, assigned)
				}
			}
			if assigned == "" && len(pending) > 0 {
				assigned = pending[0]
				pending = pending[1:]
				if name := strings.TrimSpace(part.FunctionResponse.Name); name != "" {
					perName[name] = removeFirstMatch(perName[name], assigned)
				}
			}

			if assigned != "" {
				part.FunctionResponse.ID = assigned
				missing++
			}
		}
	}

	return missing
}

func removeFirstMatch(values []string, target string) []string {
	if target == "" {
		return values
	}
	for i, v := range values {
		if v == target {
			return append(values[:i], values[i+1:]...)
		}
	}
	return values
}
