import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getTagOverview } from "../api/client";
import type { TagOverview } from "../api/types";
import { CardListItem } from "../components/CardListItem";

export function TagOverviewPage() {
  const { name } = useParams<{ name: string }>();
  const [overview, setOverview] = useState<TagOverview | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!name) return;
    setOverview(null);
    getTagOverview(name)
      .then(setOverview)
      .catch((err) => setError(err.message));
  }, [name]);

  return (
    <>
      <Link to="/browse">&larr; Back to Browse</Link>

      {error ? <p className="empty-state">{error}</p> : null}
      {!overview && !error ? <p className="empty-state">Loading…</p> : null}

      {overview ? (
        <>
          <h2>{overview.tag.name}</h2>
          {overview.tag.overview ? (
            <p>{overview.tag.overview}</p>
          ) : (
            <p className="empty-state">No synthesized overview yet for this tag.</p>
          )}
          <ul className="card-list">
            {overview.cards.map((card) => (
              <CardListItem key={card.id} card={card} />
            ))}
          </ul>
        </>
      ) : null}
    </>
  );
}
