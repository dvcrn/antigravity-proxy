package antigravity

import (
	"fmt"
	"net/http"
	"runtime"
)

const (
	endpointDaily         = "https://daily-cloudcode-pa.googleapis.com"
	endpointProd          = "https://cloudcode-pa.googleapis.com"
	userAgentVersion      = "1.15.8"
	RequestUserAgent      = "antigravity"
	RequestTypeAgent      = "agent"
	SystemInstructionText = "You are Antigravity, a powerful agentic AI coding assistant designed by the Google Deepmind team working on Advanced Agentic Coding.You are pair programming with a USER to solve their coding task. The task may require creating a new codebase, modifying or debugging an existing codebase, or simply answering a question.**Absolute paths only****Proactiveness**"
	defaultAcceptHeader   = "application/json"
)

var Endpoints = []string{
	endpointDaily,
	endpointProd,
}

const clientMetadataHeader = `{"ideType":"IDE_UNSPECIFIED","platform":"PLATFORM_UNSPECIFIED","pluginType":"GEMINI"}`

func platformUserAgent() string {
	return fmt.Sprintf("antigravity/%s %s/%s", userAgentVersion, runtime.GOOS, runtime.GOARCH)
}

func ApplyHeaders(header http.Header, token string, accept string) {
	if accept == "" {
		accept = defaultAcceptHeader
	}
	header.Set("Authorization", "Bearer "+token)
	header.Set("Content-Type", "application/json")
	header.Set("User-Agent", platformUserAgent())
	header.Set("X-Goog-Api-Client", "google-cloud-sdk vscode_cloudshelleditor/0.1")
	header.Set("Client-Metadata", clientMetadataHeader)
	header.Set("Accept", accept)
}
