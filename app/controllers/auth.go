package controllers

import (
	"fmt"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/revel/revel"
	"github.com/theothertomelliott/github-watchlists/app/models"
)

func getConfig() *oauth2.Config {
	var GITHUB = &oauth2.Config{
		ClientID:     revel.Config.StringDefault("github.oauth.clientId", ""),
		ClientSecret: revel.Config.StringDefault("github.oauth.clientSecret", ""),
		Endpoint:     github.Endpoint,
		RedirectURL:  revel.Config.StringDefault("github.oauth.redirectUrl", ""),
		Scopes:       []string{"user", "repo"},
	}
	return GITHUB
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
	GITHUB := getConfig()

	u := c.GetAuthenticatedUser()
	if u != nil && u.AccessToken != "" {
		return c.Redirect(Application.Index)
	}

	authUrl := GITHUB.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Render(authUrl)
}

func (c Auth) Logout() revel.Result {
	u := c.GetAuthenticatedUser()
	if u != nil {
		c.Session["uid"] = ""
		c.RenderArgs["user"] = nil
	}

	return c.Redirect(Auth.Index)
}

func (c Auth) Auth(code string) revel.Result {
	GITHUB := getConfig()
	tok, err := GITHUB.Exchange(oauth2.NoContext, code)
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
