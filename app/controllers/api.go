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

func (c Api) Unsubscribe(repo string) revel.Result {
	u := c.GetAuthenticatedUser()
	if u == nil {
		c.Response.Status = 401
		return c.RenderJson(ApiResult{false, "Not authenticated"})
	}
	if repo == "" {
		return c.RenderJson(ApiResult{false, "Repo is required"})
	}
	return c.RenderJson(NewApiResultSuccess())
}
