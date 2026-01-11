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
	zonesCmd.AddCommand(zonesDeleteCmd)
}
