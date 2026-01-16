package cmd

import (
	"testing"
)

func TestNotifySendCommand_Exists(t *testing.T) {
	cmd, _, err := notifyCmd.Find([]string{"send"})
	if err != nil {
		t.Fatalf("send command not found: %v", err)
	}
	if cmd.Name() != "send" {
		t.Errorf("expected command name 'send', got '%s'", cmd.Name())
	}
}

func TestNotifySendCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := notifyCmd.Find([]string{"send"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"Hello world!"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestNotifyListCommand_Exists(t *testing.T) {
	cmd, _, err := notifyCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}
	if cmd.Name() != "list" {
		t.Errorf("expected command name 'list', got '%s'", cmd.Name())
	}
}

func TestNotifyDeleteCommand_Exists(t *testing.T) {
	cmd, _, err := notifyCmd.Find([]string{"delete"})
	if err != nil {
		t.Fatalf("delete command not found: %v", err)
	}
	if cmd.Name() != "delete" {
		t.Errorf("expected command name 'delete', got '%s'", cmd.Name())
	}
}

func TestNotifyDeleteCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := notifyCmd.Find([]string{"delete"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"notification-id"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestNotifyClearCommand_Exists(t *testing.T) {
	cmd, _, err := notifyCmd.Find([]string{"clear"})
	if err != nil {
		t.Fatalf("clear command not found: %v", err)
	}
	if cmd.Name() != "clear" {
		t.Errorf("expected command name 'clear', got '%s'", cmd.Name())
	}
}

func TestNotifyOwnersCommand_Exists(t *testing.T) {
	cmd, _, err := notifyCmd.Find([]string{"owners"})
	if err != nil {
		t.Fatalf("owners command not found: %v", err)
	}
	if cmd.Name() != "owners" {
		t.Errorf("expected command name 'owners', got '%s'", cmd.Name())
	}
}
