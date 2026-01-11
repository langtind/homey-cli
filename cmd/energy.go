package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var energyCmd = &cobra.Command{
	Use:   "energy",
	Short: "Manage energy",
	Long:  `View Homey energy usage, reports, and electricity prices.`,
}

var energyLiveCmd = &cobra.Command{
	Use:   "live",
	Short: "Show live energy usage",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetEnergyLive()
		if err != nil {
			return err
		}

		if isTableFormat() {
			var report struct {
				ZoneName       string `json:"zoneName"`
				TotalConsumed  struct{ W *float64 } `json:"totalConsumed"`
				TotalGenerated struct{ W *float64 } `json:"totalGenerated"`
				Items          []struct {
					Type   string  `json:"type"`
					ID     string  `json:"id"`
					Name   *string `json:"name"`
					Values struct {
						W *float64 `json:"W"`
					} `json:"values"`
				} `json:"items"`
			}
			if err := json.Unmarshal(data, &report); err != nil {
				return fmt.Errorf("failed to parse energy data: %w", err)
			}

			fmt.Printf("Zone: %s\n", report.ZoneName)
			if report.TotalConsumed.W != nil {
				fmt.Printf("Total consumed: %.1f W\n", *report.TotalConsumed.W)
			}
			if report.TotalGenerated.W != nil && *report.TotalGenerated.W > 0 {
				fmt.Printf("Total generated: %.1f W\n", *report.TotalGenerated.W)
			}
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "DEVICE\tPOWER (W)")
			fmt.Fprintln(w, "------\t---------")
			for _, item := range report.Items {
				if item.Type == "device" && item.Name != nil && item.Values.W != nil {
					fmt.Fprintf(w, "%s\t%.1f\n", *item.Name, *item.Values.W)
				}
			}
			w.Flush()
			return nil
		}

		outputJSON(data)
		return nil
	},
}

var energyReportCmd = &cobra.Command{
	Use:   "report [day|week|month]",
	Short: "Show energy report",
	Long: `Show energy consumption report for a period.

Periods:
  day   - Daily report (default: today)
  week  - Weekly report
  month - Monthly report

Examples:
  homey energy report day
  homey energy report day --date 2025-01-10
  homey energy report week
  homey energy report month --date 2025-01`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		period := "day"
		if len(args) > 0 {
			period = args[0]
		}
		date, _ := cmd.Flags().GetString("date")

		// Default to current period
		if date == "" {
			now := time.Now()
			switch period {
			case "day":
				date = now.Format("2006-01-02")
			case "week":
				year, week := now.ISOWeek()
				date = fmt.Sprintf("%d-W%02d", year, week)
			case "month":
				date = now.Format("2006-01")
			}
		}

		var data json.RawMessage
		var err error

		switch period {
		case "day":
			data, err = apiClient.GetEnergyReportDay(date)
		case "week":
			data, err = apiClient.GetEnergyReportWeek(date)
		case "month":
			data, err = apiClient.GetEnergyReportMonth(date)
		default:
			return fmt.Errorf("invalid period: %s (use: day, week, month)", period)
		}

		if err != nil {
			return err
		}

		if isTableFormat() {
			return printEnergyReportTable(data, period, date)
		}

		outputJSON(data)
		return nil
	},
}

type deviceEnergy struct {
	Name   string   `json:"name"`
	Period *float64 `json:"period"`
	Total  *float64 `json:"total"`
}

func printEnergyReportTable(data json.RawMessage, period, date string) error {
	var report struct {
		Date        string `json:"date"`
		Electricity struct {
			ConsumedPeriod  *float64 `json:"consumedPeriod"`
			GeneratedPeriod *float64 `json:"generatedPeriod"`
			ImportedPeriod  *float64 `json:"importedPeriod"`
			Devices         struct {
				Consumed         map[string]deviceEnergy `json:"consumed"`
				EVChargerCharged map[string]deviceEnergy `json:"evChargerCharged"`
				Imported         map[string]deviceEnergy `json:"imported"`
			} `json:"devices"`
		} `json:"electricity"`
	}
	if err := json.Unmarshal(data, &report); err != nil {
		outputJSON(data)
		return nil
	}

	fmt.Printf("Energy Report: %s (%s)\n", report.Date, period)
	fmt.Println()

	if report.Electricity.ConsumedPeriod != nil {
		fmt.Printf("Total consumed: %.2f kWh\n", *report.Electricity.ConsumedPeriod)
	}
	if report.Electricity.ImportedPeriod != nil {
		fmt.Printf("Total imported: %.2f kWh\n", *report.Electricity.ImportedPeriod)
	}
	if report.Electricity.GeneratedPeriod != nil && *report.Electricity.GeneratedPeriod > 0 {
		fmt.Printf("Total generated: %.2f kWh\n", *report.Electricity.GeneratedPeriod)
	}

	// Show consumed devices
	if len(report.Electricity.Devices.Consumed) > 0 {
		fmt.Println("\nDevices:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  DEVICE\tPERIOD\tTOTAL")
		fmt.Fprintln(w, "  ------\t------\t-----")
		for _, d := range report.Electricity.Devices.Consumed {
			periodStr := "-"
			totalStr := "-"
			if d.Period != nil {
				periodStr = fmt.Sprintf("%.2f kWh", *d.Period)
			}
			if d.Total != nil {
				totalStr = fmt.Sprintf("%.2f kWh", *d.Total)
			}
			fmt.Fprintf(w, "  %s\t%s\t%s\n", d.Name, periodStr, totalStr)
		}
		w.Flush()
	}

	// Show EV chargers separately if present
	if len(report.Electricity.Devices.EVChargerCharged) > 0 {
		fmt.Println("\nEV Chargers:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  CHARGER\tPERIOD\tTOTAL")
		fmt.Fprintln(w, "  -------\t------\t-----")
		for _, d := range report.Electricity.Devices.EVChargerCharged {
			periodStr := "-"
			totalStr := "-"
			if d.Period != nil {
				periodStr = fmt.Sprintf("%.2f kWh", *d.Period)
			}
			if d.Total != nil {
				totalStr = fmt.Sprintf("%.2f kWh", *d.Total)
			}
			fmt.Fprintf(w, "  %s\t%s\t%s\n", d.Name, periodStr, totalStr)
		}
		w.Flush()
	}

	return nil
}

var energyPriceCmd = &cobra.Command{
	Use:   "price",
	Short: "Show current electricity price",
	RunE: func(cmd *cobra.Command, args []string) error {
		date := time.Now().Format("2006-01-02")
		data, err := apiClient.GetElectricityPrice(date)
		if err != nil {
			return err
		}

		if isTableFormat() {
			var prices struct {
				PriceUnit         string `json:"priceUnit"`
				PricesPerInterval []struct {
					PeriodStart string  `json:"periodStart"`
					PeriodEnd   string  `json:"periodEnd"`
					Value       float64 `json:"value"`
				} `json:"pricesPerInterval"`
			}
			if err := json.Unmarshal(data, &prices); err != nil {
				outputJSON(data)
				return nil
			}

			now := time.Now()
			fmt.Printf("Electricity prices for %s (%s)\n\n", date, prices.PriceUnit)

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIME\tPRICE")
			fmt.Fprintln(w, "----\t-----")
			for _, p := range prices.PricesPerInterval {
				start, _ := time.Parse(time.RFC3339, p.PeriodStart)
				end, _ := time.Parse(time.RFC3339, p.PeriodEnd)
				marker := ""
				if now.After(start) && now.Before(end) {
					marker = " <-- now"
				}
				fmt.Fprintf(w, "%s-%s\t%.2f%s\n",
					start.Local().Format("15:04"),
					end.Local().Format("15:04"),
					p.Value,
					marker)
			}
			w.Flush()
			return nil
		}

		outputJSON(data)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(energyCmd)
	energyCmd.AddCommand(energyLiveCmd)
	energyCmd.AddCommand(energyReportCmd)
	energyCmd.AddCommand(energyPriceCmd)

	energyReportCmd.Flags().String("date", "", "Date for report (format: YYYY-MM-DD for day, YYYY-MM for month)")
}
