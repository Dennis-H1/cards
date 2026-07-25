# Karteikarten — architecture

For what the app does, screens, and design language, see [PRODUCT.md](PRODUCT.md).

## Stack

- **Backend**: Go — single binary serving both a REST API and an MCP server over the same service layer
- **Database**: SQLite (single-user; simple backup story; can migrate to Postgres later if ever needed)
- **Frontend**: React + TypeScript, mobile-first responsive, richer 3-pane layout at desktop widths
- **Deployment**: self-hosted on the user's Hetzner VPS, behind existing nginx-proxy and Cloudflare Zero Trust setup

---

## Data model

```
Card
  id            uuid/int, pk
  front         text (markdown)
  back          text (markdown)
  source        text, nullable   -- e.g. URL or short free-text context ("from chat about X")
  created_at    timestamp
  updated_at    timestamp

Tag
  id            uuid/int, pk
  name          text, unique

CardTag
  card_id       fk -> Card
  tag_id        fk -> Tag

Review
  card_id       fk -> Card, pk
  ease_factor   float, default 2.5
  interval_days int, default 0
  repetitions   int, default 0
  due_at        timestamp
  review_count  int, default 0
  last_reviewed_at timestamp, nullable

Event
  id            uuid/int, pk
  type          enum(card_created, card_edited)
  card_id       fk -> Card
  created_at    timestamp
  seen_at       timestamp, nullable   -- null = unseen, drives the activity feed badge count
```

Card content (`front`/`back`) is stored as **Markdown**, not plain text or raw HTML — this gives real tables, fenced code blocks, and inline code in cards (important for infra content) while avoiding the XSS risk of storing raw HTML long-term. Render with a standard Markdown component (e.g. `react-markdown` + a code-highlighting plugin) in the frontend.

---

## Spaced repetition (SM-2)

Standard SM-2 algorithm. Card states: **New → Learning → Review (graduated) → Lapsed (on fail, loops back into relearning)**.

Grading buttons: **Again / Hard / Good / Easy**.

Per-review update logic:

```
if grade == Again:
    repetitions = 0
    interval_days = 1
else:
    if repetitions == 0:
        interval_days = 1
    elif repetitions == 1:
        interval_days = 6
    else:
        interval_days = round(interval_days * ease_factor)
    repetitions += 1

ease_factor = max(1.3, ease_factor + (0.1 - (5 - grade_score) * (0.08 + (5 - grade_score) * 0.02)))
due_at = now + interval_days
```

Map Again/Hard/Good/Easy to grade_score values consistent with classic SM-2 (e.g. 0/3/4/5) — exact mapping is an implementation detail, but Again must always reset repetitions and force a short interval.

Due-queue query for a review session: `SELECT * FROM cards JOIN reviews WHERE due_at <= now() ORDER BY due_at ASC`.

---

## MCP server

Exposes both tools (actions) and one resource (read-only reference data). No tool call is used purely for "drafting" — draft/iterate happens in plain conversation with Claude rendering a preview widget; only actions that actually persist data are tools.

### Tools

| Tool | Purpose |
|---|---|
| `create_card(front, back, tags[], source?)` | Persists a new card, creates default `Review` row, writes a `card_created` `Event` |
| `update_card(card_id, front?, back?, tags?)` | Updates an existing card, writes a `card_edited` `Event` |
| `list_tags()` | Returns all tags with card counts |
| `get_tag_overview(tag)` | Returns (or triggers synthesis of) a summary of all cards under a tag — used for the "living overview" feature |
| `search_cards(query)` | Text search over card front/back — used by Claude to check for near-duplicates before creating a new card |

### Resource

| Resource URI | Purpose |
|---|---|
| `resource://karteikarten/design-spec` | Returns the canonical design tokens (colors, fonts, spacing) and a reference HTML/CSS snippet for a card component. Claude fetches this before rendering any in-chat card preview widget, so previews always match the real app pixel-for-pixel. This is the single source of truth for card visual design — update it here when the app's UI changes, nothing else needs to be kept in sync manually. |

### Auth

The MCP server is reachable from the public internet (behind Cloudflare Zero Trust). It must require a credential (API key or OAuth) even though it's single-user — do not expose an open write endpoint.

---

## CI/CD

GitHub Actions (`.github/workflows/ci.yml`) runs `gofmt`/`go vet`/`go test -race` and a Docker build on every push and PR. On push to `main` (or a manual `workflow_dispatch` run), it additionally builds and pushes the image to `ghcr.io/dennis-h1/cards` (tags: `latest`, `sha-<commit>`), then a `deploy` job runs the Ansible playbook (`ansible/deploy.yml`) against the VPS directly from Actions, using repo secrets (`SSH_PRIVATE_KEY`, `VPS_HOST`, `VPS_USER`, `MCP_API_KEY`) instead of the local Ansible Vault. `docker-compose.yml` pulls the GHCR image and joins the existing `nginx-proxy` network; the playbook copies the compose file, templates `.env`, and runs `docker compose pull && up -d`. No build step runs on the VPS.

The same playbook can also be run manually from a dev machine (`ansible/README.md`), using `inventory.ini` + a local Ansible Vault instead of GitHub secrets -- useful if Actions is down or for a one-off deploy without pushing to `main`.

---

## Open items to resolve before or during build

- **Auth scheme** for the MCP server (API key vs OAuth) and how it's provisioned
- **Conflict handling**: last-write-wins is acceptable for single-user; confirm no additional locking is needed
- **Backup strategy** for the SQLite file (e.g. nightly cron copy to durable storage)
- **FSRS vs SM-2**: SM-2 is the MVP choice for simplicity/transparency; consider FSRS later if card volume grows large and scheduling accuracy matters more
