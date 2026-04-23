BEGIN;

ALTER TABLE books
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

ALTER TABLE members
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

ALTER TABLE borrowings
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_books_deleted_at
ON books(deleted_at);

CREATE INDEX IF NOT EXISTS idx_members_deleted_at
ON members(deleted_at);

CREATE INDEX IF NOT EXISTS idx_borrowings_deleted_at
ON borrowings(deleted_at);

ALTER TABLE books
DROP CONSTRAINT IF EXISTS books_isbn_key;

ALTER TABLE members
DROP CONSTRAINT IF EXISTS members_email_key;

CREATE UNIQUE INDEX IF NOT EXISTS books_isbn_active_unique
ON books(isbn)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS members_email_active_unique
ON members(email)
WHERE deleted_at IS NULL;

COMMIT;
