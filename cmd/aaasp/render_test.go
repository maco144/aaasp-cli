package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderDeploymentList_ShowsReadiness(t *testing.T) {
	deps := []any{
		map[string]any{"id": "dep-ok", "agent_def_id": "agent-1", "status": "active", "created_at": "2026-09-13T00:00:00Z", "runnable": true},
		map[string]any{
			"id": "dep-bad", "agent_def_id": "agent-2", "status": "active", "created_at": "2026-09-13T00:00:00Z",
			"runnable": false, "readiness_error": "This deployment has no anthropic key for claude-opus-4-1.",
		},
		// A server older than the readiness fields must not render as "not runnable".
		map[string]any{"id": "dep-old", "agent_def_id": "agent-3", "status": "active", "created_at": "2026-09-13T00:00:00Z"},
	}

	var buf bytes.Buffer
	renderDeploymentList(&buf, deps)
	out := buf.String()

	if !strings.Contains(out, "READY") {
		t.Fatalf("missing READY column:\n%s", out)
	}
	for _, want := range []string{"dep-ok", "yes", "dep-bad", "no", "dep-old", "?", "2026-09-13T00:00:00Z"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "dep-bad: This deployment has no anthropic key") {
		t.Errorf("readiness_error not listed under the table:\n%s", out)
	}
}

func TestRenderDeploymentList_NoFooterWhenAllRunnable(t *testing.T) {
	var buf bytes.Buffer
	renderDeploymentList(&buf, []any{map[string]any{"id": "dep-ok", "runnable": true}})
	if strings.Contains(buf.String(), "Not runnable") {
		t.Errorf("unexpected not-runnable footer:\n%s", buf.String())
	}
}

func TestRenderDeploymentDetail_ShowsReadinessAndRealFields(t *testing.T) {
	var buf bytes.Buffer
	renderDeploymentDetail(&buf, map[string]any{
		"id": "dep-bad", "agent_def_id": "agent-2", "status": "active", "credential_id": "cred-1",
		"created_at": "2026-09-13T00:00:00Z", "runnable": false, "readiness_error": "no key",
	})
	out := buf.String()
	for _, want := range []string{"dep-bad", "agent-2", "cred-1", "2026-09-13T00:00:00Z", "runnable", "no", "reason", "no key"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "schedule") {
		t.Errorf("deployments have no schedule field; should not be shown:\n%s", out)
	}
}

func TestRenderRunDetail_ShowsErrorResultAndCreated(t *testing.T) {
	var buf bytes.Buffer
	renderRunDetail(&buf, map[string]any{
		"id": "run-1", "status": "failed", "deployment_id": "dep-bad",
		"created_at": "2026-09-13T00:00:00Z", "error": "This deployment has no anthropic key",
	})
	out := buf.String()
	for _, want := range []string{"run-1", "failed", "2026-09-13T00:00:00Z", "error", "no anthropic key"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}

	buf.Reset()
	renderRunDetail(&buf, map[string]any{"id": "run-2", "status": "completed", "result": "hello world"})
	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("result not shown:\n%s", buf.String())
	}
}

func TestCredentialBody_SendsTheParamTheAPIReads(t *testing.T) {
	// The API maps "encrypted_key" to the stored key. The CLI used to send
	// "vault_key", which the API silently dropped, so every credential added via
	// the CLI was stored without a key.
	body := credentialBody("anthropic", "", "sk-test")

	if body["encrypted_key"] != "sk-test" {
		t.Errorf("encrypted_key = %q, want sk-test (body %v)", body["encrypted_key"], body)
	}
	if _, ok := body["vault_key"]; ok {
		t.Errorf("must not send the ignored vault_key param: %v", body)
	}
	// The API requires a label; default it rather than fail with a 422.
	if body["label"] != "anthropic" {
		t.Errorf("label = %q, want provider name as default", body["label"])
	}
	if got := credentialBody("openai", "work", "k")["label"]; got != "work" {
		t.Errorf("explicit label = %q, want work", got)
	}
}
