package twitter

import (
	"testing"
	"github.com/gosexy/yaml"
	"fmt"
	"github.com/gosexy/to"
	"github.com/garyburd/go-oauth/oauth"
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

	client := New(&oauth.Credentials{
		to.String(conf.Get("twitter/app/key")),
		to.String(conf.Get("twitter/app/secret")),
	})

	client.SetAuth(&oauth.Credentials{
		to.String(conf.Get("twitter/user/token")),
		to.String(conf.Get("twitter/user/secret")),
	})

	data, err := client.GetTimeline()

	if err != nil {
		fmt.Printf("Test failed.\n")
	}

	fmt.Printf("%v\n", data)
}
