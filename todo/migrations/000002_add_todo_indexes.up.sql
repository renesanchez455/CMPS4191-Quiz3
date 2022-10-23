-- Filename: migrations/000002_add_todo_indexes.up.sql
CREATE INDEX IF NOT EXISTS todo_name_idx ON todo USING GIN(to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS todo_priority_idx ON todo USING GIN(to_tsvector('simple', priority));
