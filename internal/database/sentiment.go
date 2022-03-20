package database

type InsertSentimentResultInput struct {
	RunId                     int
	Subreddit                 string
	CountComments             int
	CountPosts                int
	ScoreCompoundMean         float64
	ScoreCompoundWeightedMean float64
}

func (conn *Connection) InsertSentimentResult(input InsertSentimentResultInput) (err error) {
	stmt := `
INSERT INTO sentiment (subreddit_id, run_id, count_comments, count_posts, score_compound_weighted_mean, score_compound_mean)
VALUES (
    (SELECT 
        FIRST_VALUE(id) OVER (
            ORDER BY id
        )
    FROM 
        subreddits 
    WHERE name=$1), 
	$2,
    $3, 
    $4,
    $5, 
    $6
);
`
	_, err = conn.db.Exec(stmt, input.Subreddit, input.RunId, input.CountComments, input.CountPosts, input.ScoreCompoundWeightedMean, input.ScoreCompoundMean)
	return err
}
