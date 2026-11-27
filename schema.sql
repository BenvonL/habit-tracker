-- Database schema for the habit tracker (MySQL / MariaDB).
--
-- Fresh install:
--   mysql -u <user> -p <database> < schema.sql
--
-- Existing database: see the migration section at the bottom of this file.

CREATE TABLE IF NOT EXISTS habit (
    habit_id       INT UNSIGNED NOT NULL AUTO_INCREMENT,
    habit_name     VARCHAR(255) NOT NULL,
    habit_unit     VARCHAR(255) NOT NULL,
    habit_pomodoro BOOLEAN      NOT NULL DEFAULT FALSE,
    habit_why      TEXT         NOT NULL,
    PRIMARY KEY (habit_id),
    UNIQUE KEY uq_habit_name (habit_name)
);

-- One row per habit and day. date_count holds the sum of everything logged
-- for that habit on that day.
CREATE TABLE IF NOT EXISTS date (
    date_id    INT UNSIGNED NOT NULL AUTO_INCREMENT,
    date_date  DATE         NOT NULL,
    date_count INT UNSIGNED NOT NULL DEFAULT 0,
    habit_id   INT UNSIGNED NOT NULL,
    PRIMARY KEY (date_id),
    -- required for the atomic upsert in AddCount (INSERT ... ON DUPLICATE KEY UPDATE)
    UNIQUE KEY uq_date_habit_day (habit_id, date_date),
    CONSTRAINT fk_date_habit
        FOREIGN KEY (habit_id) REFERENCES habit (habit_id)
        ON DELETE CASCADE
);

-- ---------------------------------------------------------------------------
-- Migration for an existing database
-- ---------------------------------------------------------------------------
-- Run these statements once, in order, if your tables already exist.
-- Back up the database first.
--
-- 1. Remove entries that belong to habits that no longer exist:
--
--    DELETE d FROM date d LEFT JOIN habit h ON h.habit_id = d.habit_id
--    WHERE h.habit_id IS NULL;
--
-- 2. Merge duplicate rows for the same habit and day into one row:
--
--    UPDATE date d
--    JOIN (SELECT habit_id, date_date, MIN(date_id) AS keep_id, SUM(date_count) AS total
--          FROM date GROUP BY habit_id, date_date HAVING COUNT(*) > 1) dup
--      ON d.date_id = dup.keep_id
--    SET d.date_count = dup.total;
--
--    DELETE d FROM date d
--    JOIN (SELECT habit_id, date_date, MIN(date_id) AS keep_id
--          FROM date GROUP BY habit_id, date_date HAVING COUNT(*) > 1) dup
--      ON d.habit_id = dup.habit_id AND d.date_date = dup.date_date
--     AND d.date_id <> dup.keep_id;
--
-- 3. Align the id column types, then add the unique key and the foreign key:
--
--    ALTER TABLE habit MODIFY habit_id INT UNSIGNED NOT NULL AUTO_INCREMENT;
--    ALTER TABLE date MODIFY habit_id INT UNSIGNED NOT NULL;
--    ALTER TABLE date ADD UNIQUE KEY uq_date_habit_day (habit_id, date_date);
--    ALTER TABLE date ADD CONSTRAINT fk_date_habit
--        FOREIGN KEY (habit_id) REFERENCES habit (habit_id) ON DELETE CASCADE;
