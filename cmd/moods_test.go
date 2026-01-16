package cmd

import (
	"testing"
)

func TestMoodsListCommand_Exists(t *testing.T) {
	cmd, _, err := moodsCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}
	if cmd.Name() != "list" {
		t.Errorf("expected command name 'list', got '%s'", cmd.Name())
	}
}

func TestMoodsGetCommand_Exists(t *testing.T) {
	cmd, _, err := moodsCmd.Find([]string{"get"})
	if err != nil {
		t.Fatalf("get command not found: %v", err)
	}
	if cmd.Name() != "get" {
		t.Errorf("expected command name 'get', got '%s'", cmd.Name())
	}
}

func TestMoodsGetCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := moodsCmd.Find([]string{"get"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"mood-name"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestMoodsCreateCommand_Exists(t *testing.T) {
	cmd, _, err := moodsCmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}
	if cmd.Name() != "create" {
		t.Errorf("expected command name 'create', got '%s'", cmd.Name())
	}
}

func TestMoodsDeleteCommand_Exists(t *testing.T) {
	cmd, _, err := moodsCmd.Find([]string{"delete"})
	if err != nil {
		t.Fatalf("delete command not found: %v", err)
	}
	if cmd.Name() != "delete" {
		t.Errorf("expected command name 'delete', got '%s'", cmd.Name())
	}
}

func TestMoodsSetCommand_Exists(t *testing.T) {
	cmd, _, err := moodsCmd.Find([]string{"set"})
	if err != nil {
		t.Fatalf("set command not found: %v", err)
	}
	if cmd.Name() != "set" {
		t.Errorf("expected command name 'set', got '%s'", cmd.Name())
	}
}

func TestMoodsSetCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := moodsCmd.Find([]string{"set"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"mood-name"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}
