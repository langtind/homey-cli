package cmd

import (
	"testing"
)

func TestAppsListCommand_Exists(t *testing.T) {
	cmd, _, err := appsCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}
	if cmd.Name() != "list" {
		t.Errorf("expected command name 'list', got '%s'", cmd.Name())
	}
}

func TestAppsGetCommand_Exists(t *testing.T) {
	cmd, _, err := appsCmd.Find([]string{"get"})
	if err != nil {
		t.Fatalf("get command not found: %v", err)
	}
	if cmd.Name() != "get" {
		t.Errorf("expected command name 'get', got '%s'", cmd.Name())
	}
}

func TestAppsRestartCommand_Exists(t *testing.T) {
	cmd, _, err := appsCmd.Find([]string{"restart"})
	if err != nil {
		t.Fatalf("restart command not found: %v", err)
	}
	if cmd.Name() != "restart" {
		t.Errorf("expected command name 'restart', got '%s'", cmd.Name())
	}
}

func TestAppsInstallCommand_Exists(t *testing.T) {
	cmd, _, err := appsCmd.Find([]string{"install"})
	if err != nil {
		t.Fatalf("install command not found: %v", err)
	}
	if cmd.Name() != "install" {
		t.Errorf("expected command name 'install', got '%s'", cmd.Name())
	}
}

func TestAppsInstallCommand_RequiresOneArg(t *testing.T) {
	cmd, _, _ := appsCmd.Find([]string{"install"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"com.app.id"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}
}

func TestAppsUninstallCommand_Exists(t *testing.T) {
	cmd, _, err := appsCmd.Find([]string{"uninstall"})
	if err != nil {
		t.Fatalf("uninstall command not found: %v", err)
	}
	if cmd.Name() != "uninstall" {
		t.Errorf("expected command name 'uninstall', got '%s'", cmd.Name())
	}
}

func TestAppsEnableCommand_Exists(t *testing.T) {
	cmd, _, err := appsCmd.Find([]string{"enable"})
	if err != nil {
		t.Fatalf("enable command not found: %v", err)
	}
	if cmd.Name() != "enable" {
		t.Errorf("expected command name 'enable', got '%s'", cmd.Name())
	}
}

func TestAppsDisableCommand_Exists(t *testing.T) {
	cmd, _, err := appsCmd.Find([]string{"disable"})
	if err != nil {
		t.Fatalf("disable command not found: %v", err)
	}
	if cmd.Name() != "disable" {
		t.Errorf("expected command name 'disable', got '%s'", cmd.Name())
	}
}

func TestAppsSettingsCommand_Exists(t *testing.T) {
	cmd, _, err := appsCmd.Find([]string{"settings"})
	if err != nil {
		t.Fatalf("settings command not found: %v", err)
	}
	if cmd.Name() != "settings" {
		t.Errorf("expected command name 'settings', got '%s'", cmd.Name())
	}
}

func TestAppsUsageCommand_Exists(t *testing.T) {
	cmd, _, err := appsCmd.Find([]string{"usage"})
	if err != nil {
		t.Fatalf("usage command not found: %v", err)
	}
	if cmd.Name() != "usage" {
		t.Errorf("expected command name 'usage', got '%s'", cmd.Name())
	}
}
