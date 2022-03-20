package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Post struct {
	Subreddit   string    `json:"subreddit"`
	Title       string    `json:"title"`
	Body        string    `json:"selftext"`
	Ups         int       `json:"ups"`
	Downs       int       `json:"downs"`
	Score       int       `json:"score"`
	UpvoteRatio float64   `json:"upvote_ratio"`
	Created     float64   `json:"created_utc"`
	Author      string    `json:"author"`
	Permalink   string    `json:"permalink"`
	NumComments int       `json:"num_comments"`
	Comments    []Comment `json:"comments"`
	c           *Client
}

type SubredditService struct {
	c     *Client
	Name  string
	Limit int
}

type postResponse struct {
	Data struct {
		Children []struct {
			Data Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func parsePostResponse(b []byte, c *Client) (posts []Post, err error) {

	rawResponse := postResponse{}
	err = json.Unmarshal(b, &rawResponse)
	if err != nil {
		return posts, err
	}

	for i, child := range rawResponse.Data.Children {
		posts = append(posts, child.Data)
		posts[i].c = c
	}

	return posts, err
}

func (c *Client) Subreddit(name string, limit int) *SubredditService {
	s := SubredditService{
		c:     c,
		Name:  name,
		Limit: limit,
	}
	return &s
}

func (s *SubredditService) getPosts(sort string) (posts []Post, er error) {

	validSorts := map[string]int{
		"top":    0,
		"hot":    0,
		"new":    0,
		"rising": 0,
	}

	if _, ok := validSorts[sort]; !ok {
		return posts, fmt.Errorf(fmt.Sprintf("Sort type not supported: %s", sort))
	}

	res, err := s.c.get(fmt.Sprintf("r/%s/%s.json", s.Name, sort), map[string]string{
		"limit": strconv.FormatInt(int64(s.Limit), 10),
	})
	if err != nil {
		return posts, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return posts, err
	}

	return parsePostResponse(b, s.c)
}

func (s *SubredditService) GetHot() ([]Post, error) {
	return s.getPosts("hot")
}

func (s *SubredditService) GetTop() ([]Post, error) {
	return s.getPosts("top")
}

func (s *SubredditService) GetNew() ([]Post, error) {
	return s.getPosts("new")
}

func (s *SubredditService) GetRising() ([]Post, error) {
	return s.getPosts("rising")
}
