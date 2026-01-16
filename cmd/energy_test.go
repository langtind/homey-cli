package cmd

import (
	"testing"
)

func TestEnergyLiveCommand_Exists(t *testing.T) {
	cmd, _, err := energyCmd.Find([]string{"live"})
	if err != nil {
		t.Fatalf("live command not found: %v", err)
	}
	if cmd.Name() != "live" {
		t.Errorf("expected command name 'live', got '%s'", cmd.Name())
	}
}

func TestEnergyReportCommand_Exists(t *testing.T) {
	cmd, _, err := energyCmd.Find([]string{"report"})
	if err != nil {
		t.Fatalf("report command not found: %v", err)
	}
	if cmd.Name() != "report" {
		t.Errorf("expected command name 'report', got '%s'", cmd.Name())
	}
}

func TestEnergyReportYearCommand_Exists(t *testing.T) {
	reportCmd, _, err := energyCmd.Find([]string{"report"})
	if err != nil {
		t.Fatalf("report command not found: %v", err)
	}

	cmd, _, err := reportCmd.Find([]string{"year"})
	if err != nil {
		t.Fatalf("report year command not found: %v", err)
	}
	if cmd.Name() != "year" {
		t.Errorf("expected command name 'year', got '%s'", cmd.Name())
	}
}

func TestEnergyDeleteCommand_Exists(t *testing.T) {
	cmd, _, err := energyCmd.Find([]string{"delete"})
	if err != nil {
		t.Fatalf("delete command not found: %v", err)
	}
	if cmd.Name() != "delete" {
		t.Errorf("expected command name 'delete', got '%s'", cmd.Name())
	}
}

func TestEnergyCurrencyCommand_Exists(t *testing.T) {
	cmd, _, err := energyCmd.Find([]string{"currency"})
	if err != nil {
		t.Fatalf("currency command not found: %v", err)
	}
	if cmd.Name() != "currency" {
		t.Errorf("expected command name 'currency', got '%s'", cmd.Name())
	}
}

func TestEnergyPriceCommand_Exists(t *testing.T) {
	cmd, _, err := energyCmd.Find([]string{"price"})
	if err != nil {
		t.Fatalf("price command not found: %v", err)
	}
	if cmd.Name() != "price" {
		t.Errorf("expected command name 'price', got '%s'", cmd.Name())
	}
}

func TestEnergyPriceSetCommand_Exists(t *testing.T) {
	priceCmd, _, err := energyCmd.Find([]string{"price"})
	if err != nil {
		t.Fatalf("price command not found: %v", err)
	}

	cmd, _, err := priceCmd.Find([]string{"set"})
	if err != nil {
		t.Fatalf("price set command not found: %v", err)
	}
	if cmd.Name() != "set" {
		t.Errorf("expected command name 'set', got '%s'", cmd.Name())
	}
}

func TestEnergyPriceTypeCommand_Exists(t *testing.T) {
	priceCmd, _, err := energyCmd.Find([]string{"price"})
	if err != nil {
		t.Fatalf("price command not found: %v", err)
	}

	cmd, _, err := priceCmd.Find([]string{"type"})
	if err != nil {
		t.Fatalf("price type command not found: %v", err)
	}
	if cmd.Name() != "type" {
		t.Errorf("expected command name 'type', got '%s'", cmd.Name())
	}
}
