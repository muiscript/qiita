package main

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strconv"
)

// User represents a qiita user
type User struct {
	ID           string `json:"id"`
	PermanentID  int    `json:"permanent_id"`
	ImageURL     string `json:"profile_image_url"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	Description  string `json:"description"`
	WebsiteURL   string `json:"website_url"`
	Organization string `json:"organization"`
	TeamOnly     bool   `json:"team_only"`

	PostsCount     int `json:"items_count"`
	FolloweesCount int `json:"followees_count"`
	FollowersCount int `json:"followers_count"`

	GithubID   string `json:"github_login_name"`
	LinkedinID string `json:"linkedin_id"`
	TwitterID  string `json:"twitter_screen_name"`
}

// UsersResponse represents a response of qiita API which returns multiple users
type UsersResponse struct {
	Users      []*User
	PerPage    int
	Page       int
	FirstPage  int
	LastPage   int
	TotalCount int
}

type paginationInfo struct {
	PerPage    int
	Page       int
	FirstPage  int
	LastPage   int
	TotalCount int
}

// GetUser fetches a user having provided userID
//
// GET /api/v2/users/:user_id
// document: https://qiita.com/api/v2/docs#get-apiv2usersuser_id
func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	req, err := c.newRequest(ctx, "GET", path.Join("users", userID), nil, nil)
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

// GetUsers fetches all the users
// The number of users included in one response and the page number should be provided
//
// GET /api/v2/users
// document: https://qiita.com/api/v2/docs#get-apiv2users
func (c *Client) GetUsers(ctx context.Context, page int, perPage int) (*UsersResponse, error) {
	if err := c.validatePaginationLimit(page, 1, 100, perPage, 1, 100); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, "GET", "users", query, nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	var users []*User
	if err := c.decodeBody(resp, &users); err != nil {
		return nil, err
	}

	paginationInfo, err := c.extractPaginationInfo(resp, page, perPage)
	if err != nil {
		return nil, err
	}

	return constructUsersResponse(users, paginationInfo), nil
}

// GetFollowees fetches all the followees of the user having provided userID
// The number of users included in one response and the page number should be provided
//
// GET /api/v2/users/:user_id/followees
// document: http://qiita.com/api/v2/docs#get-apiv2usersuser_idfollowees
func (c *Client) GetFollowees(ctx context.Context, userID string, page int, perPage int) (*UsersResponse, error) {
	if err := c.validatePaginationLimit(page, 1, 100, perPage, 1, 100); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, "GET", path.Join("users", userID, "followees"), query, nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	var users []*User
	if err := c.decodeBody(resp, &users); err != nil {
		return nil, err
	}

	paginationInfo, err := c.extractPaginationInfo(resp, page, perPage)
	if err != nil {
		return nil, err
	}

	return constructUsersResponse(users, paginationInfo), nil
}

// GetFollowers fetches all the followers of the user having provided userID
// The number of users included in one response and the page number should be provided
//
// GET /api/v2/users/:user_id/followers
// document: https://qiita.com/api/v2/docs#get-apiv2usersuser_idfollowers
func (c *Client) GetFollowers(ctx context.Context, userID string, page int, perPage int) (*UsersResponse, error) {
	if err := c.validatePaginationLimit(page, 1, 100, perPage, 1, 100); err != nil {
		return nil, err
	}

	query := map[string]string{
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
	}
	req, err := c.newRequest(ctx, "GET", path.Join("users", userID, "followers"), query, nil)
	if err != nil {
		return nil, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	var users []*User
	if err := c.decodeBody(resp, &users); err != nil {
		return nil, err
	}

	paginationInfo, err := c.extractPaginationInfo(resp, page, perPage)
	if err != nil {
		return nil, err
	}

	return constructUsersResponse(users, paginationInfo), nil
}

// IsFollowingUser returns true if the authenticated user is following the user having provided userID
// This method requires authentication
//
// GET /api/v2/users/:user_id/following
// document: https://qiita.com/api/v2/docs#get-apiv2usersuser_idfollowing
func (c *Client) IsFollowingUser(ctx context.Context, userID string) (bool, error) {
	req, err := c.newRequest(ctx, "GET", path.Join("users", userID, "following"), nil, nil)
	if err != nil {
		return false, err
	}
	c.Logger.Printf("send get request to %s\n", c.URL.String())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return false, fmt.Errorf("unauthorized. you may have provided no/invalid access token (status = %d)", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else {
		return false, nil
	}
}

func constructUsersResponse(users []*User, info *paginationInfo) *UsersResponse {
	return &UsersResponse{
		Users:      users,
		PerPage:    info.PerPage,
		Page:       info.Page,
		FirstPage:  info.FirstPage,
		LastPage:   info.LastPage,
		TotalCount: info.TotalCount,
	}
}