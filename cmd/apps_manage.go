package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var appsInstallCmd = &cobra.Command{
	Use:   "install <app-id>",
	Short: "Install an app from the Homey App Store",
	Long: `Install an app from the Homey App Store.

The app-id is the unique identifier shown in the App Store URL.
For example: com.athom.homeyscript

Examples:
  homeyctl apps install com.athom.homeyscript
  homeyctl apps install com.fibaro --channel test`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID := args[0]
		channel, _ := cmd.Flags().GetString("channel")

		result, err := apiClient.InstallApp(appID, channel)
		if err != nil {
			return err
		}

		if isTableFormat() {
			var app struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Version string `json:"version"`
			}
			json.Unmarshal(result, &app)
			fmt.Printf("Installed app: %s v%s\n", app.Name, app.Version)
			return nil
		}

		outputJSON(result)
		return nil
	},
}

var appsUninstallCmd = &cobra.Command{
	Use:   "uninstall <app>",
	Short: "Uninstall an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := findApp(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.UninstallApp(app.ID); err != nil {
			return err
		}

		fmt.Printf("Uninstalled app: %s\n", app.Name)
		return nil
	},
}

var appsEnableCmd = &cobra.Command{
	Use:   "enable <app>",
	Short: "Enable an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := findApp(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.EnableApp(app.ID); err != nil {
			return err
		}

		fmt.Printf("Enabled app: %s\n", app.Name)
		return nil
	},
}

var appsDisableCmd = &cobra.Command{
	Use:   "disable <app>",
	Short: "Disable an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := findApp(args[0])
		if err != nil {
			return err
		}

		if err := apiClient.DisableApp(app.ID); err != nil {
			return err
		}

		fmt.Printf("Disabled app: %s\n", app.Name)
		return nil
	},
}

var appsUpdateCmd = &cobra.Command{
	Use:   "update <app>",
	Short: "Update app settings",
	Long: `Update app settings like autoupdate.

Examples:
  homeyctl apps update com.fibaro --autoupdate=true
  homeyctl apps update com.fibaro --autoupdate=false`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := findApp(args[0])
		if err != nil {
			return err
		}

		autoupdate, _ := cmd.Flags().GetBool("autoupdate")

		updates := map[string]interface{}{
			"autoupdate": autoupdate,
		}

		if err := apiClient.UpdateApp(app.ID, updates); err != nil {
			return err
		}

		fmt.Printf("Updated app: %s (autoupdate: %v)\n", app.Name, autoupdate)
		return nil
	},
}

var appsSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Manage app settings",
}

var appsSettingsListCmd = &cobra.Command{
	Use:   "list <app>",
	Short: "List app settings",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := findApp(args[0])
		if err != nil {
			return err
		}

		data, err := apiClient.GetAppSettings(app.ID)
		if err != nil {
			return err
		}

		outputJSON(data)
		return nil
	},
}

var appsSettingsSetCmd = &cobra.Command{
	Use:   "set <app> <setting> <value>",
	Short: "Set an app setting",
	Long: `Set an app setting value.

Examples:
  homeyctl apps settings set com.fibaro pollInterval 60
  homeyctl apps settings set myapp.id debugMode true`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := findApp(args[0])
		if err != nil {
			return err
		}

		settingName := args[1]
		value := parseValue(args[2])

		if err := apiClient.SetAppSetting(app.ID, settingName, value); err != nil {
			return err
		}

		fmt.Printf("Set %s.%s = %v\n", app.Name, settingName, value)
		return nil
	},
}

var appsUsageCmd = &cobra.Command{
	Use:   "usage <app>",
	Short: "Show app resource usage",
	Long: `Show app resource usage (CPU, memory).

Examples:
  homeyctl apps usage com.fibaro`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := findApp(args[0])
		if err != nil {
			return err
		}

		data, err := apiClient.GetAppUsage(app.ID)
		if err != nil {
			return err
		}

		if isTableFormat() {
			var usage struct {
				CPU    float64 `json:"cpu"`
				Memory int64   `json:"memory"`
			}
			if err := json.Unmarshal(data, &usage); err != nil {
				return err
			}

			fmt.Printf("App: %s\n", app.Name)
			fmt.Printf("CPU:    %.2f%%\n", usage.CPU*100)
			fmt.Printf("Memory: %.2f MB\n", float64(usage.Memory)/(1024*1024))
			return nil
		}

		outputJSON(data)
		return nil
	},
}

func init() {
	appsCmd.AddCommand(appsInstallCmd)
	appsInstallCmd.Flags().String("channel", "", "App channel (live, test)")

	appsCmd.AddCommand(appsUninstallCmd)
	appsCmd.AddCommand(appsEnableCmd)
	appsCmd.AddCommand(appsDisableCmd)

	appsCmd.AddCommand(appsUpdateCmd)
	appsUpdateCmd.Flags().Bool("autoupdate", true, "Enable/disable autoupdate")

	appsCmd.AddCommand(appsSettingsCmd)
	appsSettingsCmd.AddCommand(appsSettingsListCmd)
	appsSettingsCmd.AddCommand(appsSettingsSetCmd)

	appsCmd.AddCommand(appsUsageCmd)
}
