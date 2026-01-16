package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// Notification represents a Homey notification
type Notification struct {
	ID       string `json:"id"`
	Excerpt  string `json:"excerpt"`
	OwnerUri string `json:"ownerUri"`
	Date     string `json:"date"`
}

var notifyCmd = &cobra.Command{
	Use:     "notify",
	Aliases: []string{"notifications"},
	Short:   "Manage notifications",
	Long:    `Send and view Homey timeline notifications.`,
}

var notifySendCmd = &cobra.Command{
	Use:   "send <message>",
	Short: "Send a notification to the timeline",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := args[0]

		if err := apiClient.SendNotification(message); err != nil {
			return err
		}

		fmt.Printf("Notification sent: %s\n", message)
		return nil
	},
}

var notifyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List timeline notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetNotifications()
		if err != nil {
			return err
		}

		if isTableFormat() {
			var notifications map[string]Notification
			if err := json.Unmarshal(data, &notifications); err != nil {
				return fmt.Errorf("failed to parse notifications: %w", err)
			}

			if len(notifications) == 0 {
				fmt.Println("No notifications found.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "DATE\tMESSAGE\tID")
			fmt.Fprintln(w, "----\t-------\t--")
			for _, n := range notifications {
				excerpt := n.Excerpt
				if len(excerpt) > 50 {
					excerpt = excerpt[:47] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", n.Date, excerpt, n.ID)
			}
			w.Flush()
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var notifyDeleteCmd = &cobra.Command{
	Use:   "delete <notification-id>",
	Short: "Delete a notification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.DeleteNotification(args[0]); err != nil {
			return err
		}

		fmt.Printf("Deleted notification: %s\n", args[0])
		return nil
	},
}

var notifyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.DeleteAllNotifications(); err != nil {
			return err
		}

		fmt.Println("Cleared all notifications")
		return nil
	},
}

var notifyOwnersCmd = &cobra.Command{
	Use:   "owners",
	Short: "List notification owners",
	Long:  `List all notification sources/owners and their settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetNotificationOwners()
		if err != nil {
			return err
		}

		outputJSON(data)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)
	notifyCmd.AddCommand(notifySendCmd)
	notifyCmd.AddCommand(notifyListCmd)
	notifyCmd.AddCommand(notifyDeleteCmd)
	notifyCmd.AddCommand(notifyClearCmd)
	notifyCmd.AddCommand(notifyOwnersCmd)
}
