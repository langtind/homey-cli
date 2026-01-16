package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// parseValue converts a string value to the appropriate type (bool, number, or string)
func parseValue(valueStr string) interface{} {
	if valueStr == "true" {
		return true
	}
	if valueStr == "false" {
		return false
	}

	// Try as number
	var num float64
	if _, err := fmt.Sscanf(valueStr, "%f", &num); err == nil {
		return num
	}

	return valueStr
}

var devicesSetCmd = &cobra.Command{
	Use:   "set <name-or-id> <capability> <value>",
	Short: "Set device capability",
	Long: `Set a device capability value.

Examples:
  homeyctl devices set "PultLED" onoff true
  homeyctl devices set "PultLED" dim 0.5
  homeyctl devices set "Aksels rom" target_temperature 22`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]
		capability := args[1]
		valueStr := args[2]

		device, err := findDevice(nameOrID)
		if err != nil {
			return err
		}

		value := parseValue(valueStr)

		if err := apiClient.SetCapability(device.ID, capability, value); err != nil {
			return err
		}

		fmt.Printf("Set %s.%s = %v\n", device.Name, capability, value)
		return nil
	},
}

var devicesOnCmd = &cobra.Command{
	Use:   "on <name-or-id>",
	Short: "Turn device on",
	Long: `Turn a device on (shorthand for 'devices set <name> onoff true').

Examples:
  homeyctl devices on "Living Room Light"
  homeyctl devices on "Aksels rom"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return setDeviceOnOff(args[0], true)
	},
}

var devicesOffCmd = &cobra.Command{
	Use:   "off <name-or-id>",
	Short: "Turn device off",
	Long: `Turn a device off (shorthand for 'devices set <name> onoff false').

Examples:
  homeyctl devices off "Living Room Light"
  homeyctl devices off "Aksels rom"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return setDeviceOnOff(args[0], false)
	},
}

func setDeviceOnOff(nameOrID string, on bool) error {
	device, err := findDevice(nameOrID)
	if err != nil {
		return err
	}

	// Check if device supports onoff
	if _, hasOnOff := device.CapabilitiesObj["onoff"]; !hasOnOff {
		return fmt.Errorf("device '%s' does not support on/off", device.Name)
	}

	if err := apiClient.SetCapability(device.ID, "onoff", on); err != nil {
		return err
	}

	state := "on"
	if !on {
		state = "off"
	}
	fmt.Printf("Turned %s %s\n", device.Name, state)
	return nil
}

func init() {
	devicesCmd.AddCommand(devicesSetCmd)
	devicesCmd.AddCommand(devicesOnCmd)
	devicesCmd.AddCommand(devicesOffCmd)
}
