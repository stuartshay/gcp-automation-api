#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SECRETS_DIR="$REPO_ROOT/.secrets"
DEFAULT_SERVICE_ACCOUNT_PATH="$SECRETS_DIR/service-account.json"
ENV_FILE="$REPO_ROOT/.env"
ENV_TEMPLATE="$REPO_ROOT/.env.example"
PYTHON_VENV="$REPO_ROOT/.venv"

COLOR_BLUE="\033[0;34m"
COLOR_GREEN="\033[0;32m"
COLOR_YELLOW="\033[1;33m"
COLOR_RED="\033[0;31m"
COLOR_RESET="\033[0m"

info() {
    echo -e "${COLOR_BLUE}[INFO]${COLOR_RESET} $*"
}

success() {
    echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} $*"
}

warn() {
    echo -e "${COLOR_YELLOW}[WARN]${COLOR_RESET} $*"
}

error() {
    echo -e "${COLOR_RED}[ERROR]${COLOR_RESET} $*" >&2
    exit 1
}

usage() {
    cat <<'USAGE_BLOCK'
Usage: scripts/setup-dev-env.sh [options]

Ensures development dependencies and configuration are available locally.

Options:
  --project ID            Project ID to store in .env (defaults to $GCP_PROJECT_ID if exported)
  --service-account PATH  Destination for the service account JSON file
  --force                 Overwrite generated files if they already exist
  --skip-go-tools         Skip installing Go-based helper binaries
  --skip-python           Skip Python virtual environment setup
  -h, --help              Show this help message
USAGE_BLOCK
}

INSTALL_GO_TOOLS=1
SETUP_PYTHON=1
FORCE=0
PROJECT_OVERRIDE=""
SERVICE_ACCOUNT_DEST=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --project)
            [[ $# -lt 2 ]] && error "--project requires an argument"
            PROJECT_OVERRIDE="$2"
            shift 2
            ;;
        --service-account)
            [[ $# -lt 2 ]] && error "--service-account requires an argument"
            SERVICE_ACCOUNT_DEST="$2"
            shift 2
            ;;
        --force)
            FORCE=1
            shift
            ;;
        --skip-go-tools)
            INSTALL_GO_TOOLS=0
            shift
            ;;
        --skip-python)
            SETUP_PYTHON=0
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            usage >&2
            error "Unknown option: $1"
            ;;
    esac
done

if [[ -z "$SERVICE_ACCOUNT_DEST" ]]; then
    SERVICE_ACCOUNT_DEST="$DEFAULT_SERVICE_ACCOUNT_PATH"
fi
if [[ "$SERVICE_ACCOUNT_DEST" != /* ]]; then
    SERVICE_ACCOUNT_DEST="$REPO_ROOT/$SERVICE_ACCOUNT_DEST"
fi

mkdir -p "$SECRETS_DIR"

info "Repository root: $REPO_ROOT"

if ! command -v go >/dev/null 2>&1; then
    error "Go is not installed. Install the toolchain declared in go.mod before proceeding."
fi

GO_TOOLCHAIN_REQUIRED=$(awk '/^toolchain / {print $2}' "$REPO_ROOT/go.mod")
if [[ -z "$GO_TOOLCHAIN_REQUIRED" ]]; then
    GO_TOOLCHAIN_REQUIRED="go$(awk '/^go / {print $2}' "$REPO_ROOT/go.mod")"
fi
INSTALLED_GO=$(go env GOVERSION 2>/dev/null || go version | awk '{print $3}')
if [[ "$INSTALLED_GO" != "$GO_TOOLCHAIN_REQUIRED" ]]; then
    warn "Go $INSTALLED_GO detected but go.mod requests $GO_TOOLCHAIN_REQUIRED."
else
    success "Go toolchain $INSTALLED_GO detected."
fi

GO_BIN_DIR=$(go env GOBIN)
if [[ -z "$GO_BIN_DIR" ]]; then
    GO_BIN_DIR=$(go env GOPATH)/bin
fi
mkdir -p "$GO_BIN_DIR"
if [[ ":$PATH:" != *":$GO_BIN_DIR:"* ]]; then
    warn "Go bin directory ($GO_BIN_DIR) is not on PATH. Add it or source activate-dev.sh."
else
    info "Go bin directory resolved to $GO_BIN_DIR"
fi

if ! command -v python3 >/dev/null 2>&1; then
    error "python3 is required to provision the virtual environment."
fi

if [[ "$SETUP_PYTHON" -eq 1 ]]; then
    if [[ ! -d "$PYTHON_VENV" ]]; then
        info "Creating Python virtual environment at $PYTHON_VENV"
        python3 -m venv "$PYTHON_VENV"
        success "Created virtual environment"
    else
        info "Python virtual environment already exists at $PYTHON_VENV"
    fi

    info "Installing Python developer dependencies"
    # shellcheck disable=SC1090
    source "$PYTHON_VENV/bin/activate"
    PY_DEPS_OK=1
    if ! pip install --upgrade pip; then
        warn "Unable to upgrade pip (likely due to restricted network access)."
        PY_DEPS_OK=0
    fi
    if [[ -f "$REPO_ROOT/requirements-dev.txt" ]]; then
        if ! pip install -r "$REPO_ROOT/requirements-dev.txt"; then
            warn "Unable to install Python developer requirements (run manually once network access is available)."
            PY_DEPS_OK=0
        fi
    fi
    deactivate
    if [[ "$PY_DEPS_OK" -eq 1 ]]; then
        success "Python tooling ready"
    else
        warn "Python tooling installed with warnings; retry installs inside the virtualenv when possible."
    fi
else
    warn "Skipping Python setup as requested"
fi

install_go_tool() {
    local binary="$1"
    local module="$2"
    info "Ensuring $binary (package $module)"
    if GOBIN="$GO_BIN_DIR" go install "$module"; then
        if command -v "$binary" >/dev/null 2>&1; then
            success "$binary available at $(command -v "$binary")"
        elif [[ -x "$GO_BIN_DIR/$binary" ]]; then
            success "$binary installed to $GO_BIN_DIR"
        else
            warn "$binary was installed but is not on PATH (located in $GO_BIN_DIR)"
        fi
    else
        warn "Failed to install $binary from $module"
    fi
}

if [[ "$INSTALL_GO_TOOLS" -eq 1 ]]; then
    install_go_tool "golangci-lint" "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    install_go_tool "gosec" "github.com/securego/gosec/v2/cmd/gosec@latest"
    install_go_tool "goimports" "golang.org/x/tools/cmd/goimports@latest"
    install_go_tool "wire" "github.com/google/wire/cmd/wire@latest"
    install_go_tool "swag" "github.com/swaggo/swag/cmd/swag@latest"
else
    warn "Skipping Go tooling installation as requested"
fi

ensure_env_file() {
    if [[ -f "$ENV_FILE" ]]; then
        return
    fi

    if [[ -f "$ENV_TEMPLATE" ]]; then
        cp "$ENV_TEMPLATE" "$ENV_FILE"
        success "Created $ENV_FILE from template"
    else
        warn "No .env or template found; generating minimal file"
        cat <<'ENV_FALLBACK' > "$ENV_FILE"
# Generated by scripts/setup-dev-env.sh
ENV_FALLBACK
    fi
}

update_env_var() {
    local key="$1"
    local value="$2"
    local tmp
    tmp=$(mktemp)
    if [[ -f "$ENV_FILE" ]]; then
        awk -v key="$key" -v value="$value" '
            BEGIN { updated = 0 }
            $0 ~ "^" key "=" {
                if (!updated) {
                    print key "=" value
                    updated = 1
                }
                next
            }
            { print }
            END {
                if (!updated) {
                    print key "=" value
                }
            }
        ' "$ENV_FILE" > "$tmp"
    else
        printf '%s=%s\n' "$key" "$value" > "$tmp"
    fi
    mv "$tmp" "$ENV_FILE"
}

ensure_env_file

PROJECT_ID="$PROJECT_OVERRIDE"
if [[ -z "$PROJECT_ID" && -n "${GCP_PROJECT_ID:-}" ]]; then
    PROJECT_ID="$GCP_PROJECT_ID"
fi
if [[ -z "$PROJECT_ID" && -n "${TEST_PROJECT_ID:-}" ]]; then
    PROJECT_ID="$TEST_PROJECT_ID"
fi
if [[ -n "$PROJECT_ID" ]]; then
    update_env_var "GCP_PROJECT_ID" "$PROJECT_ID"
    update_env_var "TEST_PROJECT_ID" "$PROJECT_ID"
    success "Stored project ID $PROJECT_ID in $ENV_FILE"
else
    warn "GCP project ID not provided. Pass --project or export GCP_PROJECT_ID before running integration tests."
fi

CREDENTIALS_PATH=""
if [[ -n "${GCP_SA_KEY:-}" ]]; then
    if [[ -f "$SERVICE_ACCOUNT_DEST" && "$FORCE" -eq 0 ]]; then
        info "Service account file already exists at $SERVICE_ACCOUNT_DEST (use --force to overwrite)"
    else
        info "Writing service account credentials to $SERVICE_ACCOUNT_DEST"
        mkdir -p "$(dirname "$SERVICE_ACCOUNT_DEST")"
        old_umask=$(umask)
        umask 077
        printf '%s' "$GCP_SA_KEY" > "$SERVICE_ACCOUNT_DEST"
        umask "$old_umask"
        chmod 600 "$SERVICE_ACCOUNT_DEST"
        success "Service account key materialized"
    fi
    CREDENTIALS_PATH="$SERVICE_ACCOUNT_DEST"
elif [[ -n "${GOOGLE_APPLICATION_CREDENTIALS:-}" ]]; then
    CREDENTIALS_PATH="$GOOGLE_APPLICATION_CREDENTIALS"
    info "Using existing GOOGLE_APPLICATION_CREDENTIALS=$CREDENTIALS_PATH"
else
    warn "GCP_SA_KEY not exported; skipping credential file creation"
fi

if [[ -n "$CREDENTIALS_PATH" ]]; then
    update_env_var "GOOGLE_APPLICATION_CREDENTIALS" "$CREDENTIALS_PATH"
    success "Updated GOOGLE_APPLICATION_CREDENTIALS entry in $ENV_FILE"
else
    warn "GOOGLE_APPLICATION_CREDENTIALS not configured in $ENV_FILE"
fi

cat <<SETUP_SUMMARY

${COLOR_GREEN}Setup complete!${COLOR_RESET}

Next steps:
  1. Source activate-dev.sh to load the virtualenv and environment variables:
       source "$REPO_ROOT/activate-dev.sh"
  2. Run make targets such as 'make test' or 'make dev'.

Generated files:
  - Python virtualenv: $PYTHON_VENV
  - Service account (if created): $SERVICE_ACCOUNT_DEST
  - Environment file: $ENV_FILE

SETUP_SUMMARY
