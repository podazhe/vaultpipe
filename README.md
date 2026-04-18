# vaultpipe

> CLI tool to stream secrets from HashiCorp Vault into environment files or process envs at runtime.

---

## Installation

```bash
go install github.com/yourorg/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourorg/vaultpipe/releases).

---

## Usage

### Export secrets to a `.env` file

```bash
vaultpipe export --path secret/myapp --out .env
```

### Inject secrets into a process at runtime

```bash
vaultpipe run --path secret/myapp -- ./myapp serve
```

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | *(required)* |
| `--out` | Output file path | stdout |
| `--addr` | Vault server address | `$VAULT_ADDR` |
| `--token` | Vault token | `$VAULT_TOKEN` |
| `--format` | Output format (`env`, `json`, `yaml`) | `env` |

### Example with multiple paths

```bash
vaultpipe run \
  --path secret/myapp/db \
  --path secret/myapp/api \
  -- ./myapp serve
```

---

## Authentication

`vaultpipe` respects standard Vault environment variables:

```bash
export VAULT_ADDR=https://vault.example.com
export VAULT_TOKEN=s.xxxxxxxx
```

Token, AppRole, and AWS IAM auth methods are supported.

---

## License

MIT © yourorg