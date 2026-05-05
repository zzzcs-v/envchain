# envchain

> CLI tool to manage and chain environment configs across dev/staging/prod contexts

---

## Installation

```bash
go install github.com/yourusername/envchain@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourusername/envchain/releases).

---

## Usage

Define your environment configs in a `.envchain.yaml` file:

```yaml
base:
  APP_NAME: myapp
  LOG_LEVEL: info

staging:
  inherits: base
  API_URL: https://staging.api.example.com
  LOG_LEVEL: debug

prod:
  inherits: base
  API_URL: https://api.example.com
```

Then run your command with a resolved environment:

```bash
# Resolve and export staging environment
envchain run --env staging -- ./myapp

# Print resolved variables for an environment
envchain show --env prod

# Chain multiple contexts together
envchain run --env base,staging -- go run main.go
```

Resolved variables are merged in order, with later contexts overriding earlier ones.

---

## Commands

| Command | Description |
|---------|-------------|
| `run` | Run a command with the resolved environment |
| `show` | Print resolved environment variables |
| `validate` | Validate your `.envchain.yaml` config |

---

## License

[MIT](LICENSE)