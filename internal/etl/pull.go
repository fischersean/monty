package etl

import (
	"github.com/fischersean/monty/pkg/reddit"
)

// PullSubredditData fetches
func pullSubredditData(client *reddit.Client, subreddit string, limit int) (posts []reddit.Post, err error) {
	sub := client.Subreddit(subreddit, limit)
	posts, err = sub.GetHot()
	if err != nil {
		return posts, err
	}

	// get comments for each post
	for i := range posts {
		posts[i].Comments, err = posts[i].GetComments()
		if err != nil {
			return posts, err
		}
	}

	return posts, err
}
