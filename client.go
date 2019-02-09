package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Logger     *log.Logger
}

func New(logger *log.Logger) (*Client, error) {
	baseURL, err := url.Parse("https://qiita.com/api/v2")
	if err != nil {
		return nil, err
	}

	discardLogger := log.New(ioutil.Discard, "", log.LstdFlags)
	if logger == nil {
		logger = discardLogger
	}

	return &Client{
		URL:        baseURL,
		HTTPClient: http.DefaultClient,
		Logger:     logger,
	}, nil
}

func (c *Client) decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func (c *Client) newRequest(ctx context.Context, method string, relativePath string, body io.Reader) (*http.Request, error) {
	url := c.URL
	url.Path = path.Join(url.Path, relativePath)

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "qiita go-client (muiscript/qiita)")

	return req, nil
}

func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	req, err := c.newRequest(ctx, "GET", path.Join("users", userID), nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("user with id '%s' not found (status = %d)", userID, resp.StatusCode)
		default:
			return nil, fmt.Errorf("unknown error (status = %d)", resp.StatusCode)
		}
	}

	var user User
	if err := c.decodeBody(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) GetPost(ctx context.Context, postID string) (*Post, error) {
	req, err := c.newRequest(ctx, "GET", path.Join("items", postID), nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

    var post Post
    if err := c.decodeBody(resp, &post); err != nil {
    	return nil, err
	}

   	return &post, nil
}
