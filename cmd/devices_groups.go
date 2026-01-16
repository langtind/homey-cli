package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// DeviceGroup represents a Homey device group (virtual device that controls multiple devices)
type DeviceGroup struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Class        string   `json:"class"`
	VirtualClass string   `json:"virtualClass"`
	Zone         string   `json:"zone"`
	Devices      []string `json:"devices"`
}

var devicesGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage device groups",
	Long:  `List, create, update, and manage device groups.`,
}

var devicesGroupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all device groups",
	Long: `List all device groups.

Device groups are virtual devices that control multiple physical devices together.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetDevices()
		if err != nil {
			return err
		}

		var devices map[string]json.RawMessage
		if err := json.Unmarshal(data, &devices); err != nil {
			return fmt.Errorf("failed to parse devices: %w", err)
		}

		// Filter for groups (devices with virtualClass "group")
		var groups []DeviceGroup
		for _, raw := range devices {
			var d struct {
				ID           string   `json:"id"`
				Name         string   `json:"name"`
				Class        string   `json:"class"`
				VirtualClass string   `json:"virtualClass"`
				Zone         string   `json:"zone"`
				Devices      []string `json:"devices"`
			}
			if err := json.Unmarshal(raw, &d); err != nil {
				continue
			}
			if d.VirtualClass == "group" {
				groups = append(groups, DeviceGroup{
					ID:           d.ID,
					Name:         d.Name,
					Class:        d.Class,
					VirtualClass: d.VirtualClass,
					Zone:         d.Zone,
					Devices:      d.Devices,
				})
			}
		}

		if isTableFormat() {
			if len(groups) == 0 {
				fmt.Println("No device groups found.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tCLASS\tDEVICES\tID")
			fmt.Fprintln(w, "----\t-----\t-------\t--")
			for _, g := range groups {
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", g.Name, g.Class, len(g.Devices), g.ID)
			}
			w.Flush()
			return nil
		}

		out, _ := json.MarshalIndent(groups, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

var groupCreateClass string
var groupCreateZone string
var groupCreateDevices string

var devicesGroupsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a device group",
	Long: `Create a new device group.

A device group is a virtual device that controls multiple physical devices together.
All devices in the group must be of the same class (e.g., all lights).

Examples:
  homeyctl devices groups create "Living Room Lights" --class light --zone "Living Room" --devices "Light 1,Light 2,Light 3"
  homeyctl devices groups create "All Fans" --class fan --zone "Home" --devices "Fan 1,Fan 2"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if groupCreateClass == "" {
			return fmt.Errorf("--class is required")
		}
		if groupCreateZone == "" {
			return fmt.Errorf("--zone is required")
		}
		if groupCreateDevices == "" {
			return fmt.Errorf("--devices is required")
		}

		// Find zone ID
		zone, err := findZone(groupCreateZone)
		if err != nil {
			return err
		}

		// Find device IDs
		deviceNames := strings.Split(groupCreateDevices, ",")
		var deviceIDs []string
		for _, deviceName := range deviceNames {
			deviceName = strings.TrimSpace(deviceName)
			if deviceName == "" {
				continue
			}
			device, err := findDevice(deviceName)
			if err != nil {
				return fmt.Errorf("device '%s': %w", deviceName, err)
			}
			deviceIDs = append(deviceIDs, device.ID)
		}

		if len(deviceIDs) == 0 {
			return fmt.Errorf("at least one device is required")
		}

		group := map[string]interface{}{
			"name":      name,
			"class":     groupCreateClass,
			"zoneId":    zone.ID,
			"deviceIds": deviceIDs,
		}

		result, err := apiClient.CreateDeviceGroup(group)
		if err != nil {
			return err
		}

		if isTableFormat() {
			fmt.Printf("Created device group '%s' with %d devices\n", name, len(deviceIDs))
			return nil
		}

		outputJSON(result)
		return nil
	},
}

var groupUpdateName string
var groupUpdateAddDevices string
var groupUpdateRemoveDevices string

var devicesGroupsUpdateCmd = &cobra.Command{
	Use:   "update <group>",
	Short: "Update a device group",
	Long: `Update a device group's name or devices.

Examples:
  homeyctl devices groups update "Living Room Lights" --name "LR Lights"
  homeyctl devices groups update "Living Room Lights" --add "New Light"
  homeyctl devices groups update "Living Room Lights" --remove "Old Light"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Find the group (it's a device with virtualClass "group")
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		updates := map[string]interface{}{}

		if groupUpdateName != "" {
			updates["name"] = groupUpdateName
		}

		// For adding/removing devices, we need to get current devices and modify the list
		if groupUpdateAddDevices != "" || groupUpdateRemoveDevices != "" {
			// Get current device details to get the devices list
			data, err := apiClient.GetDevices()
			if err != nil {
				return err
			}

			var devices map[string]json.RawMessage
			if err := json.Unmarshal(data, &devices); err != nil {
				return fmt.Errorf("failed to parse devices: %w", err)
			}

			var groupData struct {
				Devices []string `json:"devices"`
			}
			if raw, ok := devices[device.ID]; ok {
				if err := json.Unmarshal(raw, &groupData); err != nil {
					return fmt.Errorf("failed to parse group: %w", err)
				}
			}

			deviceIDs := groupData.Devices

			// Add devices
			if groupUpdateAddDevices != "" {
				for _, name := range strings.Split(groupUpdateAddDevices, ",") {
					name = strings.TrimSpace(name)
					if name == "" {
						continue
					}
					d, err := findDevice(name)
					if err != nil {
						return fmt.Errorf("device '%s': %w", name, err)
					}
					// Check if already in group
					found := false
					for _, id := range deviceIDs {
						if id == d.ID {
							found = true
							break
						}
					}
					if !found {
						deviceIDs = append(deviceIDs, d.ID)
					}
				}
			}

			// Remove devices
			if groupUpdateRemoveDevices != "" {
				for _, name := range strings.Split(groupUpdateRemoveDevices, ",") {
					name = strings.TrimSpace(name)
					if name == "" {
						continue
					}
					d, err := findDevice(name)
					if err != nil {
						return fmt.Errorf("device '%s': %w", name, err)
					}
					// Remove from list
					var newDeviceIDs []string
					for _, id := range deviceIDs {
						if id != d.ID {
							newDeviceIDs = append(newDeviceIDs, id)
						}
					}
					deviceIDs = newDeviceIDs
				}
			}

			updates["deviceIds"] = deviceIDs
		}

		if len(updates) == 0 {
			return fmt.Errorf("no updates specified (use --name, --add, or --remove)")
		}

		if err := apiClient.UpdateDeviceGroup(device.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Updated device group '%s'\n", device.Name)
		return nil
	},
}

var devicesGroupsRemoveDeviceCmd = &cobra.Command{
	Use:   "remove-device <group> <device>",
	Short: "Remove a device from a group",
	Long: `Remove a device from a device group.

Examples:
  homeyctl devices groups remove-device "Living Room Lights" "Old Light"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		group, err := findDevice(args[0])
		if err != nil {
			return fmt.Errorf("group: %w", err)
		}

		device, err := findDevice(args[1])
		if err != nil {
			return fmt.Errorf("device: %w", err)
		}

		if err := apiClient.RemoveDeviceFromGroup(group.ID, device.ID); err != nil {
			return err
		}

		fmt.Printf("Removed '%s' from group '%s'\n", device.Name, group.Name)
		return nil
	},
}

func init() {
	devicesCmd.AddCommand(devicesGroupsCmd)

	devicesGroupsCmd.AddCommand(devicesGroupsListCmd)

	devicesGroupsCmd.AddCommand(devicesGroupsCreateCmd)
	devicesGroupsCreateCmd.Flags().StringVar(&groupCreateClass, "class", "", "Device class (e.g., light, socket, fan)")
	devicesGroupsCreateCmd.Flags().StringVar(&groupCreateZone, "zone", "", "Zone for the group")
	devicesGroupsCreateCmd.Flags().StringVar(&groupCreateDevices, "devices", "", "Comma-separated list of device names or IDs")

	devicesGroupsCmd.AddCommand(devicesGroupsUpdateCmd)
	devicesGroupsUpdateCmd.Flags().StringVar(&groupUpdateName, "name", "", "New name for the group")
	devicesGroupsUpdateCmd.Flags().StringVar(&groupUpdateAddDevices, "add", "", "Comma-separated list of devices to add")
	devicesGroupsUpdateCmd.Flags().StringVar(&groupUpdateRemoveDevices, "remove", "", "Comma-separated list of devices to remove")

	devicesGroupsCmd.AddCommand(devicesGroupsRemoveDeviceCmd)
}
