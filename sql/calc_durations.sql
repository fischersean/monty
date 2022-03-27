SELECT
    watermarks.run_start::DATE AS run_date
    , AVG(EXTRACT(EPOCH FROM (watermarks.run_end - watermarks.run_start))) / 60.0 as avg_duration
    , COUNT(DISTINCT EXTRACT(HOUR FROM watermarks.run_start)) as invocations
    --EXTRACT(EPOCH FROM (watermarks.run_end - watermarks.run_start)) AS duration
FROM sentiment
LEFT JOIN subreddits ON sentiment.subreddit_id=subreddits.id
LEFT JOIN watermarks ON sentiment.run_id=watermarks.id
GROUP BY run_date
ORDER BY run_date;