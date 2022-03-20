SELECT COUNT(*) FROM sentiment;

SELECT 
    id, 
    run_start,
    run_end,
    successful,
    (run_end - run_start) AS duration
FROM 
    watermarks
ORDER BY
    run_start DESC
LIMIT 5;

SELECT 
    subreddits.name,
    run_id,
    count_comments,
    count_posts,
    score_compound_weighted_mean,
    score_compound_mean,
    ROUND(score_compound_mean - score_compound_weighted_mean, 2) AS toxic_index
FROM 
    sentiment
LEFT JOIN 
    subreddits ON sentiment.subreddit_id=subreddits.id
ORDER BY 
    toxic_index DESC
LIMIT 10;

-- What is the toal comments per run?
SELECT 
    SUM(count_comments), run_id
FROM 
    sentiment
GROUP BY
    run_id
ORDER BY
    run_id ASC;