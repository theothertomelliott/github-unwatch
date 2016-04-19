package controllers

import (
	"github.com/revel/revel"
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
}

func NewApiResultSuccess() ApiResult {
	return ApiResult{true, ""}
}

func (c Api) Unsubscribe(owner string, repo string) revel.Result {
	u := c.GetAuthenticatedUser()
	if u == nil {
		c.Response.Status = 401
		return c.RenderJson(ApiResult{false, "Not authenticated"})
	}

	client := GithubClientForUser(u)

	if repo == "" || owner == "" {
		return c.RenderJson(ApiResult{false, "Repo and repo owner is required"})
	}

	_, err := client.Activity.DeleteRepositorySubscription(owner, repo)
	if err != nil {
		revel.ERROR.Println(err)
		return c.RenderJson(ApiResult{false, err.Error()})
	}

	return c.RenderJson(NewApiResultSuccess())
}
