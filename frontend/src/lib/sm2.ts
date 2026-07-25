import type { Grade, Review } from "../api/types";

// Preview-only port of internal/service/sm2.go's applySM2 -- used to show
// "this grade -> N days" hints under the grading buttons before the user
// commits. Not authoritative: the real update happens server-side via
// POST /cards/:id/grade. Keep in sync with the Go implementation by hand.
//
// Note: per the spec, only Again forks the interval calculation -- Hard/
// Good/Easy all produce the same interval for a given repetitions count and
// only diverge in the ease-factor update, which affects future reviews, not
// this one. So Hard/Good/Easy previews are expected to show the same value.
export function previewIntervalDays(review: Review, grade: Grade): number {
  if (grade === "again") return 1;
  if (review.repetitions === 0) return 1;
  if (review.repetitions === 1) return 6;
  return Math.round(review.interval_days * review.ease_factor);
}

export function formatIntervalDays(days: number): string {
  if (days < 1) return "<1d";
  if (days < 30) return `${days}d`;
  if (days < 365) return `${Math.round(days / 30)}mo`;
  return `${Math.round(days / 365)}y`;
}
