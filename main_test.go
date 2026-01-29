package main

import (
	"reflect"
	"testing"
)

func TestFirstPositionalArg(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"empty", nil, ""},
		{"flag only", []string{"-m", "opus"}, ""},
		{"bool flag only", []string{"--raw"}, ""},
		{"positional first", []string{"hello", "world"}, "hello"},
		{"flag then positional", []string{"-m", "opus", "hello"}, "hello"},
		{"bool flag then positional", []string{"--raw", "hello"}, "hello"},
		{"mixed", []string{"--dry-run", "-m", "haiku", "how", "to"}, "how"},
		{"unknown flag skipped", []string{"--unknown"}, ""},
		{"unknown flag before positional", []string{"--unknown", "hello"}, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstPositionalArg(tt.args)
			if got != tt.want {
				t.Errorf("firstPositionalArg(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestReorderArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		want           []string
		wantPassthru   []string
	}{
		{
			"flags before positional",
			[]string{"-m", "opus", "how", "to", "rebase"},
			[]string{"-m", "opus", "--", "how", "to", "rebase"},
			nil,
		},
		{
			"flags after positional",
			[]string{"how", "to", "rebase", "-m", "opus"},
			[]string{"-m", "opus", "--", "how", "to", "rebase"},
			nil,
		},
		{
			"bool flag after positional",
			[]string{"how", "to", "rebase", "--raw"},
			[]string{"--raw", "--", "how", "to", "rebase"},
			nil,
		},
		{
			"mixed flags and positional",
			[]string{"--raw", "how", "to", "-m", "haiku", "rebase"},
			[]string{"--raw", "-m", "haiku", "--", "how", "to", "rebase"},
			nil,
		},
		{
			"only flags no positional",
			[]string{"-m", "opus", "--raw"},
			[]string{"-m", "opus", "--raw"},
			nil,
		},
		{
			"only positional",
			[]string{"hello", "world"},
			[]string{"--", "hello", "world"},
			nil,
		},
		{
			"dry-run with model after prompt",
			[]string{"--dry-run", "question", "-m", "opus"},
			[]string{"--dry-run", "-m", "opus", "--", "question"},
			nil,
		},
		{
			"empty args",
			nil,
			nil,
			nil,
		},
		{
			"unknown flag with value passthrough",
			[]string{"--output-format", "json", "question"},
			[]string{"--", "question"},
			[]string{"--output-format", "json"},
		},
		{
			"unknown flag with = passthrough",
			[]string{"--output-format=json", "question"},
			[]string{"--", "question"},
			[]string{"--output-format=json"},
		},
		{
			"unknown bool flag passthrough",
			[]string{"--verbose", "--raw", "question"},
			[]string{"--raw", "--", "question"},
			[]string{"--verbose"},
		},
		{
			"mixed ask and claude flags",
			[]string{"-m", "opus", "--system-prompt", "be concise", "--raw", "question"},
			[]string{"-m", "opus", "--raw", "--", "question"},
			[]string{"--system-prompt", "be concise"},
		},
		{
			"multiple passthrough flags",
			[]string{"--output-format", "json", "--max-budget-usd", "0.5", "question"},
			[]string{"--", "question"},
			[]string{"--output-format", "json", "--max-budget-usd", "0.5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passthrough = nil
			got := reorderArgs(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reorderArgs(%v) = %v, want %v", tt.args, got, tt.want)
			}
			if !reflect.DeepEqual(passthrough, tt.wantPassthru) {
				t.Errorf("passthrough = %v, want %v", passthrough, tt.wantPassthru)
			}
		})
	}
}

func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		pipeContent string
		want        string
	}{
		{"args only", []string{"how", "to", "rebase"}, "", "how to rebase"},
		{"pipe only", nil, "some data", "Here is some data. Please analyze it:\n\n```\nsome data\n```"},
		{"pipe and args", []string{"review"}, "diff output", "Here is the input data:\n\n```\ndiff output\n```\n\nreview"},
		{"empty", nil, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPrompt(tt.args, tt.pipeContent)
			if got != tt.want {
				t.Errorf("buildPrompt(%v, %q) = %q, want %q", tt.args, tt.pipeContent, got, tt.want)
			}
		})
	}
}
