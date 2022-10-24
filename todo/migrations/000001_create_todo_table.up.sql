/*
	CMPS4191 - Quiz #3
	Rene Sanchez - 2018118383
*/
-- Filename: migrations/000001_create_todo_table.up.sql

CREATE TABLE IF NOT EXISTS todo (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    details text NOT NULL,
    priority text NOT NULL,
    status text NOT NULL,
    version integer NOT NULL DEFAULT 1
);