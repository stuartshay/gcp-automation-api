---
applyTo: "**"
---

# User Memory

## User Preferences

- Programming languages: Go, Bash
- Code style preferences: Follow idiomatic Go and shell best practices
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

- No Context7 research performed for shellcheck yet

## Conversation History

- Ran pre-commit shellcheck --all-files
- Found SC2046 in activate-dev.sh (unquoted export)
- Found SC2155 in install.sh (local assignment masking return value)
- Fixed both warnings
- Re-ran shellcheck, all warnings resolved

## Notes

- All shellcheck warnings in scripts are now fixed and pre-commit passes cleanly
