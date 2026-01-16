package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var devicesRenameCmd = &cobra.Command{
	Use:   "rename <name-or-id> <new-name>",
	Short: "Rename a device",
	Long: `Rename a device.

Examples:
  homeyctl devices rename "Old Name" "New Name"
  homeyctl devices rename abc123-device-id "New Name"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]
		newName := args[1]

		device, err := findDevice(nameOrID)
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"name": newName,
		}

		if err := apiClient.UpdateDevice(device.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Renamed device '%s' to '%s'\n", device.Name, newName)
		return nil
	},
}

var devicesMoveCmd = &cobra.Command{
	Use:   "move <device> <zone>",
	Short: "Move a device to a different zone",
	Long: `Move a device to a different zone.

The zone can be specified by name or ID.

Examples:
  homeyctl devices move "Living Room Light" "Kitchen"
  homeyctl devices move "Sensor" "Bedroom"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		zone, err := findZone(args[1])
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"zone": zone.ID,
		}

		if err := apiClient.UpdateDevice(device.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Moved device '%s' to zone '%s'\n", device.Name, zone.Name)
		return nil
	},
}

var devicesSetNoteCmd = &cobra.Command{
	Use:   "set-note <device> <note>",
	Short: "Set a note on a device",
	Long: `Set or update the note/description for a device.

Use empty string "" to clear the note.

Examples:
  homeyctl devices set-note "Living Room Light" "Main ceiling light"
  homeyctl devices set-note "Sensor" ""`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"note": args[1],
		}

		if err := apiClient.UpdateDevice(device.ID, updates); err != nil {
			return err
		}

		if args[1] == "" {
			fmt.Printf("Cleared note for device '%s'\n", device.Name)
		} else {
			fmt.Printf("Set note for device '%s': %s\n", device.Name, args[1])
		}
		return nil
	},
}

var devicesSetIconCmd = &cobra.Command{
	Use:   "set-icon <device> <icon>",
	Short: "Set a custom icon for a device",
	Long: `Set a custom icon override for a device.

Use empty string "" to reset to the default icon.

Examples:
  homeyctl devices set-icon "My Device" "light"
  homeyctl devices set-icon "My Device" ""`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"iconOverride": args[1],
		}

		if err := apiClient.UpdateDevice(device.ID, updates); err != nil {
			return err
		}

		if args[1] == "" {
			fmt.Printf("Reset icon for device '%s' to default\n", device.Name)
		} else {
			fmt.Printf("Set icon for device '%s' to '%s'\n", device.Name, args[1])
		}
		return nil
	},
}

var devicesHideCmd = &cobra.Command{
	Use:   "hide <device>",
	Short: "Hide a device from the UI",
	Long: `Hide a device from the Homey UI.

The device will still function normally but won't appear in the device list.

Examples:
  homeyctl devices hide "Hidden Sensor"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"hidden": true,
		}

		if err := apiClient.UpdateDevice(device.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Hidden device '%s' from UI\n", device.Name)
		return nil
	},
}

var devicesUnhideCmd = &cobra.Command{
	Use:   "unhide <device>",
	Short: "Show a hidden device in the UI",
	Long: `Make a hidden device visible again in the Homey UI.

Examples:
  homeyctl devices unhide "Hidden Sensor"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"hidden": false,
		}

		if err := apiClient.UpdateDevice(device.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Unhidden device '%s' - now visible in UI\n", device.Name)
		return nil
	},
}

var devicesDeleteCmd = &cobra.Command{
	Use:   "delete <device>",
	Short: "Delete a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.DeleteDevice(device.ID); err != nil {
			return err
		}

		fmt.Printf("Deleted device: %s\n", device.Name)
		return nil
	},
}

func init() {
	devicesCmd.AddCommand(devicesRenameCmd)
	devicesCmd.AddCommand(devicesMoveCmd)
	devicesCmd.AddCommand(devicesSetNoteCmd)
	devicesCmd.AddCommand(devicesSetIconCmd)
	devicesCmd.AddCommand(devicesHideCmd)
	devicesCmd.AddCommand(devicesUnhideCmd)
	devicesCmd.AddCommand(devicesDeleteCmd)
}
