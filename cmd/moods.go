package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// Mood represents a Homey mood
type Mood struct {
	ID     string                 `json:"id"`
	Name   string                 `json:"name"`
	Preset string                 `json:"preset"`
	Zone   string                 `json:"zone"`
	Active bool                   `json:"active"`
	Devices map[string]interface{} `json:"devices"`
}

var moodsCmd = &cobra.Command{
	Use:   "moods",
	Short: "Manage moods",
	Long:  `List, create, update, delete, and activate moods.`,
}

// findMood finds a mood by name or ID
func findMood(nameOrID string) (*Mood, error) {
	data, err := apiClient.GetMoods()
	if err != nil {
		return nil, err
	}

	var moods map[string]Mood
	if err := json.Unmarshal(data, &moods); err != nil {
		return nil, fmt.Errorf("failed to parse moods: %w", err)
	}

	for _, m := range moods {
		if m.ID == nameOrID || strings.EqualFold(m.Name, nameOrID) {
			return &m, nil
		}
	}

	return nil, fmt.Errorf("mood not found: %s", nameOrID)
}

var moodsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all moods",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetMoods()
		if err != nil {
			return err
		}

		if isTableFormat() {
			var moods map[string]Mood
			if err := json.Unmarshal(data, &moods); err != nil {
				return fmt.Errorf("failed to parse moods: %w", err)
			}

			if len(moods) == 0 {
				fmt.Println("No moods found.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tPRESET\tACTIVE\tID")
			fmt.Fprintln(w, "----\t------\t------\t--")
			for _, m := range moods {
				active := "no"
				if m.Active {
					active = "yes"
				}
				preset := m.Preset
				if preset == "" {
					preset = "-"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", m.Name, preset, active, m.ID)
			}
			w.Flush()
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var moodsGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get mood details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mood, err := findMood(args[0])
		if err != nil {
			return err
		}

		data, err := apiClient.GetMood(mood.ID)
		if err != nil {
			return err
		}

		outputJSON(data)
		return nil
	},
}

var moodsCreateCmd = &cobra.Command{
	Use:   "create <name> <json-file>",
	Short: "Create a new mood",
	Long: `Create a new mood from a JSON file.

The JSON file should contain device settings:
{
  "name": "Movie Night",
  "devices": {
    "<device-id>": {
      "onoff": false
    },
    "<device-id>": {
      "dim": 0.3
    }
  }
}

Examples:
  homeyctl moods create "Movie Night" mood.json
  cat mood.json | homeyctl moods create "Relax" -`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		var mood map[string]interface{}

		if len(args) == 2 {
			var data []byte
			var err error

			if args[1] == "-" {
				data, err = os.ReadFile("/dev/stdin")
			} else {
				data, err = os.ReadFile(args[1])
			}
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			if err := json.Unmarshal(data, &mood); err != nil {
				return fmt.Errorf("invalid JSON: %w", err)
			}
		} else {
			mood = map[string]interface{}{
				"devices": map[string]interface{}{},
			}
		}

		mood["name"] = name

		if _, ok := mood["devices"]; !ok {
			mood["devices"] = map[string]interface{}{}
		}

		result, err := apiClient.CreateMood(mood)
		if err != nil {
			return err
		}

		if isTableFormat() {
			fmt.Printf("Created mood: %s\n", name)
			return nil
		}

		outputJSON(result)
		return nil
	},
}

var moodsUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id> <json-file>",
	Short: "Update a mood",
	Long: `Update a mood from a JSON file.

Examples:
  homeyctl moods get "Movie Night" > mood.json
  # Edit mood.json
  homeyctl moods update "Movie Night" mood.json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mood, err := findMood(args[0])
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

		if err := apiClient.UpdateMood(mood.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Updated mood: %s\n", mood.Name)
		return nil
	},
}

var moodsDeleteCmd = &cobra.Command{
	Use:   "delete <name-or-id>",
	Short: "Delete a mood",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mood, err := findMood(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.DeleteMood(mood.ID); err != nil {
			return err
		}

		fmt.Printf("Deleted mood: %s\n", mood.Name)
		return nil
	},
}

var moodsSetCmd = &cobra.Command{
	Use:   "set <name-or-id>",
	Short: "Activate a mood",
	Long: `Activate a mood, applying its device settings.

Examples:
  homeyctl moods set "Movie Night"
  homeyctl moods set "Relax"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mood, err := findMood(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.SetMood(mood.ID); err != nil {
			return err
		}

		fmt.Printf("Activated mood: %s\n", mood.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(moodsCmd)

	moodsCmd.AddCommand(moodsListCmd)
	moodsCmd.AddCommand(moodsGetCmd)
	moodsCmd.AddCommand(moodsCreateCmd)
	moodsCmd.AddCommand(moodsUpdateCmd)
	moodsCmd.AddCommand(moodsDeleteCmd)
	moodsCmd.AddCommand(moodsSetCmd)
}
