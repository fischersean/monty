package etl

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	db "github.com/fischersean/monty/internal/database"
	"github.com/fischersean/monty/pkg/reddit"
	"log"
)

type RedditEtlConfig struct {
	Subreddit string
	Client    *reddit.Client
	Limit     int
}

type AwsEtlConfig struct {
	Config aws.Config
}

type DatabaseEtlConfig struct {
	Connection *db.Connection
}

type EtlInputs struct {
	RunId          int
	RedditConfig   RedditEtlConfig
	AwsConfig      AwsEtlConfig
	DatabaseConfig DatabaseEtlConfig
}

func RunEtl(input *EtlInputs) error {
	log.Printf("Starting work for r/%s\n", input.RedditConfig.Subreddit)
	rconfig := input.RedditConfig
	data, err := pullSubredditData(rconfig.Client, rconfig.Subreddit, rconfig.Limit)
	if err != nil {
		return err
	}

	// if the subreddit is r/all look for new subs to add to out database
	if input.RedditConfig.Subreddit == "all" || input.RedditConfig.Subreddit == "popular" {
		log.Printf("r/%s: Searching for new subs\n", input.RedditConfig.Subreddit)
		// since this list is small we will just do a UPSERT
		// we will not run the pipeline on new entries. data will be added on the next triggered run
		for _, p := range data {
			err = input.DatabaseConfig.Connection.UpsertSubreddit(p.Subreddit)
			if err != nil {
				return err
			}
		}
	}

	//log.Printf("r/%s: Calculating post scores\n", input.RedditConfig.Subreddit)
	stats, err := calcSubredditStats(data)
	if err != nil {
		return err
	}
	//log.Printf("r/%s: Finished calculating scores\n", input.RedditConfig.Subreddit)

	//log.Printf("r/%s: Executing INSERT\n", input.RedditConfig.Subreddit)
	err = storeSubredditStats(input.DatabaseConfig.Connection, input.RunId, input.RedditConfig.Subreddit, stats)
	if err != nil {
		return err
	}
	log.Printf("✅ Successfully completed r/%s\n", input.RedditConfig.Subreddit)

	return err
}
