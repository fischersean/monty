package main

import (
	"log"
	"os"
	"runtime"
	"time"

	db "github.com/fischersean/monty/internal/database"
	"github.com/fischersean/monty/internal/etl"
	_ "github.com/fischersean/monty/internal/tzinit"
	"github.com/fischersean/monty/pkg/reddit"
)

const (
	POST_LIMIT = 50
)

type Result struct {
	Error error
	Sub   string
}

type Job struct {
	RunId        int
	Sub          string
	Conn         *db.Connection
	RedditClient *reddit.Client
}

func worker(jobs <-chan Job, results chan<- Result) {
	for j := range jobs {
		err := etl.RunEtl(&etl.EtlInputs{
			RunId: j.RunId,
			RedditConfig: etl.RedditEtlConfig{
				Subreddit: j.Sub,
				Limit:     POST_LIMIT,
				Client:    j.RedditClient,
			},
			DatabaseConfig: etl.DatabaseEtlConfig{
				Connection: j.Conn,
			},
		})
		results <- Result{
			Sub:   j.Sub,
			Error: err,
		}
	}
}

func main() {
	// get db info from secrets vault
	log.Println("Establishing database connection")
	dbInfo, err := getDbConn()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.NewConnection(db.NewConnectionInput{
		Host:     dbInfo.Host,
		Port:     dbInfo.Port,
		User:     dbInfo.Username,
		Password: dbInfo.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Database connection OK")

	// generate new run id
	runId, err := conn.GenerateRunID()
	if err != nil {
		log.Fatal(err)
	}

	// reddit setup
	appId := os.Getenv("APP_ID")
	appSecret := os.Getenv("APP_SECRET")
	appAgent := os.Getenv("APP_AGENT")

	log.Println("Generating reddit API token")
	redditClient, err := reddit.NewClient(appId, appSecret, appAgent)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Reddit API OK")

	successful := true
	subs, err := conn.GetAllSubreddits()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(subs)

	// setup worker pool
	nJobs := len(subs)
	jobs := make(chan Job, nJobs)
	results := make(chan Result, nJobs)

	// deploy workers
	nWorkers := runtime.NumCPU() * 4 // 4x since this pipeline is heavily IO limited, not CPU
	for i := 0; i < nWorkers; i++ {
		go worker(jobs, results)
	}

	// load up jobs with data
	for i := 0; i < nJobs; i++ {
		jobs <- Job{
			RunId:        runId,
			Sub:          subs[i],
			Conn:         &conn,
			RedditClient: &redditClient,
		}
	}
	close(jobs)

	// listen for results
	for i := 0; i < nJobs; i++ {
		res := <-results
		if res.Error != nil {
			log.Printf("❌ Error while processing r/%s: %s", res.Error, res.Error)
			successful = false
		}
	}
	close(results)

	// finish the run by updating the watermark table
	err = conn.UpdateWatermark(db.UpdateWatermarkInput{
		Id:         runId,
		RunEnd:     time.Now(),
		Successful: successful,
	})
	if err != nil {
		log.Fatal(err)
	}
}
