SELECT 
    g.name,
    g.run_date,
    AVG(g.toxic_index) AS avg_toxic_index,
    MAX(g.stddev_toxic_index) AS stddev_toxic_index
FROM (
    SELECT
        f.name,
        f.run_date,
        f.run_hour,
        f.toxic_index,
        stddev(f.toxic_index) OVER (
            PARTITION BY f.name, f.run_date
        ) AS stddev_toxic_index
    FROM (
        SELECT
            subreddits.name,
            watermarks.run_start::DATE AS run_date,
            EXTRACT(HOUR FROM watermarks.run_start) as run_hour,
            count_comments,
            count_posts,
            score_compound_weighted_mean,
            score_compound_mean,
            score_compound_mean - score_compound_weighted_mean AS toxic_index
        FROM sentiment
        LEFT JOIN subreddits ON sentiment.subreddit_id=subreddits.id
        LEFT JOIN watermarks ON sentiment.run_id=watermarks.id
    ) AS f
) as g
GROUP BY g.name, g.run_date
ORDER BY avg_toxic_index DESC;