package cmd

import (
	"testing"
)

func TestSystemInfoCommand_Exists(t *testing.T) {
	cmd, _, err := systemCmd.Find([]string{"info"})
	if err != nil {
		t.Fatalf("info command not found: %v", err)
	}
	if cmd.Name() != "info" {
		t.Errorf("expected command name 'info', got '%s'", cmd.Name())
	}
}

func TestSystemRebootCommand_Exists(t *testing.T) {
	cmd, _, err := systemCmd.Find([]string{"reboot"})
	if err != nil {
		t.Fatalf("reboot command not found: %v", err)
	}
	if cmd.Name() != "reboot" {
		t.Errorf("expected command name 'reboot', got '%s'", cmd.Name())
	}
}

func TestSystemUsersCommand_Exists(t *testing.T) {
	cmd, _, err := systemCmd.Find([]string{"users"})
	if err != nil {
		t.Fatalf("users command not found: %v", err)
	}
	if cmd.Name() != "users" {
		t.Errorf("expected command name 'users', got '%s'", cmd.Name())
	}
}

func TestSystemInsightsCommand_Exists(t *testing.T) {
	cmd, _, err := systemCmd.Find([]string{"insights"})
	if err != nil {
		t.Fatalf("insights command not found: %v", err)
	}
	if cmd.Name() != "insights" {
		t.Errorf("expected command name 'insights', got '%s'", cmd.Name())
	}
}

func TestSystemNameCommand_Exists(t *testing.T) {
	cmd, _, err := systemCmd.Find([]string{"name"})
	if err != nil {
		t.Fatalf("name command not found: %v", err)
	}
	if cmd.Name() != "name" {
		t.Errorf("expected command name 'name', got '%s'", cmd.Name())
	}
}

func TestSystemNameGetCommand_Exists(t *testing.T) {
	nameCmd, _, err := systemCmd.Find([]string{"name"})
	if err != nil {
		t.Fatalf("name command not found: %v", err)
	}

	cmd, _, err := nameCmd.Find([]string{"get"})
	if err != nil {
		t.Fatalf("name get command not found: %v", err)
	}
	if cmd.Name() != "get" {
		t.Errorf("expected command name 'get', got '%s'", cmd.Name())
	}
}

func TestSystemNameSetCommand_Exists(t *testing.T) {
	nameCmd, _, err := systemCmd.Find([]string{"name"})
	if err != nil {
		t.Fatalf("name command not found: %v", err)
	}

	cmd, _, err := nameCmd.Find([]string{"set"})
	if err != nil {
		t.Fatalf("name set command not found: %v", err)
	}
	if cmd.Name() != "set" {
		t.Errorf("expected command name 'set', got '%s'", cmd.Name())
	}
}

func TestSystemNameSetCommand_RequiresOneArg(t *testing.T) {
	nameCmd, _, _ := systemCmd.Find([]string{"name"})
	cmd, _, _ := nameCmd.Find([]string{"set"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"New Name"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}
