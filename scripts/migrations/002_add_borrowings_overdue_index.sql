BEGIN;

CREATE INDEX IF NOT EXISTS idx_borrowings_overdue
    ON borrowings(due_date)
    WHERE status = 'borrowed' AND due_date IS NOT NULL;

COMMIT;