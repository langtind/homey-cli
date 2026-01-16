package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// Dashboard represents a Homey dashboard
type Dashboard struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Columns []interface{} `json:"columns"`
}

var dashboardsCmd = &cobra.Command{
	Use:   "dashboards",
	Short: "Manage dashboards",
	Long:  `List, create, update, and delete dashboards.`,
}

// findDashboard finds a dashboard by name or ID
func findDashboard(nameOrID string) (*Dashboard, error) {
	data, err := apiClient.GetDashboards()
	if err != nil {
		return nil, err
	}

	var dashboards map[string]Dashboard
	if err := json.Unmarshal(data, &dashboards); err != nil {
		return nil, fmt.Errorf("failed to parse dashboards: %w", err)
	}

	for _, d := range dashboards {
		if d.ID == nameOrID || strings.EqualFold(d.Name, nameOrID) {
			return &d, nil
		}
	}

	return nil, fmt.Errorf("dashboard not found: %s", nameOrID)
}

var dashboardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all dashboards",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetDashboards()
		if err != nil {
			return err
		}

		if isTableFormat() {
			var dashboards map[string]Dashboard
			if err := json.Unmarshal(data, &dashboards); err != nil {
				return fmt.Errorf("failed to parse dashboards: %w", err)
			}

			if len(dashboards) == 0 {
				fmt.Println("No dashboards found.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tCOLUMNS\tID")
			fmt.Fprintln(w, "----\t-------\t--")
			for _, d := range dashboards {
				fmt.Fprintf(w, "%s\t%d\t%s\n", d.Name, len(d.Columns), d.ID)
			}
			w.Flush()
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var dashboardsGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get dashboard details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dashboard, err := findDashboard(args[0])
		if err != nil {
			return err
		}

		data, err := apiClient.GetDashboard(dashboard.ID)
		if err != nil {
			return err
		}

		outputJSON(data)
		return nil
	},
}

var dashboardsCreateCmd = &cobra.Command{
	Use:   "create <name> [json-file]",
	Short: "Create a new dashboard",
	Long: `Create a new dashboard.

The optional JSON file should contain the dashboard structure:
{
  "name": "My Dashboard",
  "columns": [
    {
      "widgets": [...]
    }
  ]
}

Examples:
  homeyctl dashboards create "My Dashboard"
  homeyctl dashboards create "My Dashboard" dashboard.json`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		var dashboard map[string]interface{}

		if len(args) == 2 {
			data, err := os.ReadFile(args[1])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			if err := json.Unmarshal(data, &dashboard); err != nil {
				return fmt.Errorf("invalid JSON: %w", err)
			}
		} else {
			dashboard = map[string]interface{}{
				"columns": []interface{}{},
			}
		}

		dashboard["name"] = name

		result, err := apiClient.CreateDashboard(dashboard)
		if err != nil {
			return err
		}

		if isTableFormat() {
			fmt.Printf("Created dashboard: %s\n", name)
			return nil
		}

		outputJSON(result)
		return nil
	},
}

var dashboardsUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id> <json-file>",
	Short: "Update a dashboard",
	Long: `Update a dashboard from a JSON file.

Examples:
  homeyctl dashboards get "My Dashboard" > dashboard.json
  # Edit dashboard.json
  homeyctl dashboards update "My Dashboard" dashboard.json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dashboard, err := findDashboard(args[0])
		if err != nil {
			return err
		}

		data, err := os.ReadFile(args[1])
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		var updates map[string]interface{}
		if err := json.Unmarshal(data, &updates); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		if err := apiClient.UpdateDashboard(dashboard.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Updated dashboard: %s\n", dashboard.Name)
		return nil
	},
}

var dashboardsDeleteCmd = &cobra.Command{
	Use:   "delete <name-or-id>",
	Short: "Delete a dashboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dashboard, err := findDashboard(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.DeleteDashboard(dashboard.ID); err != nil {
			return err
		}

		fmt.Printf("Deleted dashboard: %s\n", dashboard.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dashboardsCmd)

	dashboardsCmd.AddCommand(dashboardsListCmd)
	dashboardsCmd.AddCommand(dashboardsGetCmd)
	dashboardsCmd.AddCommand(dashboardsCreateCmd)
	dashboardsCmd.AddCommand(dashboardsUpdateCmd)
	dashboardsCmd.AddCommand(dashboardsDeleteCmd)
}
