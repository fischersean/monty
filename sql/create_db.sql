CREATE DATABASE monty;

-- Delete everything
DROP TABLE IF EXISTS monty.sentiment;
DROP TABLE IF EXISTS monty.watermarks;
DROP TABLE IF EXISTS monty.subreddits;


-- Create the subreddits table
CREATE TABLE IF NOT EXISTS monty.subreddits (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Watermark table
CREATE TABLE IF NOT EXISTS monty.watermarks (
    id SERIAL PRIMARY KEY,
    run_start TIMESTAMP,
    run_end TIMESTAMP,
    successful BOOLEAN
);

-- Sentiment table
CREATE TABLE IF NOT EXISTS monty.sentiment (
    id SERIAL PRIMARY KEY,
    subreddit_id INTEGER REFERENCES subreddits(id) NOT NULL,
    run_id INTEGER REFERENCES watermarks(id) NOT NULL,
    count_comments INTEGER NOT NULL,
    count_posts INTEGER NOT NULL,
    score_compound_weighted_mean DECIMAL NOT NULL,
    score_compound_mean DECIMAL NOT NULL
);

-- Seed db with some initial values
INSERT INTO monty.subreddits (name) VALUES ('all');
INSERT INTO monty.subreddits (name) VALUES ('popular');