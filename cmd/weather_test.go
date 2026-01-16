package cmd

import (
	"testing"
)

func TestWeatherCurrentCommand_Exists(t *testing.T) {
	cmd, _, err := weatherCmd.Find([]string{"current"})
	if err != nil {
		t.Fatalf("current command not found: %v", err)
	}
	if cmd.Name() != "current" {
		t.Errorf("expected command name 'current', got '%s'", cmd.Name())
	}
}

func TestWeatherForecastCommand_Exists(t *testing.T) {
	cmd, _, err := weatherCmd.Find([]string{"forecast"})
	if err != nil {
		t.Fatalf("forecast command not found: %v", err)
	}
	if cmd.Name() != "forecast" {
		t.Errorf("expected command name 'forecast', got '%s'", cmd.Name())
	}
}
