package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/maco144/aaasp-cli/internal/api"
)

// Server runs the MCP stdio loop: read one JSON-RPC message per line from
// in, dispatch it, write one JSON-RPC response per line to out. Logging goes
// to a separate writer since stdout is reserved for protocol messages.
type Server struct {
	client  *api.Client
	tools   []toolDef
	version string
	logger  *log.Logger
}

func NewServer(client *api.Client, version string, logOut io.Writer) *Server {
	return &Server{
		client:  client,
		tools:   buildTools(),
		version: version,
		logger:  log.New(logOut, "[mcp] ", log.LstdFlags),
	}
}

func (s *Server) Serve(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			s.logger.Printf("parse error: %v", err)
			s.writeResponse(out, &Response{
				JSONRPC: "2.0",
				Error:   &RPCError{Code: CodeParseError, Message: "invalid JSON"},
			})
			continue
		}

		resp := s.handle(&req)
		if resp != nil {
			s.writeResponse(out, resp)
		}
	}

	return scanner.Err()
}

func (s *Server) handle(req *Request) *Response {
	switch req.Method {
	case "initialize":
		return s.reply(req, map[string]any{
			"protocolVersion": ProtocolVersion,
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    "aaasp",
				"version": s.version,
			},
		})

	case "notifications/initialized", "notifications/cancelled":
		return nil // notifications get no response

	case "ping":
		return s.reply(req, map[string]any{})

	case "tools/list":
		tools := make([]Tool, len(s.tools))
		for i, t := range s.tools {
			tools[i] = t.Tool
		}
		return s.reply(req, map[string]any{"tools": tools})

	case "tools/call":
		return s.handleToolCall(req)

	default:
		if req.IsNotification() {
			return nil
		}
		return &Response{JSONRPC: "2.0", ID: req.ID, Error: &RPCError{
			Code:    CodeMethodNotFound,
			Message: fmt.Sprintf("method not found: %s", req.Method),
		}}
	}
}

func (s *Server) handleToolCall(req *Request) *Response {
	var params ToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &Response{JSONRPC: "2.0", ID: req.ID, Error: &RPCError{
			Code: CodeInvalidParams, Message: "invalid tools/call params",
		}}
	}

	for _, t := range s.tools {
		if t.Name != params.Name {
			continue
		}

		result, err := t.handler(s.client, params.Arguments)
		if err != nil {
			s.logger.Printf("tool %s failed: %v", params.Name, err)
			result = errorResult(err.Error())
		}
		return s.reply(req, result)
	}

	return &Response{JSONRPC: "2.0", ID: req.ID, Error: &RPCError{
		Code:    CodeMethodNotFound,
		Message: fmt.Sprintf("unknown tool: %s", params.Name),
	}}
}

func (s *Server) reply(req *Request, result any) *Response {
	if req.IsNotification() {
		return nil
	}
	return &Response{JSONRPC: "2.0", ID: req.ID, Result: result}
}

func (s *Server) writeResponse(out io.Writer, resp *Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		s.logger.Printf("marshal error: %v", err)
		return
	}
	data = append(data, '\n')
	if _, err := out.Write(data); err != nil {
		s.logger.Printf("write error: %v", err)
	}
}
