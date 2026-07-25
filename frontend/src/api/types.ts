export interface Card {
  id: number;
  front: string;
  back: string;
  source: string | null;
  created_at: string;
  updated_at: string;
  tags: string[];
}

export interface Tag {
  id: number;
  name: string;
  overview: string | null;
  overview_updated_at: string | null;
  card_count: number;
}

export interface Review {
  card_id: number;
  ease_factor: number;
  interval_days: number;
  repetitions: number;
  due_at: string;
  review_count: number;
  last_reviewed_at: string | null;
}

export interface DueCard {
  card: Card;
  review: Review;
}

export interface TagOverview {
  tag: Tag;
  cards: Card[];
}

export type EventType = "card_created" | "card_edited";

export interface ActivityEvent {
  id: number;
  type: EventType;
  card_id: number;
  created_at: string;
  seen_at: string | null;
}

export type Grade = "again" | "hard" | "good" | "easy";
