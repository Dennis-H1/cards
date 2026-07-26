import { Link } from "react-router-dom";
import type { Card } from "../api/types";

export function CardView({ card, showBack }: { card: Card; showBack: boolean }) {
  return (
    <div className="kk-card">
      <div className="kk-card__front">{card.front}</div>
      {showBack ? <div className="kk-card__back">{card.back}</div> : null}
      {card.tags && card.tags.length > 0 ? (
        <div style={{ marginTop: "1rem" }}>
          {card.tags.map((tag) => (
            <Link
              className="kk-tag"
              key={tag}
              to={`/tags/${encodeURIComponent(tag)}`}
              onClick={(e) => e.stopPropagation()}
            >
              {tag}
            </Link>
          ))}
        </div>
      ) : null}
    </div>
  );
}
