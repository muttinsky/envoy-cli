# envoy-cli

> A CLI tool for managing and diffing .env files across multiple environments with secret masking support.

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envoy-cli/releases).

## Usage

```bash
# Diff two .env files
envoy diff .env.staging .env.production

# Diff with secrets masked
envoy diff .env.staging .env.production --mask-secrets

# List all variables in an env file
envoy list .env.staging

# Validate an env file against a template
envoy validate .env --template .env.example

# Sync missing keys from one env to another
envoy sync .env.example .env.local
```

### Example Output

```
~ DB_HOST        staging.db.local  →  prod.db.internal
+ NEW_FEATURE_FLAG               →  true
- DEPRECATED_KEY  old_value       →  (missing)
```

Keys marked with `~` are changed, `+` are added, and `-` are removed.

## Configuration

`envoy-cli` looks for a `.envoy.yaml` config file in the project root to define secret patterns and environment groups.

```yaml
secrets:
  - "*_KEY"
  - "*_SECRET"
  - "*_TOKEN"
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)