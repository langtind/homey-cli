package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

type InsightLog struct {
	ID       string `json:"id"`
	OwnerURI string `json:"ownerUri"`
	OwnerID  string `json:"ownerId"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Units    string `json:"units"`
}

var insightsCmd = &cobra.Command{
	Use:   "insights",
	Short: "Manage insights",
	Long:  `View Homey insights logs and historical data.`,
}

var insightsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all insight logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetInsights()
		if err != nil {
			return err
		}

		if isTableFormat() {
			var logs []InsightLog
			if err := json.Unmarshal(data, &logs); err != nil {
				return fmt.Errorf("failed to parse insights: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TITLE\tTYPE\tUNITS\tID")
			fmt.Fprintln(w, "-----\t----\t-----\t--")
			for _, l := range logs {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", l.Title, l.Type, l.Units, l.ID)
			}
			w.Flush()
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var insightsGetCmd = &cobra.Command{
	Use:   "get <log-id>",
	Short: "Get insight log entries",
	Long: `Get historical data entries for an insight log.

The log-id is from 'homeyctl insights list' output.
Example: homey:device:abc123:measure_power

Resolutions:
  - last24Hours (default)
  - lastWeek
  - lastMonth
  - lastYear
  - last2Years

Examples:
  homeyctl insights get "homey:device:abc123:measure_power"
  homeyctl insights get "homey:device:abc123:measure_power" --resolution lastWeek`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logID := args[0]
		resolution, _ := cmd.Flags().GetString("resolution")

		// First, look up the log to get ownerUri and ownerId
		data, err := apiClient.GetInsights()
		if err != nil {
			return err
		}

		var logs []InsightLog
		if err := json.Unmarshal(data, &logs); err != nil {
			return fmt.Errorf("failed to parse insights: %w", err)
		}

		var ownerURI, ownerID string
		for _, log := range logs {
			if log.ID == logID {
				ownerURI = log.OwnerURI
				ownerID = log.OwnerID
				break
			}
		}

		if ownerURI == "" {
			return fmt.Errorf("log not found: %s\nUse 'homeyctl insights list' to see available logs", logID)
		}

		entries, err := apiClient.GetInsightEntries(ownerURI, ownerID, resolution)
		if err != nil {
			return err
		}

		if isTableFormat() {
			var entryList []struct {
				T time.Time   `json:"t"`
				V interface{} `json:"v"`
			}
			if err := json.Unmarshal(entries, &entryList); err != nil {
				return fmt.Errorf("failed to parse entries: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIME\tVALUE")
			fmt.Fprintln(w, "----\t-----")
			for _, e := range entryList {
				fmt.Fprintf(w, "%s\t%v\n", e.T.Local().Format("2006-01-02 15:04"), e.V)
			}
			w.Flush()
			return nil
		}

		outputJSON(entries)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(insightsCmd)
	insightsCmd.AddCommand(insightsListCmd)
	insightsCmd.AddCommand(insightsGetCmd)

	insightsGetCmd.Flags().String("resolution", "last24Hours", "Resolution: last24Hours, lastWeek, lastMonth, lastYear, last2Years")
}
