package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var devicesGetSettingsCmd = &cobra.Command{
	Use:   "get-settings <name-or-id>",
	Short: "Get device settings",
	Long: `Get all settings for a device.

This shows configurable settings like zone_activity_disabled, climate_exclude,
and driver-specific settings.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := findDevice(args[0])
		if err != nil {
			return err
		}

		settings, err := apiClient.GetDeviceSettings(device.ID)
		if err != nil {
			return err
		}

		if isTableFormat() {
			var settingsMap map[string]interface{}
			if err := json.Unmarshal(settings, &settingsMap); err != nil {
				return fmt.Errorf("failed to parse settings: %w", err)
			}

			fmt.Printf("Settings for %s:\n\n", device.Name)
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SETTING\tVALUE")
			fmt.Fprintln(w, "-------\t-----")
			for key, val := range settingsMap {
				fmt.Fprintf(w, "%s\t%v\n", key, val)
			}
			w.Flush()
			return nil
		}

		outputJSON(settings)
		return nil
	},
}

var devicesSetSettingCmd = &cobra.Command{
	Use:   "set-setting <name-or-id> <setting-key> <value>",
	Short: "Set device setting",
	Long: `Set a device setting value.

Device settings are different from capabilities - they configure device behavior
rather than control it. Common settings include:
  - zone_activity_disabled: Exclude sensor from zone activity detection
  - climate_exclude: Exclude device from climate control

Examples:
  homeyctl devices set-setting "Motion Sensor" zone_activity_disabled true
  homeyctl devices set-setting "Thermostat" climate_exclude false`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]
		settingKey := args[1]
		valueStr := args[2]

		device, err := findDevice(nameOrID)
		if err != nil {
			return err
		}

		value := parseValue(valueStr)

		settings := map[string]interface{}{
			settingKey: value,
		}

		if err := apiClient.SetDeviceSetting(device.ID, settings); err != nil {
			if strings.Contains(err.Error(), "Missing Scopes") {
				return fmt.Errorf(`permission denied: changing device settings requires 'homey.device' scope

OAuth tokens only support 'homey.device.control' (for on/off, dim, etc.),
not full device access needed for settings.

To change device settings, create an API key at my.homey.app:
  1. Go to https://my.homey.app
  2. Select your Homey → Settings → API Keys
  3. Create a new API key (it will have full access)
  4. Run: homeyctl config set-token <your-api-key>`)
			}
			return err
		}

		fmt.Printf("Set %s setting %s = %v\n", device.Name, settingKey, value)
		return nil
	},
}

func init() {
	devicesCmd.AddCommand(devicesGetSettingsCmd)
	devicesCmd.AddCommand(devicesSetSettingCmd)
}
