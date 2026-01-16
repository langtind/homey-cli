package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// Device represents a Homey device
type Device struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Class           string                `json:"class"`
	Zone            string                `json:"zone"`
	CapabilitiesObj map[string]Capability `json:"capabilitiesObj"`
}

// Capability represents a device capability
type Capability struct {
	ID    string      `json:"id"`
	Value interface{} `json:"value"`
	Title string      `json:"title"`
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Manage devices",
	Long:  `List, view, control, and manage Homey devices.`,
}

var devicesMatchFilter string

// findDevice finds a device by name or ID from the list of all devices
func findDevice(nameOrID string) (*Device, error) {
	data, err := apiClient.GetDevices()
	if err != nil {
		return nil, err
	}

	var devices map[string]Device
	if err := json.Unmarshal(data, &devices); err != nil {
		return nil, fmt.Errorf("failed to parse devices: %w", err)
	}

	for _, d := range devices {
		if d.ID == nameOrID || strings.EqualFold(d.Name, nameOrID) {
			return &d, nil
		}
	}

	return nil, fmt.Errorf("device not found: %s", nameOrID)
}

var devicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all devices",
	Long: `List all devices, optionally filtered by name.

Examples:
  homeyctl devices list
  homeyctl devices list --match "kitchen"
  homeyctl devices list --match "light"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetDevices()
		if err != nil {
			return err
		}

		var devices map[string]Device
		if err := json.Unmarshal(data, &devices); err != nil {
			return fmt.Errorf("failed to parse devices: %w", err)
		}

		// Filter devices if --match is provided
		var filtered []Device
		for _, d := range devices {
			if devicesMatchFilter == "" || strings.Contains(strings.ToLower(d.Name), strings.ToLower(devicesMatchFilter)) {
				filtered = append(filtered, d)
			}
		}

		if isTableFormat() {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tCLASS\tID")
			fmt.Fprintln(w, "----\t-----\t--")
			for _, d := range filtered {
				fmt.Fprintf(w, "%s\t%s\t%s\n", d.Name, d.Class, d.ID)
			}
			w.Flush()
			return nil
		}

		out, _ := json.MarshalIndent(filtered, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

var devicesGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get device details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		if isTableFormat() {
			fmt.Printf("Name:  %s\n", device.Name)
			fmt.Printf("Class: %s\n", device.Class)
			fmt.Printf("ID:    %s\n", device.ID)
			fmt.Println("\nCapabilities:")

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "  CAPABILITY\tVALUE")
			fmt.Fprintln(w, "  ----------\t-----")
			for _, cap := range device.CapabilitiesObj {
				fmt.Fprintf(w, "  %s\t%v\n", cap.ID, cap.Value)
			}
			w.Flush()
			return nil
		}

		out, _ := json.MarshalIndent(device, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

var devicesValuesCmd = &cobra.Command{
	Use:   "values <name-or-id>",
	Short: "Get all capability values for a device",
	Long: `Get all current capability values for a device.

Useful for multi-sensors and devices with many capabilities.

Examples:
  homeyctl devices values "PultLED"
  homeyctl devices values "Multisensor 6"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		if isTableFormat() {
			fmt.Printf("Values for %s:\n\n", device.Name)
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "CAPABILITY\tVALUE")
			fmt.Fprintln(w, "----------\t-----")
			for _, cap := range device.CapabilitiesObj {
				fmt.Fprintf(w, "%s\t%v\n", cap.ID, cap.Value)
			}
			w.Flush()
			return nil
		}

		// JSON output - just the values
		values := make(map[string]interface{})
		for _, cap := range device.CapabilitiesObj {
			values[cap.ID] = cap.Value
		}
		out, _ := json.MarshalIndent(map[string]interface{}{
			"id":     device.ID,
			"name":   device.Name,
			"values": values,
		}, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.AddCommand(devicesListCmd)
	devicesListCmd.Flags().StringVar(&devicesMatchFilter, "match", "", "Filter devices by name (case-insensitive)")
	devicesCmd.AddCommand(devicesGetCmd)
	devicesCmd.AddCommand(devicesValuesCmd)
}
