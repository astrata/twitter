/*
	Written by José Carlos Nieto <xiam@menteslibres.org>
	(c) 2012 Astrata Software http://astrata.mx

	MIT License
*/

package twitter

import (
	//"fmt"
	"github.com/garyburd/go-oauth/oauth"
	//"github.com/gosexy/sugar"
	"github.com/gosexy/to"
	"github.com/gosexy/yaml"
	"testing"
	//"time"
)

var SettingsFile = "settings.yaml"

var conf *yaml.Yaml

func TestSettings(t *testing.T) {
	var err error
	conf, err = yaml.Open(SettingsFile)

	if err != nil {
		panic(err.Error())
	}

}

func TestApi(t *testing.T) {

	var err error

	client := New(&oauth.Credentials{
		to.String(conf.Get("twitter/app/key")),
		to.String(conf.Get("twitter/app/secret")),
	})

	client.SetAuth(&oauth.Credentials{
		to.String(conf.Get("twitter/user/token")),
		to.String(conf.Get("twitter/user/secret")),
	})

	_, err = client.VerifyCredentials(nil)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}
	/*

	_, err = client.HomeTimeline(nil)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}

	_, err = client.MentionsTimeline(nil)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}

	_, err = client.UserTimeline(nil)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}

	_, err = client.RetweetsOfMe(nil)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}

	_, err = client.Retweets(int64(21947795900469248), nil)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}

	var status *sugar.Map
	status, err = client.Update(fmt.Sprintf("Test message @ %s", time.Now()), nil)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}

	tweetId := to.Int64(status.Get("id_str"))

	_, err = client.Destroy(tweetId, nil)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}

	tweetId := to.Int64(status.Get("id_str"))
	*/

	files := []string{
		"_resources/test.jpg",
	}

	_, err = client.UpdateWithMedia("Hello", nil, files)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}

	/*
	_, err = client.Retweet(int64(21947795900469248), nil)

	if err != nil {
		t.Errorf("Test failed: %s\n", err.Error())
	}
	*/

}
