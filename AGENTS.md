# Agent Instructions

## Repository-wide Setup

Before running tests or development scripts, configure credentials and tooling with the helper
script:

1. Provide Google Cloud credentials. The Go unit tests instantiate a real `GCPService`, so
   `GOOGLE_APPLICATION_CREDENTIALS` must point at a usable service-account JSON on disk. Either
   export the JSON contents directly:

   ```bash
   export GCP_SA_KEY="$(cat /path/to/service-account.json)"
   ```

   or run the script with an existing file:
   `scripts/setup-dev-env.sh --service-account /path/to/service-account.json`. Secrets that are
   stored as base64 (common in CI environments) can also be assigned to `GCP_SA_KEY`; the setup
   script automatically decodes and validates the credentials before writing the JSON file.

2. Supply a project identifier by either exporting `GCP_PROJECT_ID`/`TEST_PROJECT_ID` or passing
   `--project <your-gcp-project-id>`.
3. Run the setup script:

   ```bash
   scripts/setup-dev-env.sh --project <your-gcp-project-id>
   ```

The script materializes the key at `.secrets/service-account.json`, keeps `.env` in sync (including
`GOOGLE_APPLICATION_CREDENTIALS`), prepares `activate-dev.sh`, and ensures Go/Python tooling is
available. Use `--force` if you need to overwrite existing credentials.

### After running the script

1. Load the generated environment and PATH adjustments:

   ```bash
   source activate-dev.sh
   ```

2. Confirm that `GOOGLE_APPLICATION_CREDENTIALS` points to the service-account file (for example
   `echo $GOOGLE_APPLICATION_CREDENTIALS`). If it is unset or the file is missing, `go test` will
   fail with errors such as `failed to create resource manager client`.
3. Run the project tests. Typical sequences are:

   ```bash
   make test          # run unit tests (requires credentials)
   make test-all      # run unit + integration suites
   ```

Integration tests default to mocked services; exporting `TEST_MODE=integration` and a real
`TEST_PROJECT_ID` switches them to real GCP APIs.
