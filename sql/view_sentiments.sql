SELECT 
    g.name,
    g.run_date,
    g.run_hour,
    g.toxic_index,
    g.stddev_toxic_index,
    g.stddev_toxic_index / g.toxic_index AS stdev_percent
FROM (
    SELECT
        f.name,
        f.run_date,
        f.run_hour,
        f.toxic_index,
        stddev(f.toxic_index) OVER (
            PARTITION BY f.name, f.run_date
        ) as stddev_toxic_index
    FROM (
        SELECT
            subreddits.name,
            watermarks.run_start::DATE as run_date,
            EXTRACT(HOUR FROM watermarks.run_start) as run_hour,
            count_comments,
            count_posts,
            score_compound_weighted_mean,
            score_compound_mean,
            score_compound_mean - score_compound_weighted_mean AS toxic_index
        FROM sentiment
        LEFT JOIN subreddits ON sentiment.subreddit_id=subreddits.id
        LEFT JOIN watermarks on sentiment.run_id=watermarks.id
    ) AS f
) as g
ORDER BY g.name, g.run_date ASC, g.run_hour ASC;