import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { createCard } from "../api/client";

export function NewCardPage() {
  const navigate = useNavigate();
  const [front, setFront] = useState("");
  const [back, setBack] = useState("");
  const [tagsInput, setTagsInput] = useState("");
  const [source, setSource] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  async function save() {
    setSaving(true);
    setError(null);
    try {
      const tags = tagsInput
        .split(",")
        .map((t) => t.trim())
        .filter(Boolean);
      const card = await createCard({ front, back, tags, source: source || undefined });
      navigate(`/cards/${card.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setSaving(false);
    }
  }

  const valid = front.trim() !== "" && back.trim() !== "";

  return (
    <>
      <Link to="/browse">&larr; Back to Browse</Link>
      <h2 style={{ marginTop: "1rem" }}>New card</h2>

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

      <label htmlFor="card-source">Source (optional)</label>
      <input
        id="card-source"
        className="card-edit-field"
        value={source}
        onChange={(e) => setSource(e.target.value)}
      />

      {error ? <p className="empty-state">{error}</p> : null}

      <div style={{ display: "flex", gap: "0.5rem", marginTop: "1rem" }}>
        <button className="reveal-button" disabled={!valid || saving} onClick={save}>
          {saving ? "Saving…" : "Create"}
        </button>
        <button className="button-secondary" onClick={() => navigate(-1)}>
          Cancel
        </button>
      </div>
    </>
  );
}
