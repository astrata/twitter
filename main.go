/*
	Written by José Carlos Nieto <xiam@menteslibres.org>
	(c) 2012 Astrata Software http://astrata.mx

	MIT License
*/
package twitter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gosexy/sugar"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type multipartBody struct {
	ContentType string
	Data        io.Reader
}

/*
	API prefix.
*/
var Prefix = "https://api.twitter.com/1.1/"

/*
	If true, prints connection messages.
*/
var Debug = false

/*
	A Twitter client.
*/
type Client struct {
	client oauth.Client
	auth   *oauth.Credentials
}

/*
	Creates a new Twitter client.
*/
func New(credentials *oauth.Credentials) *Client {
	self := &Client{}
	self.client = oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
		Credentials:                   *credentials,
	}
	return self
}

/*
	Sets the User credentials for the Client.
*/
func (self *Client) SetAuth(credentials *oauth.Credentials) {
	self.auth = credentials
}

/*
	A helper for getting User credentials.
*/
func (self *Client) Setup() error {

	tmpCred, err := self.client.RequestTemporaryCredentials(http.DefaultClient, "", nil)

	if err != nil {
		return err
	}

	authURI := self.client.AuthorizationURL(tmpCred, nil)

	fmt.Printf("Hello, we are about to obtain your Twitter token.\n\n")
	fmt.Printf("Please open this URL in your browser:\n")
	fmt.Printf("%v\n\n", authURI)

	stdin := bufio.NewReader(os.Stdin)

	fmt.Printf("What's the PIN?\n")

	codebuf, err := stdin.ReadBytes('\n')

	if err != nil {
		return err
	}

	code := strings.Trim(string(codebuf), "\r\n ")

	auth, _, err := self.client.RequestToken(http.DefaultClient, tmpCred, code)

	if err != nil {
		return err
	}

	fmt.Printf("\nHere's your data:\n\n")

	self.SetAuth(auth)

	fmt.Printf("Token: %s\n", auth.Token)
	fmt.Printf("Secret: %s\n", auth.Secret)

	return nil
}

/*
	API 1.1 implementation

	https://dev.twitter.com/docs/api/1.1
*/

/*
	Returns an HTTP 200 OK response code and a representation of the requesting user if authentication was successful; returns a 401 status code and an error message if not. Use this method to test if supplied user credentials are valid.

	https://dev.twitter.com/docs/api/1.1/get/account/verify_credentials
*/
func (self *Client) VerifyCredentials(params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.get("/account/verify_credentials", merge(url.Values{}, params), data)
	return data, err
}

/*
	Returns a collection of the most recent Tweets and retweets posted by the authenticating user and the users they follow.

	https://dev.twitter.com/docs/api/1.1/get/statuses/home_timeline
*/
func (self *Client) HomeTimeline(params url.Values) (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get("/statuses/home_timeline", merge(url.Values{}, params), data)
	return data, err
}

/*
	Returns the 20 most recent mentions (tweets containing a users's @screen_name) for the authenticating user.

	https://dev.twitter.com/docs/api/1.1/get/statuses/mentions_timeline
*/
func (self *Client) MentionsTimeline(params url.Values) (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get("/statuses/mentions_timeline", merge(url.Values{}, params), data)
	return data, err
}

/*
	Returns a collection of the most recent Tweets posted by the user indicated by the screen_name or user_id parameters.

	https://dev.twitter.com/docs/api/1.1/get/statuses/user_timeline
*/
func (self *Client) UserTimeline(params url.Values) (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get("/statuses/user_timeline", merge(url.Values{}, params), data)
	return data, err
}

/*
	Returns the most recent tweets authored by the authenticating user that have recently been retweeted by others.

	https://dev.twitter.com/docs/api/1.1/get/statuses/retweets_of_me
*/
func (self *Client) RetweetsOfMe(params url.Values) (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get("/statuses/retweets_of_me", merge(url.Values{}, params), data)
	return data, err
}

/*
	Returns up to 100 of the first retweets of a given tweet.

	https://dev.twitter.com/docs/api/1.1/get/statuses/retweets/%3Aid
*/
func (self *Client) Retweets(id int64, params url.Values) (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get(fmt.Sprintf("/statuses/retweets/%d", id), merge(url.Values{}, params), data)
	return data, err
}

/*
	Returns a single Tweet, specified by the id parameter. The Tweet's author will also be embedded within the tweet.

	https://dev.twitter.com/docs/api/1.1/get/statuses/show/%3Aid
*/
func (self *Client) Show(id int64, params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.get(fmt.Sprintf("/statuses/show/%d", id), merge(url.Values{}, params), data)
	return data, err
}

/*
	Destroys the status specified by the required ID parameter.

	https://dev.twitter.com/docs/api/1.1/post/statuses/destroy/%3Aid
*/
func (self *Client) Destroy(id int64, params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.post(fmt.Sprintf("/statuses/destroy/%d", id), nil, merge(url.Values{}, params), data)
	return data, err
}

/*
	Retweets a tweet. Returns the original tweet with retweet details embedded.

	https://dev.twitter.com/docs/api/1.1/post/statuses/retweet/%3Aid
*/
func (self *Client) Retweet(id int64, params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.post(fmt.Sprintf("/statuses/retweet/%d", id), nil, merge(url.Values{}, params), data)
	return data, err
}

/*
	Updates the authenticating user's current status and attaches media for upload. In other words, it creates a Tweet with a picture attached.

	https://dev.twitter.com/docs/api/1.1/post/statuses/update_with_media
*/
func (self *Client) UpdateWithMedia(status string, params url.Values, files []string) (data *sugar.Map, err error) {

	endpoint := "/statuses/update_with_media"

	data = &sugar.Map{}

	buf := bytes.NewBuffer(nil)
	body := multipart.NewWriter(buf)

	for _, file := range files {

		writer, err := body.CreateFormFile("media[]", path.Base(file))

		if err != nil {
			return nil, err
		}

		reader, err := os.Open(file)

		if err != nil {
			return nil, err
		}

		io.Copy(writer, reader)

		reader.Close()
	}

	params = merge(url.Values{"status": {status}}, params)

	//fullURI := Prefix + strings.Trim(endpoint, "/") + ".json"

	//self.client.SignParam(self.auth, "POST", fullURI, params)

	for k, _ := range params {
		for _, value := range params[k] {
			body.WriteField(k, value)
		}
	}

	body.Close()

	//fmt.Printf("%v\n", buf)

	req := &multipartBody{body.FormDataContentType(), buf}

	err = self.request("POST", endpoint, nil, nil, req, data)

	return data, err
}

/*
	Updates the authenticating user's current status, also known as tweeting.

	https://dev.twitter.com/docs/api/1.1/post/statuses/update
*/
func (self *Client) Update(status string, params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}

	local := url.Values{
		"status": {status},
	}

	err = self.post("/statuses/update", nil, merge(local, params), data)

	return data, err
}

/*
	Returns information allowing the creation of an embedded representation of a Tweet on third party sites.

	https://dev.twitter.com/docs/api/1.1/get/statuses/oembed
*/
func (self *Client) Oembed(params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.get("/statuses/oembed", params, data)
	return data, err
}

/*
	Returns a collection of relevant Tweets matching a specified query.

	https://dev.twitter.com/docs/api/1.1/get/search/tweets
*/
func (self *Client) Search(params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.get("/search/tweets", params, data)
	return data, err
}

/*
	Returns a cursored collection of user IDs for every user the specified user is following (otherwise known as their "friends")

	https://dev.twitter.com/docs/api/1.1/get/friends/ids
*/
func (self *Client) Friends(params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.get("/friends/ids", params, data)
	return data, err
}

/*
	Returns a cursored collection of user IDs for every user following the specified user.

	https://dev.twitter.com/docs/api/1.1/get/followers/ids
*/
func (self *Client) Followers(params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.get("/followers/ids", params, data)
	return data, err
}

/*
	Returns fully-hydrated user objects for up to 100 users per request, as specified by comma-separated values passed to the user_id and/or screen_name parameters.

	https://dev.twitter.com/docs/api/1.1/get/users/lookup
*/
func (self *Client) LookupUser(params url.Values) (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get("/users/lookup", params, data)
	return data, err
}

/*
	Returns a variety of information about the user specified by the required user_id or screen_name parameter.

	https://dev.twitter.com/docs/api/1.1/get/users/show
*/
func (self *Client) ShowUser(params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.get("/users/show", params, data)
	return data, err
}

/*
	Connections.
*/

/*
	An HTTP request from the client to the Twitter API.
*/
func (self *Client) request(method string, endpoint string, getParams url.Values, postParams url.Values, buf *multipartBody, data interface{}) error {

	var requestURI string
	var res *http.Response
	var req *http.Request
	var err error

	if method != "GET" && method != "POST" {
		return fmt.Errorf("Unknown request method: %s\n", method)
	}

	if getParams == nil {
		getParams = url.Values{}
	}

	if postParams == nil {
		postParams = url.Values{}
	}

	fullURI := Prefix + strings.Trim(endpoint, "/") + ".json"

	if buf == nil {
		if method == "POST" {
			buf = &multipartBody{
				"application/x-www-form-urlencoded; charset=UTF-8",
				strings.NewReader(postParams.Encode()),
			}
		}
	}

	requestURI = fullURI + "?" + getParams.Encode()

	if buf == nil {
		req, err = http.NewRequest(method, requestURI, nil)
		if err != nil {
			return err
		}
	} else {
		req, err = http.NewRequest(method, requestURI, buf.Data)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", buf.ContentType)
	}

	addr, _ := url.Parse(requestURI)

	req.Header.Set("Authorization", self.client.AuthorizationHeader(self.auth, method, addr, postParams))

	client := &http.Client{}

	res, err = client.Do(req)

	if Debug == true {
		fmt.Printf("%s %s\n", method, requestURI)
	}

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		if Debug == true {
			fmt.Printf("Response: %s\n", body)
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("ERROR: %s returned status %d, %s", fullURI, res.StatusCode, body)
		}

		err = json.Unmarshal(body, data)

		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Something wrong happened.\n")
}

/*
	Shortcut for GET requests.
*/
func (self *Client) get(endpoint string, params url.Values, data interface{}) error {
	return self.request("GET", endpoint, params, nil, nil, data)
}

/*
	Shortcut for POST requests.
*/
func (self *Client) post(endpoint string, getParams url.Values, postParams url.Values, data interface{}) error {
	return self.request("POST", endpoint, getParams, postParams, nil, data)
}

/*
	Can merge default values with user provides ones.
*/
func merge(into url.Values, from url.Values) url.Values {
	if from != nil {
		for k, _ := range from {
			into.Set(k, from.Get(k))
		}
	}
	return into
}
