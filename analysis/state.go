package analysis

import (
	"fmt"
	"strings"
	"unicode"
	"github.com/ViktorTomkovic/go-firstlsp/lsp"
)

type State struct {
	// Map of file names to contents
	Documents map[string]string
}

func NewState() State {
	return State {Documents: map[string]string{}}
}

func getDiagnostics(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		if idx := strings.Index(line, "VS Code"); idx >= 0 {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range: LineRange(row, idx, idx+len("VS Code")),
				Severity: 2,
				Source: "Common Sense",
				Message: "We do not speak profanities here.",
			})
		}
		lineLen := len(line)
		trimLen := len(strings.TrimRightFunc(line, unicode.IsSpace))
		if trimLen < lineLen {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range: LineRange(row, trimLen, lineLen),
				Severity: 2,
				Source: "Whitespace check",
				Message: "Trailing whitespace",
			})
		}
	}
	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnostics(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnostics(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := s.Documents[uri]
	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID: &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID: &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line: position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line: position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) CodeAction(id int, uri string) lsp.CodeActionResponse {
	text := s.Documents[uri]
	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range: LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Use more customazible editor",
				Edit: &lsp.WorkspaceEdit{
					Changes: replaceChange,
				},
			})
			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range: LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Censor profanity",
				Edit: &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}
	}
	return lsp.CodeActionResponse {
		Response: lsp.Response{
			RPC: "2.0",
			ID: &id,
		},
		Result: actions,
	}
}

func (s *State) Completion(id int, uri string) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label: "Neovim",
			Detail: "Very cool editor",
			Documentation: "More description about how awesome this editor is.",
		},
		{
			Label: "Neovim2",
			Detail: "Very cool editor2",
			Documentation: "More desc2ription about how awesome this editor is.",
		},
	}
	return lsp.CompletionResponse {
		Response: lsp.Response{
			RPC: "2.0",
			ID: &id,
		},
		Result: items,
	}
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line: line,
			Character: start,
		},
		End: lsp.Position{
			Line: line,
			Character: end,
		},
	}
}

