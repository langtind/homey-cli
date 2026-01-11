# homeyctl

![homeyctl banner](banner.png)

A command-line interface for controlling [Homey](https://homey.app) smart home devices via the local API.

> **Note:** The binary is named `homeyctl` to avoid conflicts with Athom's official `homey` CLI tool used for app development.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap langtind/tap
brew install homey-cli
```

### Download Binary

Download from [Releases](https://github.com/langtind/homey-cli/releases) and add to your PATH.

### Build from Source

```bash
go install github.com/langtind/homey-cli@latest
```

## Configuration

### 1. Find your Homey's IP address

Open the Homey app → Settings → General → scroll down to find the local IP address (e.g., `192.168.1.100`).

### 2. Get your API token

1. Go to [my.homey.app](https://my.homey.app/)
2. Log in and select your Homey
3. Click **Settings** (gear icon, bottom left)
4. Click **API Keys**
5. Click **+ New API Key** to create a new token
6. Copy the generated token

### 3. Configure homeyctl

```bash
# Set your Homey's IP address
homeyctl config set-host 192.168.1.100

# Set your API token (paste the token you copied)
homeyctl config set-token <your-token>

# Verify configuration
homeyctl config show
```

Configuration is stored in `~/.config/homeyctl/config.toml`.

## Usage

### Devices

```bash
# List all devices
homeyctl devices list

# Get device details
homeyctl devices get "Living Room Light"

# Control devices
homeyctl devices set "Living Room Light" onoff true
homeyctl devices set "Living Room Light" dim 0.5
homeyctl devices set "Thermostat" target_temperature 22

# Delete a device
homeyctl devices delete "Old Device"
```

### Flows

```bash
# List all flows
homeyctl flows list

# Get flow details
homeyctl flows get "My Flow"

# Trigger a flow
homeyctl flows trigger "Good Morning"

# Update a flow (partial/merge update)
homeyctl flows update "My Flow" updated-flow.json

# Delete a flow
homeyctl flows delete "Old Flow"

# List available flow cards (for creating flows)
homeyctl flows cards --type trigger
homeyctl flows cards --type condition
homeyctl flows cards --type action
```

#### Creating Flows

Create flows from JSON files. See `homeyctl flows create --help` for full documentation.

```bash
# Create a simple flow
homeyctl flows create my-flow.json

# Create an advanced flow
homeyctl flows create --advanced my-advanced-flow.json
```

Example flow JSON (turn on heater when user arrives and temperature is low):

```json
{
  "name": "Heat office on arrival",
  "trigger": {
    "id": "homey:manager:presence:user_enter",
    "args": {
      "user": {"id": "<user-id>", "name": "User Name"}
    }
  },
  "conditions": [
    {
      "id": "homey:manager:logic:lt",
      "droptoken": "homey:device:<device-id>|measure_temperature",
      "args": {"comparator": 20}
    }
  ],
  "actions": [
    {"id": "homey:device:<device-id>:on", "args": {}},
    {"id": "homey:device:<device-id>:target_temperature_set", "args": {"target_temperature": 23}}
  ]
}
```

Note: For droptokens, use pipe (`|`) before capability name.

### Users

```bash
# List users (useful for getting user IDs for flows)
homeyctl users list
```

### Zones

```bash
# List zones
homeyctl zones list

# Delete a zone
homeyctl zones delete "Unused Room"
```

### Apps

```bash
# List installed apps
homeyctl apps list

# Restart an app
homeyctl apps restart com.some.app
```

### Energy

```bash
# Show live power usage
homeyctl energy live

# Show today's energy report
homeyctl energy report day

# Show report for specific date
homeyctl energy report day --date 2025-01-10

# Weekly and monthly reports
homeyctl energy report week
homeyctl energy report month --date 2025-12

# Show electricity prices (dynamic)
homeyctl energy price

# Set fixed electricity price (e.g., Norgespris 0.50 kr/kWh)
homeyctl energy price set 0.50

# Get/set price type (fixed, dynamic, disabled)
homeyctl energy price type
homeyctl energy price type fixed
```

### Insights

```bash
# List all insight logs
homeyctl insights list

# Get historical data for a log
homeyctl insights get "homey:device:abc123:measure_power"

# With different resolutions
homeyctl insights get "homey:device:abc123:measure_power" --resolution lastWeek
```

### Variables

```bash
# List logic variables
homeyctl variables list

# Get/set variable
homeyctl variables get "my_variable"
homeyctl variables set "my_variable" 42

# Create/delete variable
homeyctl variables create "new_var" number 0
homeyctl variables delete "new_var"
```

### Notifications

```bash
# Send notification to Homey timeline
homeyctl notifications send "Hello from CLI"

# List notifications
homeyctl notifications list
```

### System

```bash
# System info
homeyctl system info

# Reboot Homey (use with caution)
homeyctl system reboot
```

## Output Formats

```bash
# JSON output (default)
homeyctl devices list

# Table output
homeyctl devices list --format table

# Set default format
homeyctl config set-format table
```

## AI Assistant Support

Get context for AI assistants (Claude, ChatGPT, etc.) to help them use homeyctl:

```bash
homeyctl ai
```

This outputs documentation, examples, and flow JSON format - perfect for pasting into an AI chat or including in your project's context.

## Environment Variables

All config options can be set via environment variables:

```bash
export HOMEY_HOST=192.168.1.100
export HOMEY_TOKEN=your-token
export HOMEY_FORMAT=table
```