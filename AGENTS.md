# Repository Guidelines

## Project Structure
- `cmd/` - Cobra command definitions. Each file is a command group (devices, flows, zones, etc.)
- `internal/client/` - HTTP client for Homey's local REST API
- `internal/config/` - Configuration management with Viper
- `main.go` - Entry point with version info

## Build & Test Commands
- `go build -o homey .` - Build the CLI binary
- `go test ./...` - Run the full test suite
- `make test` - Same as above
- `make fmt` - Format with gofumpt
- `make lint` - Run golangci-lint

## Coding Style
- Go formatting via `gofumpt` (run `make fmt`)
- Lint via `golangci-lint` (run `make lint`)
- Follow standard Go naming: exported `CamelCase`, unexported `camelCase`

## Testing Guidelines
- Use Go's `testing` package
- Run `go test ./...` before shipping changes
- Tests live alongside code as `*_test.go`

## Commit Guidelines
- Conventional Commit style: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`
- Keep commits scoped and descriptive
- Include `Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>` when AI-assisted

## Configuration
- Config file: `~/.config/homey-cli/config.toml`
- Environment variables: `HOMEY_HOST`, `HOMEY_TOKEN`, `HOMEY_FORMAT`
- Override config path: `--config /path/to/config.toml`

## Flow Creation (AI Agents)
When creating flows via CLI:

1. **Discover IDs first:**
   ```bash
   homey devices list          # Get device IDs
   homey users list            # Get user IDs
   homey flows cards --type trigger|condition|action  # Get card IDs
   homey zones list            # Get zone IDs
   ```

2. **Droptoken format for logic conditions:**
   - Use pipe (`|`) before capability: `homey:device:<device-id>|<capability>`
   - Example: `homey:device:abc123|measure_temperature`

3. **Flow JSON structure:**
   ```json
   {
     "name": "Flow Name",
     "trigger": {"id": "homey:manager:presence:user_enter", "args": {...}},
     "conditions": [{"id": "...", "droptoken": "...", "args": {...}}],
     "actions": [{"id": "homey:device:<id>:on", "args": {}}]
   }
   ```

4. **Validation:** CLI validates JSON before sending to API

## Common Tasks
- List devices: `homey devices list`
- Control device: `homey devices set "Device Name" capability value`
- Trigger flow: `homey flows trigger "Flow Name"`
- Create flow: `homey flows create flow.json`
