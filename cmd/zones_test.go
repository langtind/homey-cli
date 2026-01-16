package cmd

import (
	"testing"
)

func TestZonesRenameCommand_Exists(t *testing.T) {
	cmd, _, err := zonesCmd.Find([]string{"rename"})
	if err != nil {
		t.Fatalf("rename command not found: %v", err)
	}
	if cmd.Name() != "rename" {
		t.Errorf("expected command name 'rename', got '%s'", cmd.Name())
	}
}

func TestZonesRenameCommand_RequiresTwoArgs(t *testing.T) {
	cmd, _, _ := zonesCmd.Find([]string{"rename"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	// Test with wrong number of args
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"one"})
	if err == nil {
		t.Error("expected error with 1 arg")
	}

	err = cmd.Args(cmd, []string{"one", "two", "three"})
	if err == nil {
		t.Error("expected error with 3 args")
	}

	// Test with correct number of args
	err = cmd.Args(cmd, []string{"zone-name", "new-name"})
	if err != nil {
		t.Errorf("expected no error with 2 args, got: %v", err)
	}
}

func TestZonesRenameCommand_Usage(t *testing.T) {
	cmd, _, _ := zonesCmd.Find([]string{"rename"})

	expected := "rename <name-or-id> <new-name>"
	if cmd.Use != expected {
		t.Errorf("expected Use '%s', got '%s'", expected, cmd.Use)
	}
}

func TestZonesDeleteCommand_Exists(t *testing.T) {
	cmd, _, err := zonesCmd.Find([]string{"delete"})
	if err != nil {
		t.Fatalf("delete command not found: %v", err)
	}
	if cmd.Name() != "delete" {
		t.Errorf("expected command name 'delete', got '%s'", cmd.Name())
	}
}

func TestZonesListCommand_Exists(t *testing.T) {
	cmd, _, err := zonesCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}
	if cmd.Name() != "list" {
		t.Errorf("expected command name 'list', got '%s'", cmd.Name())
	}
}

func TestZonesIconsCommand_Exists(t *testing.T) {
	cmd, _, err := zonesCmd.Find([]string{"icons"})
	if err != nil {
		t.Fatalf("icons command not found: %v", err)
	}
	if cmd.Name() != "icons" {
		t.Errorf("expected command name 'icons', got '%s'", cmd.Name())
	}
}

func TestZonesSetIconCommand_Exists(t *testing.T) {
	cmd, _, err := zonesCmd.Find([]string{"set-icon"})
	if err != nil {
		t.Fatalf("set-icon command not found: %v", err)
	}
	if cmd.Name() != "set-icon" {
		t.Errorf("expected command name 'set-icon', got '%s'", cmd.Name())
	}
}

func TestZonesSetIconCommand_RequiresTwoArgs(t *testing.T) {
	cmd, _, _ := zonesCmd.Find([]string{"set-icon"})

	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	err = cmd.Args(cmd, []string{"zone-name", "icon"})
	if err != nil {
		t.Errorf("expected no error with 2 args, got: %v", err)
	}
}

func TestZonesRenameCommand_HasIconFlag(t *testing.T) {
	cmd, _, _ := zonesCmd.Find([]string{"rename"})

	flag := cmd.Flags().Lookup("icon")
	if flag == nil {
		t.Error("expected --icon flag to be defined")
	}
}

func TestKnownZoneIcons_NotEmpty(t *testing.T) {
	if len(knownZoneIcons) == 0 {
		t.Error("expected knownZoneIcons to not be empty")
	}

	// Check that common icons are present
	expectedIcons := []string{"home", "livingRoom", "kitchen", "bedroom", "office"}
	for _, expected := range expectedIcons {
		found := false
		for _, icon := range knownZoneIcons {
			if icon == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected icon '%s' to be in knownZoneIcons", expected)
		}
	}
}
