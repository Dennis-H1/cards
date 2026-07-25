import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getCard, listActivity, markActivitySeen } from "../api/client";
import type { ActivityEvent, Card } from "../api/types";

export function ActivityPage({ onSeen }: { onSeen: () => void }) {
  const [events, setEvents] = useState<ActivityEvent[] | null>(null);
  const [cardsById, setCardsById] = useState<Record<number, Card>>({});
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    listActivity()
      .then((evts) => {
        setEvents(evts);

        const uniqueCardIds = [...new Set(evts.map((e) => e.card_id))];
        Promise.all(uniqueCardIds.map((id) => getCard(id).catch(() => null))).then((cards) => {
          const byId: Record<number, Card> = {};
          for (const card of cards) {
            if (card) byId[card.id] = card;
          }
          setCardsById(byId);
        });

        const hadUnseen = evts.some((e) => e.seen_at === null);
        if (hadUnseen) {
          markActivitySeen().then(onSeen);
        }
      })
      .catch((err) => setError(err.message));
  }, [onSeen]);

  if (error) return <p className="empty-state">{error}</p>;
  if (!events) return <p className="empty-state">Loading…</p>;
  if (events.length === 0) return <p className="empty-state">No activity yet.</p>;

  return (
    <ul className="event-list">
      {events.map((event) => {
        const card = cardsById[event.card_id];
        return (
          <li key={event.id} className={`event-item ${event.seen_at ? "seen" : ""}`}>
            <Link to={`/cards/${event.card_id}`}>
              <strong>{event.type === "card_created" ? "Card created" : "Card edited"}</strong>
              <div>{card ? card.front : `card #${event.card_id}`}</div>
              <small>{new Date(event.created_at).toLocaleString()}</small>
            </Link>
          </li>
        );
      })}
    </ul>
  );
}
