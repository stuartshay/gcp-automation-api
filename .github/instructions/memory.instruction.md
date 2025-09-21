- Development environment: Linux, zsh, pre-commit hooks
- Communication style: Concise, thorough, no unnecessary repetition

## Project Context

- Current project type: GCP Automation API (Go backend)
- Tech stack: Go 1.24.7, Gin, GCP SDK, Bash scripts
- Architecture patterns: Clean architecture, separation of concerns
- Key requirements: Code quality, security, automation, CI/CD

## Coding Patterns

- Preferred patterns and practices: Structured error handling, dependency injection, pre-commit
  hooks
- Code organization preferences: Modular, clean separation of logic
- Testing approaches: Unit and integration tests, pre-commit test hooks
- Documentation style: Markdown, OpenAPI

## Context7 Research History

- Structlog-gcp reviewed for GCP logging format (Python)
- Chose Go GCP logging client for structured logging
- Logging fields: severity, message, trace, request/response metadata, error details
- Logging injected via middleware, handlers log actionable events

## Conversation History

- Implemented structured logging for Cloud Run API using GCP logging client
- Logging middleware injects logger into context
- Handlers emit logs for requests, responses, errors
- All code errors resolved, build and tests pass
- Next step: Validate logs in GCP dashboard

## Notes

- Logging is now actionable and compatible with GCP dashboards
- Memory updated after implementation

# Unit Test Setup (2025-09-21)

## Project Test Setup

- AGENTS.md reviewed for test requirements
- .env and GOOGLE_APPLICATION_CREDENTIALS verified
- Development environment activated via 'Activate Dev Environment' task
- Unit tests run successfully using VS Code task
- All test suites passed (handlers, integration, validation, etc.)
- No missing credentials or configuration issues detected

## Conversation History

- Unit test setup and execution validated for GCP Automation API
- Project is ready for further development and PR review
