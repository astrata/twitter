/*
	Written by José Carlos Nieto <xiam@menteslibres.org>
	(c) 2012 Astrata Software http://astrata.mx

	MIT License
*/
package twitter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gosexy/sugar"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

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
func (self *Client) VerifyCredentials() (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.get("/account/verify_credentials", nil, data)
	return data, err
}

/*
	Returns a collection of the most recent Tweets and retweets posted by the authenticating user and the users they follow.

	https://dev.twitter.com/docs/api/1.1/get/statuses/home_timeline
*/
func (self *Client) HomeTimeline() (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get("/statuses/home_timeline", nil, data)
	return data, err
}

/*
	Returns the 20 most recent mentions (tweets containing a users's @screen_name) for the authenticating user.

	https://dev.twitter.com/docs/api/1.1/get/statuses/mentions_timeline
*/
func (self *Client) MentionsTimeline() (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get("/statuses/mentions_timeline", nil, data)
	return data, err
}

/*
	Returns a collection of the most recent Tweets posted by the user indicated by the screen_name or user_id parameters.

	https://dev.twitter.com/docs/api/1.1/get/statuses/user_timeline
*/
func (self *Client) UserTimeline() (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get("/statuses/user_timeline", nil, data)
	return data, err
}

/*
	Returns the most recent tweets authored by the authenticating user that have recently been retweeted by others.

	https://dev.twitter.com/docs/api/1.1/get/statuses/retweets_of_me
*/
func (self *Client) RetweetsOfMe() (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get("/statuses/retweets_of_me", nil, data)
	return data, err
}

/*
	Returns up to 100 of the first retweets of a given tweet.

	https://dev.twitter.com/docs/api/1.1/get/statuses/retweets/%3Aid
*/
func (self *Client) Retweets(id int64) (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.get(fmt.Sprintf("/statuses/retweets/%d", id), nil, data)
	return data, err
}

/*
	Returns a single Tweet, specified by the id parameter. The Tweet's author will also be embedded within the tweet.

	https://dev.twitter.com/docs/api/1.1/get/statuses/show/%3Aid
*/
func (self *Client) Show(id int64) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.get(fmt.Sprintf("/statuses/show/%d", id), nil, data)
	return data, err
}

/*
	Destroys the status specified by the required ID parameter.

	https://dev.twitter.com/docs/api/1.1/post/statuses/destroy/%3Aid
*/
func (self *Client) Destroy(id int64) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.post(fmt.Sprintf("/statuses/destroy/%d", id), nil, nil, data)
	return data, err
}

/*
	Retweets a tweet. Returns the original tweet with retweet details embedded.

	https://dev.twitter.com/docs/api/1.1/post/statuses/retweet/%3Aid
*/
func (self *Client) Retweet(id int64) (data *sugar.Map, err error) {
	data = &sugar.Map{}
	err = self.post(fmt.Sprintf("/statuses/retweet/%d", id), nil, nil, data)
	return data, err
}

/*
	Updates the authenticating user's current status and attaches media for upload. In other words, it creates a Tweet with a picture attached.

	https://dev.twitter.com/docs/api/1.1/post/statuses/update_with_media
*/
func (self *Client) UpdateWithMedia(status string, params url.Values) (data *sugar.Map, err error) {
	data = &sugar.Map{}

	if params == nil {
		params = url.Values{}
	}

	params.Add("status", status)

	err = self.post("/statuses/update_with_media", nil, params, data)

	return data, err
}

/*
	Updates the authenticating user's current status, also known as tweeting.

	https://dev.twitter.com/docs/api/1.1/post/statuses/update
*/
func (self *Client) Update(message string, params url.Values) (data *sugar.Map, err error) {

	if params == nil {
		params = url.Values{}
	}

	params.Add("status", message)

	data = &sugar.Map{}
	err = self.post("/statuses/update", nil, params, data)

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
func (self *Client) request(method string, endpoint string, getParams url.Values, postParams url.Values, data interface{}) error {

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

	var requestURI string
	var res *http.Response
	var err error

	if method == "GET" {
		self.client.SignParam(self.auth, method, fullURI, getParams)
		requestURI = fullURI + "?" + getParams.Encode()
		res, err = http.Get(requestURI)
	} else {
		self.client.SignParam(self.auth, method, fullURI, postParams)
		requestURI = fullURI
		res, err = http.PostForm(requestURI, postParams)
	}

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
	return self.request("GET", endpoint, params, nil, data)
}

/*
	Shortcut for POST requests.
*/
func (self *Client) post(endpoint string, getParams url.Values, postParams url.Values, data interface{}) error {
	return self.request("POST", endpoint, getParams, postParams, data)
}
