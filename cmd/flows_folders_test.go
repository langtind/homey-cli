package cmd

import (
	"testing"
)

func TestFlowFoldersListCommand_Exists(t *testing.T) {
	cmd, _, err := flowsFoldersCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}
	if cmd.Name() != "list" {
		t.Errorf("expected command name 'list', got '%s'", cmd.Name())
	}
}

func TestFlowFoldersGetCommand_Exists(t *testing.T) {
	cmd, _, err := flowsFoldersCmd.Find([]string{"get"})
	if err != nil {
		t.Fatalf("get command not found: %v", err)
	}
	if cmd.Name() != "get" {
		t.Errorf("expected command name 'get', got '%s'", cmd.Name())
	}
}

func TestFlowFoldersGetCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := flowsFoldersCmd.Find([]string{"get"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"folder-id"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestFlowFoldersCreateCommand_Exists(t *testing.T) {
	cmd, _, err := flowsFoldersCmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}
	if cmd.Name() != "create" {
		t.Errorf("expected command name 'create', got '%s'", cmd.Name())
	}
}

func TestFlowFoldersCreateCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := flowsFoldersCmd.Find([]string{"create"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"folder-name"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestFlowFoldersDeleteCommand_Exists(t *testing.T) {
	cmd, _, err := flowsFoldersCmd.Find([]string{"delete"})
	if err != nil {
		t.Fatalf("delete command not found: %v", err)
	}
	if cmd.Name() != "delete" {
		t.Errorf("expected command name 'delete', got '%s'", cmd.Name())
	}
}

func TestFlowFoldersUpdateCommand_Exists(t *testing.T) {
	cmd, _, err := flowsFoldersCmd.Find([]string{"update"})
	if err != nil {
		t.Fatalf("update command not found: %v", err)
	}
	if cmd.Name() != "update" {
		t.Errorf("expected command name 'update', got '%s'", cmd.Name())
	}
}

func TestFlowFoldersUpdateCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := flowsFoldersCmd.Find([]string{"update"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"folder-name"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}
