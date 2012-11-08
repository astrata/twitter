# Twitter 1.1 API for Golang

This is a work in progress. The final product aims to be a [Go][1] library for the Twitter [1.1][2] API.

## Want to help?

Register your APP at https://dev.twitter.com/apps to get your Key and Secret. Leave "Callback URL" blank, you don't need it for the command line tests.

Get the source with [git][3]:

```
# Creating a dev directory
% mkdir -p $GOPATH/src/github.com/astrata
% cd $GOPATH/src/github.com/astrata
% git clone git://github.com/astrata/twitter.git twitter
% cd twitter

# This may require you some additional work.
% go build

# This will ask you for your PIN and will give you your user credentials.
% twitter -key AAAA -secret BBBB

# Create a settings.yaml file with your app keys and user credentials.
% cat settings.yaml
twitter:
  app:
    key: ZerGYGhZytwFrsaR4xAse
    secret: PCadfTgdxAsercATs4Asr5dAx
  user:
    token: 12345678-rOaRx4saKTTuNJlhiuI7ehumzOV5xSp6dOtlk1Rs
    secret: fmt5pMcEbXer4DmmRFls7KesjXcQ4utgqrTf0KcR8

# Run the tests.
% go test

# Hack what you need and send pull requests :-).
vim main.go
```

[1]: http://golang.org
[2]: https://dev.twitter.com/docs/api/1.1
[3]: http://git-scm.com

