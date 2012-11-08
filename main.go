package twitter

import (
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gosexy/sugar"
	"net/http"
	"os"
	"net/url"
	"strings"
	"io/ioutil"
	"encoding/json"
	"bufio"
	"fmt"
)

var Prefix = "https://api.twitter.com/1/"

type Client struct {
	client oauth.Client
	auth *oauth.Credentials
}

func New(credentials *oauth.Credentials) *Client {
	self := &Client{}
	self.client = oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
		TokenRequestURI: "https://api.twitter.com/oauth/access_token",
		Credentials: *credentials,
	}
	return self
}

func (self *Client) GetTimeline() (data *sugar.List, err error) {
	data = &sugar.List{}
	err = self.Request("/statuses/home_timeline", nil, data)
	return data, err
}

func (self *Client) Request(endpoint string, params url.Values, data interface{}) (error) {

	if params == nil {
		params = url.Values{}
	}

	fullURI := Prefix + endpoint + ".json"

	self.client.SignParam(self.auth, "GET", fullURI, params)

	res, err := http.Get(fullURI + "?" + params.Encode())

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

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

func (self *Client) SetAuth(credentials *oauth.Credentials) {
	self.auth = credentials
}

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
