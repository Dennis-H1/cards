import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getCard, updateCard } from "../api/client";
import type { Card } from "../api/types";

export function CardDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [card, setCard] = useState<Card | null>(null);
  const [front, setFront] = useState("");
  const [back, setBack] = useState("");
  const [tagsInput, setTagsInput] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    if (!id) return;
    getCard(Number(id))
      .then((c) => {
        setCard(c);
        setFront(c.front);
        setBack(c.back);
        setTagsInput(c.tags.join(", "));
      })
      .catch((err) => setError(err.message));
  }, [id]);

  async function save() {
    if (!card) return;
    setSaving(true);
    setError(null);
    try {
      const tags = tagsInput
        .split(",")
        .map((t) => t.trim())
        .filter(Boolean);
      const updated = await updateCard(card.id, { front, back, tags });
      setCard(updated);
      setSaved(true);
      setTimeout(() => setSaved(false), 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setSaving(false);
    }
  }

  if (error && !card) return <p className="empty-state">{error}</p>;
  if (!card) return <p className="empty-state">Loading…</p>;

  const dirty = front !== card.front || back !== card.back || tagsInput !== card.tags.join(", ");

  return (
    <>
      <Link to="/browse">&larr; Back to Browse</Link>
      <h2 style={{ marginTop: "1rem" }}>Edit card</h2>

      <label htmlFor="card-front">Front</label>
      <textarea
        id="card-front"
        className="card-edit-field"
        rows={4}
        value={front}
        onChange={(e) => setFront(e.target.value)}
      />

      <label htmlFor="card-back">Back</label>
      <textarea
        id="card-back"
        className="card-edit-field"
        rows={4}
        value={back}
        onChange={(e) => setBack(e.target.value)}
      />

      <label htmlFor="card-tags">Tags (comma-separated)</label>
      <input
        id="card-tags"
        className="card-edit-field"
        value={tagsInput}
        onChange={(e) => setTagsInput(e.target.value)}
      />

      {error ? <p className="empty-state">{error}</p> : null}

      <div style={{ display: "flex", gap: "0.5rem", marginTop: "1rem" }}>
        <button className="reveal-button" disabled={!dirty || saving} onClick={save}>
          {saving ? "Saving…" : saved ? "Saved" : "Save"}
        </button>
        <button className="button-secondary" onClick={() => navigate(-1)}>
          Cancel
        </button>
      </div>
    </>
  );
}
