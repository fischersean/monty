package reddit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

type authToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	DeviceId    string `json:"device_id"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

const (
	baseUrl = "oauth.reddit.com"
)

type Client struct {
	token authToken
	agent string
}

func NewClient(appId, appSecret, agent string) (c Client, err error) {

	c.agent = agent

	url := "https://www.reddit.com/api/v1/access_token?scope=read"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("grant_type", "https://oauth.reddit.com/grants/installed_client")
	_ = writer.WriteField("device_id", "DO_NOT_TRACK_THIS_DEVICE")
	err = writer.Close()
	if err != nil {
		return c, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return c, err
	}

	req.SetBasicAuth(appId, appSecret)
	req.Header.Set("User-agent", c.agent)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return c, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return c, err
	}

	token := authToken{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return c, err
	}

	if token.AccessToken == "" {
		return c, fmt.Errorf("could not retrieve token")
	}

	c.token = token
	return c, err
}

func (c *Client) oauthGet(u *url.URL) (resp *http.Response, err error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.token.TokenType, c.token.AccessToken))
	req.Header.Add("User-Agent", c.agent)

	resp, err = client.Do(req)

	// **** Commenting this out since reddit doesn't appear to be providing these headers anymore ****
	//rateLimitUsed, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Used"))
	//rateLimitRemaining, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	//rateLimitReset, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Reset"))

	//fmt.Printf("Used %d requests. %d remaining", rateLimitUsed, rateLimitRemaining)

	//if rateLimitRemaining == 0 {
	//time.Sleep(time.Duration(rateLimitReset) * time.Second)
	//}

	return resp, err
}

func (c *Client) get(path string, headers map[string]string) (resp *http.Response, err error) {

	u := &url.URL{
		Scheme: "https",
		Host:   baseUrl,
		Path:   path,
	}

	v := url.Values{}
	for k, h := range headers {
		v.Add(k, h)
	}
	u.RawQuery = v.Encode()

	return c.oauthGet(u)
}
