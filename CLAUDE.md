# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Project

Karteikarten: a personal, self-hosted spaced-repetition flashcard app where cards are drafted and created through conversation with Claude (via MCP), then reviewed in a dedicated web app. Go backend (REST API + MCP server, one binary), SQLite, React/TypeScript frontend.

This repo currently holds only the spec — no code has been written yet.

## Docs

- **[PRODUCT.md](PRODUCT.md)** — what the app does, end-to-end flow, screens, design language, MVP scope. Read before making product or UX decisions.
- **[ARCHITECTURE.md](ARCHITECTURE.md)** — stack, data model, SM-2 scheduling algorithm, MCP server contract (tools/resource/auth). Read before writing backend code, touching the data model, or implementing the MCP server.
