# Karteikarten

[![CI](https://github.com/Dennis-H1/cards/actions/workflows/ci.yml/badge.svg)](https://github.com/Dennis-H1/cards/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Dennis-H1/cards)](go.mod)

A personal, self-hosted spaced-repetition flashcard app for learning infrastructure/technical topics. Cards are drafted through conversation with Claude (via MCP), reviewed/edited inline, and saved into a library you review with an SM-2 scheduler.

Single user, self-hosted, mobile-first. See [PRODUCT.md](PRODUCT.md) for what it does and [ARCHITECTURE.md](ARCHITECTURE.md) for how it's built.

**Status**: spec only — no code yet.

## Stack

- Go backend — single binary, REST API + MCP server over one service layer
- SQLite
- React + TypeScript frontend

## Docs

- [PRODUCT.md](PRODUCT.md) — features, screens, design language, MVP scope
- [ARCHITECTURE.md](ARCHITECTURE.md) — stack, data model, SM-2 algorithm, MCP server contract
- [CLAUDE.md](CLAUDE.md) — guidance for Claude Code working in this repo
