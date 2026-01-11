package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/langtind/homey-cli/internal/client"
	"github.com/langtind/homey-cli/internal/config"
)

var (
	cfg       *config.Config
	apiClient *client.Client

	formatFlag string

	versionInfo struct {
		Version string
		Commit  string
		Date    string
	}
)

func SetVersionInfo(version, commit, date string) {
	versionInfo.Version = version
	versionInfo.Commit = commit
	versionInfo.Date = date
}

const setupInstructions = `
Welcome to homey-cli! To get started, you need to configure your Homey connection.

1. SET YOUR HOMEY'S IP ADDRESS
   Find it in the Homey app: Settings → General → scroll down
   Then run:
     homey config set-host <ip-address>

2. CREATE AN API KEY
   Go to https://my.homey.app/
   → Settings (gear icon) → API Keys → + New API Key
   Then run:
     homey config set-token <your-token>

3. VERIFY YOUR SETUP
     homey config show
     homey devices list

For more help: homey --help
`

var rootCmd = &cobra.Command{
	Use:   "homey",
	Short: "CLI for Homey smart home",
	Long:  `A command-line interface for controlling Homey devices, flows, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if configured, show setup instructions if not
		loadedCfg, _ := config.Load()
		if loadedCfg == nil || loadedCfg.Token == "" {
			fmt.Print(setupInstructions)
			return
		}
		// If configured, show normal help
		cmd.Help()
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config for config and version commands
		cmdPath := cmd.CommandPath()
		if cmd.Name() == "config" || cmd.Name() == "version" || cmd.Name() == "help" ||
			cmd.Name() == "set-token" || cmd.Name() == "set-host" || cmd.Name() == "show" ||
			cmd.Name() == "completion" || cmdPath == "homey" {
			return nil
		}

		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.Token == "" {
			return fmt.Errorf("no API token configured. Run: homey config set-token <token>")
		}

		if formatFlag != "" {
			cfg.Format = formatFlag
		}

		apiClient = client.New(cfg)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "", "Output format: json, table (default: json)")
}

// outputJSON pretty-prints JSON data
func outputJSON(data []byte) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err == nil {
		pretty, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(pretty))
		return
	}
	fmt.Println(string(data))
}

// isTableFormat returns true if table format is requested
func isTableFormat() bool {
	return cfg != nil && cfg.Format == "table"
}
