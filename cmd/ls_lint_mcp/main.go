package main

import (
	"context"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var Version = "dev"

type LintInput struct {
	Paths []string `json:"paths" jsonschema:"the paths to lint"`
}

type LintOutput struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func Lint(ctx context.Context, req *mcp.CallToolRequest, input LintInput) (
	*mcp.CallToolResult,
	LintOutput,
	error,
) {
	return nil, LintOutput{Greeting: "Hi " + strings.Join(input.Paths, ", ")}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "ls-lint-mcp",
		Version: Version,
	}, nil)

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "lint", Description: "lint paths"},
		Lint,
	)

	log.Println("Starting server 2")

	if err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	); err != nil {
		log.Fatal(err)
	}
}
