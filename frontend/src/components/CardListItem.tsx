import { useNavigate } from "react-router-dom";
import type { Card } from "../api/types";
import { CardView } from "./CardView";

// A card row that navigates to the card's detail/edit view on click, while
// still letting the tag chips inside CardView navigate to their own route
// (they stop propagation so this handler doesn't also fire).
export function CardListItem({ card }: { card: Card }) {
  const navigate = useNavigate();
  return (
    <li
      onClick={() => navigate(`/cards/${card.id}`)}
      style={{ cursor: "pointer" }}
    >
      <CardView card={card} showBack />
    </li>
  );
}
