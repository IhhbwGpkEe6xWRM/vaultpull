# vaultpull

> CLI tool to sync HashiCorp Vault secrets into local `.env` files with namespace filtering

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Set your Vault address and token, then run `vaultpull` with a namespace path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-token-here"

# Pull secrets from a namespace into a local .env file
vaultpull sync --namespace secret/myapp/production --output .env
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--namespace` | Vault secret path / namespace to pull from | *(required)* |
| `--output` | Output `.env` file path | `.env` |
| `--prefix` | Filter keys by prefix | *(none)* |
| `--dry-run` | Preview output without writing to disk | `false` |

### Example Output

```env
DATABASE_URL=postgres://user:pass@host:5432/db
API_KEY=abc123
REDIS_URL=redis://localhost:6379
```

---

## Requirements

- Go 1.21+
- HashiCorp Vault with a valid token and appropriate policies

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)