package mcp

import (
	"strings"
	"testing"
)

func TestListDeploymentsDescription_TellsAgentsAboutReadiness(t *testing.T) {
	// The tool relays the API JSON, which carries runnable/readiness_error. An
	// agent only uses those fields if the tool description says they exist.
	for _, td := range buildTools() {
		if td.Name != "list_deployments" {
			continue
		}
		for _, want := range []string{"runnable", "readiness_error"} {
			if !strings.Contains(td.Description, want) {
				t.Errorf("list_deployments description missing %q: %s", want, td.Description)
			}
		}
		return
	}
	t.Fatal("list_deployments tool not found")
}
