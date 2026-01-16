package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type Zone struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Parent string `json:"parent"`
	Icon   string `json:"icon"`
}

// Known zone icons available in Homey
var knownZoneIcons = []string{
	"home",
	"livingRoom",
	"kitchen",
	"bedroom",
	"bedroomSingle",
	"bedroomDouble",
	"bedroomKids",
	"bathroom",
	"toilet",
	"office",
	"garage",
	"garden",
	"gardenShed",
	"basement",
	"attic",
	"hallway",
	"laundryRoom",
	"gameRoom",
	"diningRoom",
	"closet",
	"staircase",
	"balcony",
	"terrace",
	"pool",
	"gym",
	"sauna",
	"workshop",
	"storage",
	"groundFloor",
	"firstFloor",
	"secondFloor",
	"thirdFloor",
	"default",
}

var zoneRenameIcon string

var zonesCmd = &cobra.Command{
	Use:   "zones",
	Short: "Manage zones",
	Long:  `List, view, and delete Homey zones.`,
}

var zonesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all zones",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetZones()
		if err != nil {
			return err
		}

		if isTableFormat() {
			var zones map[string]Zone
			if err := json.Unmarshal(data, &zones); err != nil {
				return fmt.Errorf("failed to parse zones: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tICON\tID")
			fmt.Fprintln(w, "----\t----\t--")
			for _, z := range zones {
				fmt.Fprintf(w, "%s\t%s\t%s\n", z.Name, z.Icon, z.ID)
			}
			w.Flush()
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var zonesIconsCmd = &cobra.Command{
	Use:   "icons",
	Short: "List available zone icons",
	Long: `List all known zone icons that can be used with the --icon flag.

Note: This list may not be exhaustive. Homey may support additional icons.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if isTableFormat() {
			fmt.Println("Available zone icons:")
			fmt.Println()
			for _, icon := range knownZoneIcons {
				fmt.Printf("  %s\n", icon)
			}
			fmt.Println()
			fmt.Println("Use: homeyctl zones rename <zone> <new-name> --icon <icon>")
			fmt.Println("Or:  homeyctl zones set-icon <zone> <icon>")
			return nil
		}

		jsonData, _ := json.MarshalIndent(knownZoneIcons, "", "  ")
		fmt.Println(string(jsonData))
		return nil
	},
}

var zonesSetIconCmd = &cobra.Command{
	Use:   "set-icon <name-or-id> <icon>",
	Short: "Set the icon for a zone",
	Long: `Set the icon for a zone without changing its name.

Use 'homeyctl zones icons' to see available icons.

Examples:
  homeyctl zones set-icon "Aksels rom" bedroomSingle
  homeyctl zones set-icon "Garden" garden`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]
		newIcon := args[1]

		// Get all zones to find by name
		data, err := apiClient.GetZones()
		if err != nil {
			return err
		}

		var zones map[string]Zone
		if err := json.Unmarshal(data, &zones); err != nil {
			return fmt.Errorf("failed to parse zones: %w", err)
		}

		// Find zone by name or ID
		var zone *Zone
		for _, z := range zones {
			if z.ID == nameOrID || strings.EqualFold(z.Name, nameOrID) {
				zone = &z
				break
			}
		}

		if zone == nil {
			return fmt.Errorf("zone not found: %s", nameOrID)
		}

		updates := map[string]interface{}{
			"icon": newIcon,
		}

		if err := apiClient.UpdateZone(zone.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Changed icon for zone '%s' from '%s' to '%s'\n", zone.Name, zone.Icon, newIcon)
		return nil
	},
}

var zonesRenameCmd = &cobra.Command{
	Use:   "rename <name-or-id> <new-name>",
	Short: "Rename a zone",
	Long: `Rename a zone, optionally changing its icon.

Use 'homeyctl zones icons' to see available icons.

Examples:
  homeyctl zones rename "Office" "Home Office"
  homeyctl zones rename "Office" "Aksels rom" --icon bedroomSingle
  homeyctl zones rename abc123-zone-id "New Name"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]
		newName := args[1]

		// Get all zones to find by name
		data, err := apiClient.GetZones()
		if err != nil {
			return err
		}

		var zones map[string]Zone
		if err := json.Unmarshal(data, &zones); err != nil {
			return fmt.Errorf("failed to parse zones: %w", err)
		}

		// Find zone by name or ID
		var zone *Zone
		for _, z := range zones {
			if z.ID == nameOrID || strings.EqualFold(z.Name, nameOrID) {
				zone = &z
				break
			}
		}

		if zone == nil {
			return fmt.Errorf("zone not found: %s", nameOrID)
		}

		updates := map[string]interface{}{
			"name": newName,
		}

		if zoneRenameIcon != "" {
			updates["icon"] = zoneRenameIcon
		}

		if err := apiClient.UpdateZone(zone.ID, updates); err != nil {
			return err
		}

		if zoneRenameIcon != "" {
			fmt.Printf("Renamed zone '%s' to '%s' with icon '%s'\n", zone.Name, newName, zoneRenameIcon)
		} else {
			fmt.Printf("Renamed zone '%s' to '%s'\n", zone.Name, newName)
		}
		return nil
	},
}

var zonesDeleteCmd = &cobra.Command{
	Use:   "delete <name-or-id>",
	Short: "Delete a zone",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]

		// Get all zones to find by name
		data, err := apiClient.GetZones()
		if err != nil {
			return err
		}

		var zones map[string]Zone
		if err := json.Unmarshal(data, &zones); err != nil {
			return fmt.Errorf("failed to parse zones: %w", err)
		}

		// Find zone by name or ID
		for _, z := range zones {
			if z.ID == nameOrID || strings.EqualFold(z.Name, nameOrID) {
				if err := apiClient.DeleteZone(z.ID); err != nil {
					return err
				}
				fmt.Printf("Deleted zone: %s\n", z.Name)
				return nil
			}
		}

		return fmt.Errorf("zone not found: %s", nameOrID)
	},
}

func init() {
	rootCmd.AddCommand(zonesCmd)
	zonesCmd.AddCommand(zonesListCmd)
	zonesCmd.AddCommand(zonesIconsCmd)
	zonesCmd.AddCommand(zonesSetIconCmd)
	zonesCmd.AddCommand(zonesRenameCmd)
	zonesRenameCmd.Flags().StringVar(&zoneRenameIcon, "icon", "", "Change the zone icon (use 'zones icons' to see available icons)")
	zonesCmd.AddCommand(zonesDeleteCmd)
}
