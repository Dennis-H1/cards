# Karteikarten — product spec

## Summary

A personal spaced-repetition flashcard app ("Karteikarten") for learning infrastructure/technical topics. The card creation pipeline runs through conversation with Claude: the user chats or reads something worth remembering, asks Claude to make a card, reviews/edits it inline in chat, and on save it lands in the app's library. The app itself handles browsing, tagging, and spaced-repetition review sessions.

Two things distinguish this from a plain Anki clone:
1. **Tag overviews** — a Claude-synthesized living summary per tag, not just a pile of disconnected cards.
2. **Design-spec sharing via MCP resource** — Claude's in-chat card previews are pixel-consistent with the real app's UI, because both read from the same source of truth.

Single user, self-hosted, mobile-first. The app sits behind a single-account
login screen — still one user, just not openly reachable by anyone who finds
the URL.

For stack, data model, the SM-2 algorithm, and the MCP server contract, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## End-to-end flow

1. User: "make a card about X" (mid-conversation, anywhere)
2. Claude fetches `resource://karteikarten/design-spec`, drafts front/back content, renders an interactive preview widget (matching the real app's visual design) with edit / save / discard actions
3. User iterates conversationally ("change the back to mention Y") until satisfied
4. User taps Save → Claude calls `create_card` → row persisted, `Event` logged
5. Next time the user opens the app, the card appears in the library with the same visual language as the chat preview
6. A bell/activity icon in the app's top bar shows unseen `Event`s (cards created/edited since last viewed) — front text, tags, timestamp, source snippet; dismisses on view
7. Review session pulls the due-queue and runs the SM-2 grading loop

---

## App structure (screens)

- **Login** — username/password gate for the single account; no signup flow
- **Review** (home) — due-queue card, tap-to-reveal, Again/Hard/Good/Easy grading, progress indicator
- **Browse** — card list, filterable by tag, search
- **Tag overview** — synthesized summary for a tag (Claude-generated via `get_tag_overview`), plus the cards under it
- **Activity feed** — chronological list of `Event`s, unseen-count badge in top bar

Mobile: bottom nav with three items (Review, Browse, Capture/Activity). No hamburger menu.

Desktop: three-pane layout — tag tree/list (left), review or browse content (center, width-constrained, not full-bleed), tag overview panel (right).

---

## Design language

- Dark-mode-first
- Background `#12151A`, card surface `#1B1F27`, primary text `#E8E6DF`
- Accent (success/correct/graduated): sage `#7C9A82`
- Accent (fail/due/warning): muted rust `#C77B58`
- Secondary/borders: slate `#5B6472`
- UI chrome font: Inter or IBM Plex Sans
- Card content font: IBM Plex Mono or JetBrains Mono (functionally correct — much card content is code/commands/config)
- Motion: flat, fast cross-fade on reveal; no gamified animation, confetti, or streak effects — the tone is "competent tool," not "gamified app"
- Signature element: due-count shown as a terminal-prompt-style line, e.g. `14 due · 3 new`

This palette and the reference card markup should live in the codebase and be exposed via the `design-spec` MCP resource described in [ARCHITECTURE.md](ARCHITECTURE.md) — it is the actual source of truth, not this document.

---

## MVP scope

1. Card CRUD + tags (many-to-many)
2. SM-2 scheduler + review session UI with grading
3. MCP server: the 5 tools + 1 resource (see [ARCHITECTURE.md](ARCHITECTURE.md))
4. Tag overview: stored markdown field, regenerable via `get_tag_overview`
5. Activity feed (Event log + top-bar badge)
6. Mobile-first responsive UI, desktop 3-pane layout

## Explicitly out of scope for MVP (backlog)

- Cloze deletion card type
- Bulk extraction (paste long text → propose N candidate cards → approve/reject)
- Auto-regeneration of tag overviews on a threshold (vs. manual trigger)
- Weak-spot detection / auto-generated easier prerequisite cards
- Related-card linking
- Anki `.apkg` import
- Semantic search / embeddings over card content
- Image/attachment support on cards
- Push/email due-card notifications
- Duplicate/near-duplicate merge tooling (beyond `search_cards` at creation time)
