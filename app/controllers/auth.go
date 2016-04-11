package controllers

import (
	"fmt"
	"strconv"

	"code.google.com/p/goauth2/oauth"
	"github.com/revel/revel"
	"github.com/theothertomelliott/github-watchlists/app/models"
)

var GITHUB = &oauth.Config{
	ClientId:     "3fc8cf5e52ff07137a40",
	ClientSecret: "a4051132da4583860259cad737ca5666258c443a",
	AuthURL:      "https://github.com/login/oauth/authorize",
	TokenURL:     "https://github.com/login/oauth/access_token",
	RedirectURL:  "http://localhost:9000/Auth/Auth",
	Scope:        "user repo",
}

type Auth struct {
	*revel.Controller
}

func init() {
	RegisterAuth(&Auth{})
}

type AuthController interface {
	GetAuthenticatedUser() *models.User
}

func RegisterAuth(target AuthController) {
	revel.InterceptFunc(setUserOrNil, revel.BEFORE, target)
}

func (c Auth) Index() revel.Result {
	u := c.GetAuthenticatedUser()
	if u != nil && u.AccessToken != "" {
		// TODO: Do redirect based on input
		return c.Redirect(Application.Index)
	}

	// TODO: Make this a pure redirect, deal with rendering separately
	authUrl := GITHUB.AuthCodeURL("foo")
	return c.Render(authUrl)
}

func (c Auth) Auth(code string) revel.Result {
	t := &oauth.Transport{Config: GITHUB}
	tok, err := t.Exchange(code)
	if err != nil {
		revel.ERROR.Println(err)
		return c.Redirect(Auth.Index)
	}

	user := c.GetAuthenticatedUser()
	if user == nil {
		user = models.NewUser()
		c.Session["uid"] = fmt.Sprintf("%d", user.Uid)
		c.RenderArgs["user"] = user
	}
	user.AccessToken = tok.AccessToken
	return c.Redirect(Auth.Index)
}

func (c Auth) GetAuthenticatedUser() *models.User {
	return c.RenderArgs["user"].(*models.User)
}

func setUserOrNil(c *revel.Controller) revel.Result {
	var user *models.User = nil
	if _, ok := c.Session["uid"]; ok {
		uid, _ := strconv.ParseInt(c.Session["uid"], 10, 0)
		user = models.GetUser(int(uid))
	}
	c.RenderArgs["user"] = user
	return nil
}
