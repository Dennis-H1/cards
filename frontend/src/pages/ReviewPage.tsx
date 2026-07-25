import { useCallback, useEffect, useState } from "react";
import { gradeCard, getDueCards } from "../api/client";
import type { DueCard, Grade } from "../api/types";
import { CardView } from "../components/CardView";
import { formatIntervalDays, previewIntervalDays } from "../lib/sm2";

const GRADES: { grade: Grade; label: string; className: string }[] = [
  { grade: "again", label: "Again", className: "grade-again" },
  { grade: "hard", label: "Hard", className: "grade-hard" },
  { grade: "good", label: "Good", className: "grade-good" },
  { grade: "easy", label: "Easy", className: "grade-easy" },
];

const EMPTY_GRADE_COUNTS: Record<Grade, number> = { again: 0, hard: 0, good: 0, easy: 0 };

export function ReviewPage() {
  const [queue, setQueue] = useState<DueCard[] | null>(null);
  const [sessionTotal, setSessionTotal] = useState(0);
  const [showBack, setShowBack] = useState(false);
  const [gradeCounts, setGradeCounts] = useState(EMPTY_GRADE_COUNTS);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getDueCards()
      .then((cards) => {
        setQueue(cards);
        setSessionTotal(cards.length);
      })
      .catch((err) => setError(err.message));
  }, []);

  const grade = useCallback(
    async (g: Grade) => {
      if (!queue || queue.length === 0) return;
      const [current, ...rest] = queue;
      try {
        await gradeCard(current.card.id, g);
        setGradeCounts((counts) => ({ ...counts, [g]: counts[g] + 1 }));
        setQueue(rest);
        setShowBack(false);
      } catch (err) {
        setError(err instanceof Error ? err.message : String(err));
      }
    },
    [queue],
  );

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (!queue || queue.length === 0) return;
      if (e.key === " " || e.key === "Enter") {
        e.preventDefault();
        if (!showBack) setShowBack(true);
        return;
      }
      if (!showBack) return;
      const byKey: Record<string, Grade> = { "1": "again", "2": "hard", "3": "good", "4": "easy" };
      const g = byKey[e.key];
      if (g) grade(g);
    }
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [queue, showBack, grade]);

  if (error) {
    return <p className="empty-state">{error}</p>;
  }

  if (queue === null) {
    return <p className="empty-state">Loading…</p>;
  }

  if (queue.length === 0) {
    if (sessionTotal === 0) {
      return (
        <>
          <p className="due-count">0 due</p>
          <p className="empty-state">Nothing due. Go make some cards.</p>
        </>
      );
    }
    return (
      <div className="session-summary">
        <p className="session-summary__headline">
          Reviewed {sessionTotal} {sessionTotal === 1 ? "card" : "cards"}
        </p>
        <ul className="session-summary__breakdown">
          {GRADES.map(({ grade: g, label }) => (
            <li key={g}>
              {label}: {gradeCounts[g]}
            </li>
          ))}
        </ul>
      </div>
    );
  }

  const newCount = queue.filter((c) => c.review.review_count === 0).length;
  const dueCount = queue.length - newCount;
  const reviewedCount = sessionTotal - queue.length;
  const progressPercent = sessionTotal === 0 ? 0 : (reviewedCount / sessionTotal) * 100;
  const current = queue[0];

  return (
    <>
      <p className="due-count">
        {dueCount} due · {newCount} new
      </p>
      <div className="progress-bar">
        <div className="progress-bar__fill" style={{ width: `${progressPercent}%` }} />
      </div>
      <p className="progress-label">
        Card {reviewedCount + 1} of {sessionTotal}
      </p>

      <CardView key={current.card.id} card={current.card} showBack={showBack} />

      {showBack ? (
        <div className="grade-buttons">
          {GRADES.map(({ grade: g, label, className }) => (
            <button key={g} className={className} onClick={() => grade(g)}>
              {label}
              <span className="grade-buttons__interval">
                {formatIntervalDays(previewIntervalDays(current.review, g))}
              </span>
            </button>
          ))}
        </div>
      ) : (
        <button className="reveal-button" onClick={() => setShowBack(true)}>
          Show answer
        </button>
      )}
    </>
  );
}
