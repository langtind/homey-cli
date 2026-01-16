package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// findUser finds a user by name or ID from the list of all users
func findUser(nameOrID string) (*User, error) {
	data, err := apiClient.GetUsers()
	if err != nil {
		return nil, err
	}

	var users map[string]User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users: %w", err)
	}

	for _, u := range users {
		if u.ID == nameOrID || strings.EqualFold(u.Name, nameOrID) {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("user not found: %s", nameOrID)
}

var usersGetCmd = &cobra.Command{
	Use:   "get <name-or-id>",
	Short: "Get user details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := findUser(args[0])
		if err != nil {
			return err
		}

		data, err := apiClient.GetUser(user.ID)
		if err != nil {
			return err
		}

		if isTableFormat() {
			var u struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Role    string `json:"role"`
				Enabled bool   `json:"enabled"`
				AthomID string `json:"athomId"`
			}
			if err := json.Unmarshal(data, &u); err != nil {
				return err
			}

			fmt.Printf("Name:    %s\n", u.Name)
			fmt.Printf("ID:      %s\n", u.ID)
			fmt.Printf("Role:    %s\n", u.Role)
			fmt.Printf("Enabled: %v\n", u.Enabled)
			fmt.Printf("AthomID: %s\n", u.AthomID)
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var usersMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Get current user details",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetUserMe()
		if err != nil {
			return err
		}

		if isTableFormat() {
			var u struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Role    string `json:"role"`
				Enabled bool   `json:"enabled"`
			}
			if err := json.Unmarshal(data, &u); err != nil {
				return err
			}

			fmt.Printf("Name:    %s\n", u.Name)
			fmt.Printf("ID:      %s\n", u.ID)
			fmt.Printf("Role:    %s\n", u.Role)
			fmt.Printf("Enabled: %v\n", u.Enabled)
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var usersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long: `Create a new user with a specific role.

Available roles:
  guest   - Can control devices but limited access
  user    - Normal user with standard access
  manager - Extended permissions, can manage users

Examples:
  homeyctl users create --role guest
  homeyctl users create --role user`,
	RunE: func(cmd *cobra.Command, args []string) error {
		role, _ := cmd.Flags().GetString("role")
		if role == "" {
			return fmt.Errorf("--role is required (guest, user, manager)")
		}

		validRoles := []string{"guest", "user", "manager"}
		isValid := false
		for _, r := range validRoles {
			if r == role {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid role: %s (valid: guest, user, manager)", role)
		}

		user := map[string]interface{}{
			"role": role,
		}

		result, err := apiClient.CreateUser(user)
		if err != nil {
			return err
		}

		if isTableFormat() {
			fmt.Printf("Created user with role: %s\n", role)
			return nil
		}

		outputJSON(result)
		return nil
	},
}

var usersUpdateCmd = &cobra.Command{
	Use:   "update <name-or-id>",
	Short: "Update a user",
	Long: `Update a user's role or enabled status.

Examples:
  homeyctl users update "John" --role manager
  homeyctl users update "Guest" --enabled=false`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := findUser(args[0])
		if err != nil {
			return err
		}

		updates := map[string]interface{}{}

		if role, _ := cmd.Flags().GetString("role"); role != "" {
			updates["role"] = role
		}

		if cmd.Flags().Changed("enabled") {
			enabled, _ := cmd.Flags().GetBool("enabled")
			updates["enabled"] = enabled
		}

		if len(updates) == 0 {
			return fmt.Errorf("no updates specified (use --role or --enabled)")
		}

		if err := apiClient.UpdateUser(user.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Updated user: %s\n", user.Name)
		return nil
	},
}

var usersDeleteCmd = &cobra.Command{
	Use:   "delete <name-or-id>",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := findUser(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.DeleteUser(user.ID); err != nil {
			return err
		}

		fmt.Printf("Deleted user: %s\n", user.Name)
		return nil
	},
}

var usersPresenceCmd = &cobra.Command{
	Use:   "presence",
	Short: "Show user presence status",
	RunE: func(cmd *cobra.Command, args []string) error {
		usersData, err := apiClient.GetUsers()
		if err != nil {
			return err
		}

		var users map[string]User
		if err := json.Unmarshal(usersData, &users); err != nil {
			return fmt.Errorf("failed to parse users: %w", err)
		}

		if isTableFormat() {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tPRESENT\tASLEEP\tID")
			fmt.Fprintln(w, "----\t-------\t------\t--")

			for _, u := range users {
				// Get presence status
				presentData, _ := apiClient.GetPresent(u.ID)
				asleepData, _ := apiClient.GetAsleep(u.ID)

				var present, asleep struct {
					Value bool `json:"value"`
				}
				json.Unmarshal(presentData, &present)
				json.Unmarshal(asleepData, &asleep)

				presentStr := "no"
				if present.Value {
					presentStr = "yes"
				}
				asleepStr := "no"
				if asleep.Value {
					asleepStr = "yes"
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", u.Name, presentStr, asleepStr, u.ID)
			}
			w.Flush()
			return nil
		}

		// JSON output - build presence map
		presenceMap := make(map[string]interface{})
		for _, u := range users {
			presentData, _ := apiClient.GetPresent(u.ID)
			asleepData, _ := apiClient.GetAsleep(u.ID)

			var present, asleep struct {
				Value bool `json:"value"`
			}
			json.Unmarshal(presentData, &present)
			json.Unmarshal(asleepData, &asleep)

			presenceMap[u.ID] = map[string]interface{}{
				"name":    u.Name,
				"present": present.Value,
				"asleep":  asleep.Value,
			}
		}

		out, _ := json.MarshalIndent(presenceMap, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	usersCmd.AddCommand(usersGetCmd)
	usersCmd.AddCommand(usersMeCmd)

	usersCmd.AddCommand(usersCreateCmd)
	usersCreateCmd.Flags().String("role", "", "User role (guest, user, manager)")

	usersCmd.AddCommand(usersUpdateCmd)
	usersUpdateCmd.Flags().String("role", "", "New role (guest, user, manager, owner)")
	usersUpdateCmd.Flags().Bool("enabled", true, "Enable/disable user")

	usersCmd.AddCommand(usersDeleteCmd)
	usersCmd.AddCommand(usersPresenceCmd)
}
