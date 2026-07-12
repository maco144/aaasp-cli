// Package mcp implements a minimal MCP (Model Context Protocol) server over
// stdio — newline-delimited JSON-RPC 2.0, matching the wire format aaasp-ex's
// MCP client already speaks (see AaaspEx.Mcp.Transport.Stdio).
package mcp

import "encoding/json"

const ProtocolVersion = "2024-11-05"

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// IsNotification reports whether this message expects no response (no "id").
func (r *Request) IsNotification() bool {
	return len(r.ID) == 0
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	CodeParseError     = -32700
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603
)

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

type ToolCallParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ToolCallResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

func textResult(text string) *ToolCallResult {
	return &ToolCallResult{Content: []ContentBlock{{Type: "text", Text: text}}}
}

func errorResult(text string) *ToolCallResult {
	return &ToolCallResult{Content: []ContentBlock{{Type: "text", Text: text}}, IsError: true}
}
