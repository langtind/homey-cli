package cmd

import (
	"testing"
)

func TestUsersListCommand_Exists(t *testing.T) {
	cmd, _, err := usersCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}
	if cmd.Name() != "list" {
		t.Errorf("expected command name 'list', got '%s'", cmd.Name())
	}
}

func TestUsersGetCommand_Exists(t *testing.T) {
	cmd, _, err := usersCmd.Find([]string{"get"})
	if err != nil {
		t.Fatalf("get command not found: %v", err)
	}
	if cmd.Name() != "get" {
		t.Errorf("expected command name 'get', got '%s'", cmd.Name())
	}
}

func TestUsersGetCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := usersCmd.Find([]string{"get"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"user-name"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestUsersMeCommand_Exists(t *testing.T) {
	cmd, _, err := usersCmd.Find([]string{"me"})
	if err != nil {
		t.Fatalf("me command not found: %v", err)
	}
	if cmd.Name() != "me" {
		t.Errorf("expected command name 'me', got '%s'", cmd.Name())
	}
}

func TestUsersCreateCommand_Exists(t *testing.T) {
	cmd, _, err := usersCmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}
	if cmd.Name() != "create" {
		t.Errorf("expected command name 'create', got '%s'", cmd.Name())
	}
}

func TestUsersDeleteCommand_Exists(t *testing.T) {
	cmd, _, err := usersCmd.Find([]string{"delete"})
	if err != nil {
		t.Fatalf("delete command not found: %v", err)
	}
	if cmd.Name() != "delete" {
		t.Errorf("expected command name 'delete', got '%s'", cmd.Name())
	}
}
