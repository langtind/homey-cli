package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// FlowFolder represents a Homey flow folder
type FlowFolder struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

var flowsFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "Manage flow folders",
	Long:  `List, create, update, and delete flow folders.`,
}

// findFlowFolder finds a flow folder by name or ID
func findFlowFolder(nameOrID string) (*FlowFolder, error) {
	data, err := apiClient.GetFlowFolders()
	if err != nil {
		return nil, err
	}

	var folders map[string]FlowFolder
	if err := json.Unmarshal(data, &folders); err != nil {
		return nil, fmt.Errorf("failed to parse flow folders: %w", err)
	}

	for _, f := range folders {
		if f.ID == nameOrID || strings.EqualFold(f.Name, nameOrID) {
			return &f, nil
		}
	}

	return nil, fmt.Errorf("flow folder not found: %s", nameOrID)
}

var flowsFoldersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all flow folders",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetFlowFolders()
		if err != nil {
			return err
		}

		if isTableFormat() {
			var folders map[string]FlowFolder
			if err := json.Unmarshal(data, &folders); err != nil {
				return fmt.Errorf("failed to parse flow folders: %w", err)
			}

			if len(folders) == 0 {
				fmt.Println("No flow folders found.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tPARENT\tID")
			fmt.Fprintln(w, "----\t------\t--")
			for _, f := range folders {
				parent := f.Parent
				if parent == "" {
					parent = "(root)"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", f.Name, parent, f.ID)
			}
			w.Flush()
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var flowsFoldersGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get flow folder details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		folder, err := findFlowFolder(args[0])
		if err != nil {
			return err
		}

		data, err := apiClient.GetFlowFolder(folder.ID)
		if err != nil {
			return err
		}

		outputJSON(data)
		return nil
	},
}

var flowsFoldersCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new flow folder",
	Long: `Create a new flow folder.

Examples:
  homeyctl flows folders create "Lighting"
  homeyctl flows folders create "Kitchen Automations" --parent "Lighting"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		parentName, _ := cmd.Flags().GetString("parent")

		folder := map[string]interface{}{
			"name": name,
		}

		if parentName != "" {
			parent, err := findFlowFolder(parentName)
			if err != nil {
				return fmt.Errorf("parent folder: %w", err)
			}
			folder["parent"] = parent.ID
		}

		result, err := apiClient.CreateFlowFolder(folder)
		if err != nil {
			return err
		}

		if isTableFormat() {
			fmt.Printf("Created flow folder: %s\n", name)
			return nil
		}

		outputJSON(result)
		return nil
	},
}

var flowsFoldersUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a flow folder",
	Long: `Update a flow folder's name or parent.

Examples:
  homeyctl flows folders update "Lighting" --name "Light Automations"
  homeyctl flows folders update "Kitchen" --parent "Lighting"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		folder, err := findFlowFolder(args[0])
		if err != nil {
			return err
		}

		updates := map[string]interface{}{}

		if name, _ := cmd.Flags().GetString("name"); name != "" {
			updates["name"] = name
		}

		if parentName, _ := cmd.Flags().GetString("parent"); parentName != "" {
			parent, err := findFlowFolder(parentName)
			if err != nil {
				return fmt.Errorf("parent folder: %w", err)
			}
			updates["parent"] = parent.ID
		}

		if len(updates) == 0 {
			return fmt.Errorf("no updates specified (use --name or --parent)")
		}

		if err := apiClient.UpdateFlowFolder(folder.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Updated flow folder: %s\n", folder.Name)
		return nil
	},
}

var flowsFoldersDeleteCmd = &cobra.Command{
	Use:   "delete <name-or-id>",
	Short: "Delete a flow folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		folder, err := findFlowFolder(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.DeleteFlowFolder(folder.ID); err != nil {
			return err
		}

		fmt.Printf("Deleted flow folder: %s\n", folder.Name)
		return nil
	},
}

func init() {
	flowsCmd.AddCommand(flowsFoldersCmd)

	flowsFoldersCmd.AddCommand(flowsFoldersListCmd)
	flowsFoldersCmd.AddCommand(flowsFoldersGetCmd)

	flowsFoldersCmd.AddCommand(flowsFoldersCreateCmd)
	flowsFoldersCreateCmd.Flags().String("parent", "", "Parent folder name or ID")

	flowsFoldersCmd.AddCommand(flowsFoldersUpdateCmd)
	flowsFoldersUpdateCmd.Flags().String("name", "", "New name for the folder")
	flowsFoldersUpdateCmd.Flags().String("parent", "", "New parent folder name or ID")

	flowsFoldersCmd.AddCommand(flowsFoldersDeleteCmd)
}
