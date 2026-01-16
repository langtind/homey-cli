package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var presenceCmd = &cobra.Command{
	Use:   "presence",
	Short: "Manage user presence",
	Long:  `Get and set user presence (home/away) and sleep status.`,
}

var presenceGetCmd = &cobra.Command{
	Use:   "get <user>",
	Short: "Get user presence status",
	Long: `Get presence status for a specific user.

Examples:
  homeyctl presence get "Arild"
  homeyctl presence get me`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]

		var userID string
		var userName string

		if nameOrID == "me" {
			// Get current user
			data, err := apiClient.GetUserMe()
			if err != nil {
				return err
			}
			var u struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			json.Unmarshal(data, &u)
			userID = u.ID
			userName = u.Name
		} else {
			user, err := findUser(nameOrID)
			if err != nil {
				return err
			}
			userID = user.ID
			userName = user.Name
		}

		presentData, err := apiClient.GetPresent(userID)
		if err != nil {
			return err
		}

		asleepData, err := apiClient.GetAsleep(userID)
		if err != nil {
			return err
		}

		var present, asleep struct {
			Value bool `json:"value"`
		}
		json.Unmarshal(presentData, &present)
		json.Unmarshal(asleepData, &asleep)

		if isTableFormat() {
			presentStr := "away"
			if present.Value {
				presentStr = "home"
			}
			asleepStr := "awake"
			if asleep.Value {
				asleepStr = "asleep"
			}

			fmt.Printf("User:    %s\n", userName)
			fmt.Printf("Present: %s\n", presentStr)
			fmt.Printf("Asleep:  %s\n", asleepStr)
			return nil
		}

		result := map[string]interface{}{
			"user":    userName,
			"present": present.Value,
			"asleep":  asleep.Value,
		}
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

var presenceSetCmd = &cobra.Command{
	Use:   "set <user> <home|away>",
	Short: "Set user presence status",
	Long: `Set presence status for a specific user.

Use "me" to set your own presence.

Examples:
  homeyctl presence set me home
  homeyctl presence set me away
  homeyctl presence set "Arild" home`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]
		status := args[1]

		var present bool
		switch status {
		case "home", "true", "1", "yes":
			present = true
		case "away", "false", "0", "no":
			present = false
		default:
			return fmt.Errorf("invalid status: %s (use: home, away)", status)
		}

		var err error
		if nameOrID == "me" {
			err = apiClient.SetPresentMe(present)
		} else {
			user, findErr := findUser(nameOrID)
			if findErr != nil {
				return findErr
			}
			err = apiClient.SetPresent(user.ID, present)
		}

		if err != nil {
			return err
		}

		statusStr := "away"
		if present {
			statusStr = "home"
		}
		fmt.Printf("Set %s presence to: %s\n", nameOrID, statusStr)
		return nil
	},
}

var presenceAsleepCmd = &cobra.Command{
	Use:   "asleep",
	Short: "Manage sleep status",
}

var presenceAsleepGetCmd = &cobra.Command{
	Use:   "get <user>",
	Short: "Get user sleep status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]

		var userID string
		var userName string

		if nameOrID == "me" {
			data, err := apiClient.GetUserMe()
			if err != nil {
				return err
			}
			var u struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			json.Unmarshal(data, &u)
			userID = u.ID
			userName = u.Name
		} else {
			user, err := findUser(nameOrID)
			if err != nil {
				return err
			}
			userID = user.ID
			userName = user.Name
		}

		data, err := apiClient.GetAsleep(userID)
		if err != nil {
			return err
		}

		var asleep struct {
			Value bool `json:"value"`
		}
		json.Unmarshal(data, &asleep)

		if isTableFormat() {
			status := "awake"
			if asleep.Value {
				status = "asleep"
			}
			fmt.Printf("%s is %s\n", userName, status)
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var presenceAsleepSetCmd = &cobra.Command{
	Use:   "set <user> <asleep|awake>",
	Short: "Set user sleep status",
	Long: `Set sleep status for a specific user.

Use "me" to set your own status.

Examples:
  homeyctl presence asleep set me asleep
  homeyctl presence asleep set me awake
  homeyctl presence asleep set "Arild" asleep`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]
		status := args[1]

		var asleep bool
		switch status {
		case "asleep", "true", "1", "yes", "sleeping":
			asleep = true
		case "awake", "false", "0", "no", "woke":
			asleep = false
		default:
			return fmt.Errorf("invalid status: %s (use: asleep, awake)", status)
		}

		var err error
		if nameOrID == "me" {
			err = apiClient.SetAsleepMe(asleep)
		} else {
			user, findErr := findUser(nameOrID)
			if findErr != nil {
				return findErr
			}
			err = apiClient.SetAsleep(user.ID, asleep)
		}

		if err != nil {
			return err
		}

		statusStr := "awake"
		if asleep {
			statusStr = "asleep"
		}
		fmt.Printf("Set %s sleep status to: %s\n", nameOrID, statusStr)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(presenceCmd)

	presenceCmd.AddCommand(presenceGetCmd)
	presenceCmd.AddCommand(presenceSetCmd)

	presenceCmd.AddCommand(presenceAsleepCmd)
	presenceAsleepCmd.AddCommand(presenceAsleepGetCmd)
	presenceAsleepCmd.AddCommand(presenceAsleepSetCmd)
}
