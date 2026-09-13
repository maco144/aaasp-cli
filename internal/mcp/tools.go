package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/maco144/aaasp-cli/internal/api"
)

type toolHandler func(client *api.Client, args map[string]any) (*ToolCallResult, error)

type toolDef struct {
	Tool
	handler toolHandler
}

func buildTools() []toolDef {
	return []toolDef{
		{
			Tool: Tool{
				Name: "dispatch_run",
				Description: "Submit a goal to AAASP. The platform classifies the goal and " +
					"auto-routes it to the best-matching agent deployment, then starts a run. " +
					"Returns immediately with a run_id — poll get_run to check progress.",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"goal": map[string]any{
							"type":        "string",
							"description": "The task or question to send to an agent.",
						},
						"skill_hint": map[string]any{
							"type":        "string",
							"description": "Optional — skip goal classification and route directly to this skill/specialist.",
						},
					},
					"required": []string{"goal"},
				},
			},
			handler: dispatchRun,
		},
		{
			Tool: Tool{
				Name:        "get_run",
				Description: "Fetch the current status and result of a run by ID.",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"run_id": map[string]any{
							"type":        "string",
							"description": "The run ID returned by dispatch_run.",
						},
					},
					"required": []string{"run_id"},
				},
			},
			handler: getRun,
		},
		{
			Tool: Tool{
				Name: "list_deployments",
				Description: "List all agent deployments for the authenticated tenant. Each deployment " +
					"has `runnable` (bool) and `readiness_error` (string or null). A deployment with " +
					"runnable=false will fail every run; readiness_error says why (e.g. no LLM key for " +
					"its provider) and what to fix. Check it before dispatching to a deployment.",
				InputSchema: map[string]any{
					"type":       "object",
					"properties": map[string]any{},
				},
			},
			handler: listDeployments,
		},
		{
			Tool: Tool{
				Name:        "list_agents",
				Description: "List all agent definitions for the authenticated tenant.",
				InputSchema: map[string]any{
					"type":       "object",
					"properties": map[string]any{},
				},
			},
			handler: listAgents,
		},
	}
}

func dispatchRun(client *api.Client, args map[string]any) (*ToolCallResult, error) {
	goal, _ := args["goal"].(string)
	if goal == "" {
		return errorResult("goal is required"), nil
	}

	body := map[string]any{"goal": goal}
	if hint, ok := args["skill_hint"].(string); ok && hint != "" {
		body["skill_hint"] = hint
	}

	var result map[string]any
	if err := client.Post("/dispatch", body, &result); err != nil {
		return errorResult(apiErrText(err)), nil
	}

	return jsonResult(result)
}

func getRun(client *api.Client, args map[string]any) (*ToolCallResult, error) {
	runID, _ := args["run_id"].(string)
	if runID == "" {
		return errorResult("run_id is required"), nil
	}

	var result map[string]any
	if err := client.Get("/runs/"+runID, &result); err != nil {
		return errorResult(apiErrText(err)), nil
	}

	return jsonResult(result)
}

func listDeployments(client *api.Client, _ map[string]any) (*ToolCallResult, error) {
	var result map[string]any
	if err := client.Get("/deployments", &result); err != nil {
		return errorResult(apiErrText(err)), nil
	}

	return jsonResult(result["deployments"])
}

func listAgents(client *api.Client, _ map[string]any) (*ToolCallResult, error) {
	var result []map[string]any
	if err := client.Get("/agents", &result); err != nil {
		return errorResult(apiErrText(err)), nil
	}

	return jsonResult(result)
}

func jsonResult(v any) (*ToolCallResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return textResult(string(data)), nil
}

func apiErrText(err error) string {
	if apiErr, ok := err.(*api.APIError); ok {
		return fmt.Sprintf("AAASP API error (%d): %s", apiErr.Status, apiErr.Message)
	}
	return err.Error()
}
