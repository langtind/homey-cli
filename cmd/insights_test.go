package cmd

import (
	"testing"
)

func TestInsightsListCommand_Exists(t *testing.T) {
	cmd, _, err := insightsCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}
	if cmd.Name() != "list" {
		t.Errorf("expected command name 'list', got '%s'", cmd.Name())
	}
}

func TestInsightsGetCommand_Exists(t *testing.T) {
	cmd, _, err := insightsCmd.Find([]string{"get"})
	if err != nil {
		t.Fatalf("get command not found: %v", err)
	}
	if cmd.Name() != "get" {
		t.Errorf("expected command name 'get', got '%s'", cmd.Name())
	}
}

func TestInsightsGetCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := insightsCmd.Find([]string{"get"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"log-id"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestInsightsDeleteCommand_Exists(t *testing.T) {
	cmd, _, err := insightsCmd.Find([]string{"delete"})
	if err != nil {
		t.Fatalf("delete command not found: %v", err)
	}
	if cmd.Name() != "delete" {
		t.Errorf("expected command name 'delete', got '%s'", cmd.Name())
	}
}

func TestInsightsDeleteCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := insightsCmd.Find([]string{"delete"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"log-id"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestInsightsClearCommand_Exists(t *testing.T) {
	cmd, _, err := insightsCmd.Find([]string{"clear"})
	if err != nil {
		t.Fatalf("clear command not found: %v", err)
	}
	if cmd.Name() != "clear" {
		t.Errorf("expected command name 'clear', got '%s'", cmd.Name())
	}
}

func TestInsightsClearCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := insightsCmd.Find([]string{"clear"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"log-id"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}
