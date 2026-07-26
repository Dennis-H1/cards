import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listTags, searchCards } from "../api/client";
import type { Card, Tag } from "../api/types";
import { CardListItem } from "../components/CardListItem";

export function BrowsePage() {
  const [query, setQuery] = useState("");
  const [cards, setCards] = useState<Card[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    listTags().then(setTags).catch((err) => setError(err.message));
  }, []);

  useEffect(() => {
    const timeout = setTimeout(() => {
      searchCards(query).then(setCards).catch((err) => setError(err.message));
    }, 200);
    return () => clearTimeout(timeout);
  }, [query]);

  return (
    <>
      <div style={{ display: "flex", gap: "0.5rem", marginBottom: "1rem" }}>
        <input
          className="search-input"
          style={{ marginBottom: 0 }}
          placeholder="Search cards…"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <Link to="/cards/new" className="reveal-button" style={{ width: "auto", marginTop: 0, whiteSpace: "nowrap" }}>
          + New
        </Link>
      </div>

      {tags.length > 0 ? (
        <div style={{ marginBottom: "1rem" }}>
          {tags.map((tag) => (
            <Link key={tag.id} to={`/tags/${encodeURIComponent(tag.name)}`} className="kk-tag">
              {tag.name} ({tag.card_count})
            </Link>
          ))}
        </div>
      ) : null}

      {error ? <p className="empty-state">{error}</p> : null}

      {cards.length === 0 && !error ? (
        <p className="empty-state">No cards found.</p>
      ) : (
        <ul className="card-list">
          {cards.map((card) => (
            <CardListItem key={card.id} card={card} />
          ))}
        </ul>
      )}
    </>
  );
}
