package cmd

import (
	"testing"
)

func TestDashboardsListCommand_Exists(t *testing.T) {
	cmd, _, err := dashboardsCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}
	if cmd.Name() != "list" {
		t.Errorf("expected command name 'list', got '%s'", cmd.Name())
	}
}

func TestDashboardsGetCommand_Exists(t *testing.T) {
	cmd, _, err := dashboardsCmd.Find([]string{"get"})
	if err != nil {
		t.Fatalf("get command not found: %v", err)
	}
	if cmd.Name() != "get" {
		t.Errorf("expected command name 'get', got '%s'", cmd.Name())
	}
}

func TestDashboardsGetCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := dashboardsCmd.Find([]string{"get"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"dashboard-name"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestDashboardsCreateCommand_Exists(t *testing.T) {
	cmd, _, err := dashboardsCmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}
	if cmd.Name() != "create" {
		t.Errorf("expected command name 'create', got '%s'", cmd.Name())
	}
}

func TestDashboardsDeleteCommand_Exists(t *testing.T) {
	cmd, _, err := dashboardsCmd.Find([]string{"delete"})
	if err != nil {
		t.Fatalf("delete command not found: %v", err)
	}
	if cmd.Name() != "delete" {
		t.Errorf("expected command name 'delete', got '%s'", cmd.Name())
	}
}
