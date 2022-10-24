@@ -0,0 +1,4 @@
/*
	CMPS4191 - Quiz #3
	Rene Sanchez - 2018118383
*/
-- Filename: migrations/000002_add_todo_indexes.down.sql
DROP INDEX If EXISTS todo_name_idx;
DROP INDEX If EXISTS todo_priority_idx;