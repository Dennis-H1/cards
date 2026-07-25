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
      <input
        className="search-input"
        placeholder="Search cards…"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
      />

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
