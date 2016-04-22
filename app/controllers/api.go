package controllers

import (
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"github.com/theothertomelliott/github-watchlists/app/models"
)

type Api struct {
	*revel.Controller
}

func init() {
	RegisterAuth(&Api{})
}

func (c Api) GetAuthenticatedUser() *models.User {
	return c.RenderArgs["user"].(*models.User)
}

type ApiResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Owner   string `json:"owner"`
	Repo    string `json:"repo"`
}

func NewApiResultSuccess(owner string, repo string) ApiResult {
	return ApiResult{true, "", owner, repo}
}

func (c Api) Unsubscribe(owner string, repo string, dryRun bool) revel.Result {
	u := c.GetAuthenticatedUser()
	if u == nil {
		c.Response.Status = 401
		return c.RenderJson(ApiResult{false, "Not authenticated", owner, repo})
	}

	client := GithubClientForUser(u)

	if repo == "" || owner == "" {
		return c.RenderJson(ApiResult{false, "Repo and repo owner is required", "owner", "repo"})
	}

	if !dryRun {
		_, err := client.Activity.DeleteRepositorySubscription(owner, repo)
		if err != nil {
			revel.ERROR.Println(err)
			return c.RenderJson(ApiResult{false, err.Error(), owner, repo})
		}
	}

	go cache.Delete("watched_" + u.AccessToken)

	return c.RenderJson(NewApiResultSuccess(owner, repo))
}
