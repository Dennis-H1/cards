import type { ActivityEvent, Card, DueCard, Grade, Review, Tag, TagOverview } from "./types";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`/api${path}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? `request failed: ${res.status}`);
  }
  if (res.status === 204) {
    return undefined as T;
  }
  return res.json();
}

export function getDueCards(): Promise<DueCard[]> {
  return request("/cards/due");
}

export function getCard(id: number): Promise<Card> {
  return request(`/cards/${id}`);
}

export function searchCards(query: string): Promise<Card[]> {
  return request(`/cards/search?q=${encodeURIComponent(query)}`);
}

export function createCard(input: {
  front: string;
  back: string;
  tags: string[];
  source?: string;
}): Promise<Card> {
  return request("/cards", { method: "POST", body: JSON.stringify(input) });
}

export function updateCard(
  id: number,
  input: { front?: string; back?: string; tags?: string[] },
): Promise<Card> {
  return request(`/cards/${id}`, { method: "PATCH", body: JSON.stringify(input) });
}

export function gradeCard(id: number, grade: Grade): Promise<Review> {
  return request(`/cards/${id}/grade`, { method: "POST", body: JSON.stringify({ grade }) });
}

export function listTags(): Promise<Tag[]> {
  return request("/tags");
}

export function getTagOverview(name: string): Promise<TagOverview> {
  return request(`/tags/${encodeURIComponent(name)}/overview`);
}

export function listActivity(unseenOnly = false): Promise<ActivityEvent[]> {
  return request(`/activity${unseenOnly ? "?unseen=true" : ""}`);
}

export function markActivitySeen(): Promise<void> {
  return request("/activity/seen", { method: "POST" });
}
