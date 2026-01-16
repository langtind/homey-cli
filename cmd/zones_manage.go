package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var zoneRenameIcon string

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

		zone, err := findZone(nameOrID)
		if err != nil {
			return err
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

		zone, err := findZone(nameOrID)
		if err != nil {
			return err
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

var zoneCreateIcon string

var zonesCreateCmd = &cobra.Command{
	Use:   "create <name> --parent <parent-zone>",
	Short: "Create a new zone",
	Long: `Create a new zone under a parent zone.

Use 'homeyctl zones icons' to see available icons.

Examples:
  homeyctl zones create "Office" --parent "Home"
  homeyctl zones create "Kids Room" --parent "Second Floor" --icon bedroomKids`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		parentName, _ := cmd.Flags().GetString("parent")
		if parentName == "" {
			return fmt.Errorf("--parent is required")
		}

		parent, err := findZone(parentName)
		if err != nil {
			return fmt.Errorf("parent zone: %w", err)
		}

		icon := zoneCreateIcon
		if icon == "" {
			icon = "default"
		}

		zone := map[string]interface{}{
			"name":   name,
			"parent": parent.ID,
			"icon":   icon,
		}

		result, err := apiClient.CreateZone(zone)
		if err != nil {
			return err
		}

		if isTableFormat() {
			fmt.Printf("Created zone '%s' under '%s'\n", name, parent.Name)
			return nil
		}

		outputJSON(result)
		return nil
	},
}

var zonesMoveCmd = &cobra.Command{
	Use:   "move <zone> <new-parent>",
	Short: "Move a zone to a different parent",
	Long: `Move a zone to be under a different parent zone.

Examples:
  homeyctl zones move "Office" "Second Floor"
  homeyctl zones move "Kids Room" "Home"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, err := findZone(args[0])
		if err != nil {
			return err
		}

		newParent, err := findZone(args[1])
		if err != nil {
			return fmt.Errorf("new parent zone: %w", err)
		}

		updates := map[string]interface{}{
			"parent": newParent.ID,
		}

		if err := apiClient.UpdateZone(zone.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Moved zone '%s' to parent '%s'\n", zone.Name, newParent.Name)
		return nil
	},
}

var zonesDeleteCmd = &cobra.Command{
	Use:   "delete <zone>",
	Short: "Delete a zone",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, err := findZone(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.DeleteZone(zone.ID); err != nil {
			return err
		}

		fmt.Printf("Deleted zone: %s\n", zone.Name)
		return nil
	},
}

func init() {
	zonesCmd.AddCommand(zonesCreateCmd)
	zonesCreateCmd.Flags().String("parent", "", "Parent zone (required)")
	zonesCreateCmd.Flags().StringVar(&zoneCreateIcon, "icon", "", "Zone icon (use 'zones icons' to see available)")

	zonesCmd.AddCommand(zonesRenameCmd)
	zonesRenameCmd.Flags().StringVar(&zoneRenameIcon, "icon", "", "Change the zone icon (use 'zones icons' to see available icons)")

	zonesCmd.AddCommand(zonesMoveCmd)
	zonesCmd.AddCommand(zonesSetIconCmd)
	zonesCmd.AddCommand(zonesDeleteCmd)
}
