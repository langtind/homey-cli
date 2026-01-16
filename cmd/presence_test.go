package cmd

import (
	"testing"
)

func TestPresenceGetCommand_Exists(t *testing.T) {
	cmd, _, err := presenceCmd.Find([]string{"get"})
	if err != nil {
		t.Fatalf("get command not found: %v", err)
	}
	if cmd.Name() != "get" {
		t.Errorf("expected command name 'get', got '%s'", cmd.Name())
	}
}

func TestPresenceGetCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := presenceCmd.Find([]string{"get"})

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

func TestPresenceSetCommand_Exists(t *testing.T) {
	cmd, _, err := presenceCmd.Find([]string{"set"})
	if err != nil {
		t.Fatalf("set command not found: %v", err)
	}
	if cmd.Name() != "set" {
		t.Errorf("expected command name 'set', got '%s'", cmd.Name())
	}
}

func TestPresenceSetCommand_RequiresTwoArgs(t *testing.T) {
	cmd, _, _ := presenceCmd.Find([]string{"set"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"user"})
	if err == nil {
		t.Error("expected error with 1 arg")
	}

	err = cmd.Args(cmd, []string{"user", "home"})
	if err != nil {
		t.Errorf("expected no error with 2 args, got: %v", err)
	}
}

func TestPresenceAsleepGetCommand_Exists(t *testing.T) {
	asleepCmd, _, err := presenceCmd.Find([]string{"asleep"})
	if err != nil {
		t.Fatalf("asleep command not found: %v", err)
	}

	cmd, _, err := asleepCmd.Find([]string{"get"})
	if err != nil {
		t.Fatalf("asleep get command not found: %v", err)
	}
	if cmd.Name() != "get" {
		t.Errorf("expected command name 'get', got '%s'", cmd.Name())
	}
}

func TestPresenceAsleepSetCommand_Exists(t *testing.T) {
	asleepCmd, _, err := presenceCmd.Find([]string{"asleep"})
	if err != nil {
		t.Fatalf("asleep command not found: %v", err)
	}

	cmd, _, err := asleepCmd.Find([]string{"set"})
	if err != nil {
		t.Fatalf("asleep set command not found: %v", err)
	}
	if cmd.Name() != "set" {
		t.Errorf("expected command name 'set', got '%s'", cmd.Name())
	}
}
