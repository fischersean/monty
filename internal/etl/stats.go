package etl

import (
	db "github.com/fischersean/monty/internal/database"
	"github.com/fischersean/monty/pkg/reddit"
	"github.com/jonreiter/govader"
)

// MAX_DEPTH is 3 to reduce total comments searched
const MAX_DEPTH = 3

type postStats struct {
	// countComments is the denominator for the normal mean
	countComments int
	// countPosts is the total number of posts processed
	countPosts int
	// cummulativeScore is the denominator for the weighted mean
	cummulativeScore int
	// this is the running numerator for normal and weighted mean compound scores
	rawCummulativeCompound         float64
	rawCummulativeWeightedCompound float64
	// weighted average compound score per comment by score
	scoreWeightedMeanCompoundScore float64
	// simple average compound score per comment
	meanCompoundScore float64
}

type subredditStats postStats

func calcCommentStats(analyzer *govader.SentimentIntensityAnalyzer, c reddit.Comment, ps *postStats, depth int) (s *postStats, err error) {
	if depth > MAX_DEPTH {
		return s, err
	}
	ps.countComments += 1
	// detect sentiment and add it to ps
	if len(c.Body) != 0 {
		sentiment := analyzer.PolarityScores(c.Body)
		if err != nil {
			return s, err
		}

		ps.rawCummulativeCompound += sentiment.Compound
		ps.rawCummulativeWeightedCompound += (float64(c.Score) * sentiment.Compound)
		ps.cummulativeScore += c.Score
	}

	for _, c := range c.Replies {
		s, err = calcCommentStats(analyzer, c, ps, depth+1)
	}

	return s, err
}

func calcPostStats(analyzer *govader.SentimentIntensityAnalyzer, post reddit.Post) (*postStats, error) {
	stats := &postStats{}
	for _, c := range post.Comments {
		stats, err := calcCommentStats(analyzer, c, stats, 0)
		if err != nil {
			return stats, err
		}
	}
	return stats, nil
}

func calcSubredditStats(data []reddit.Post) (stats subredditStats, err error) {
	// shared analzer
	analyzer := govader.NewSentimentIntensityAnalyzer()
	// process each post individually.
	for _, p := range data {
		postStats, err := calcPostStats(analyzer, p)
		if err != nil {
			return stats, err
		}
		stats.countPosts += 1
		stats.countComments += postStats.countComments
		stats.cummulativeScore += postStats.cummulativeScore
		stats.rawCummulativeCompound += postStats.rawCummulativeCompound
		stats.rawCummulativeWeightedCompound += postStats.rawCummulativeWeightedCompound
	}

	// finalize stats
	stats.meanCompoundScore = stats.rawCummulativeCompound / float64(stats.countComments)
	stats.scoreWeightedMeanCompoundScore = stats.rawCummulativeWeightedCompound / float64(stats.cummulativeScore)
	return stats, err
}

func storeSubredditStats(conn *db.Connection, runId int, subreddit string, stats subredditStats) (err error) {
	return conn.InsertSentimentResult(db.InsertSentimentResultInput{
		Subreddit:                 subreddit,
		CountPosts:                stats.countPosts,
		CountComments:             stats.countComments,
		ScoreCompoundMean:         stats.meanCompoundScore,
		ScoreCompoundWeightedMean: stats.scoreWeightedMeanCompoundScore,
		RunId:                     runId,
	})
}
