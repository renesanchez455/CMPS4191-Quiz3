@@ -0,0 +1,4 @@
-- Filename: migrations/000002_add_todo_indexes.down.sql
DROP INDEX If EXISTS todo_name_idx;
DROP INDEX If EXISTS todo_priority_idx;