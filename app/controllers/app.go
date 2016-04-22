package controllers

import (
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"github.com/theothertomelliott/github-watchlists/app/models"
)

type Application struct {
	*revel.Controller
}

func init() {
	RegisterAuth(&Application{})
}

func (c Application) GetAuthenticatedUser() *models.User {
	return c.RenderArgs["user"].(*models.User)
}

func getAllWatchedByUser(client *github.Client, user string) ([]github.Repository, error) {

	var response *github.Response = nil
	var reposOut []github.Repository
	page := 1
	for response == nil || response.LastPage != 0 {
		var repos []github.Repository
		var err error
		repos, response, err = client.Activity.ListWatched(user, &github.ListOptions{Page: page, PerPage: 100})
		if err != nil {
			return nil, err
		}
		page = response.NextPage

		reposOut = append(reposOut, repos...)
	}

	return reposOut, nil
}

func GithubClientForUser(user *models.User) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc)
}

func (c Application) Index() revel.Result {
	u := c.GetAuthenticatedUser()
	if u == nil {
		return c.Redirect(Auth.Index)
	}

	client := GithubClientForUser(u)

	var err error
	var me *github.User
	me, _, err = client.Users.Get("")
	if err != nil {
		revel.ERROR.Println(err)
	}

	var watched []github.Repository
	if err = cache.Get("watched_"+u.AccessToken, &watched); err != nil {
		watched, err = getAllWatchedByUser(client, "")
		if err != nil {
			revel.ERROR.Println(err)
		}

		go cache.Set("watched_"+u.AccessToken, watched, 60*time.Minute)
	}
	return c.Render(me, watched)

}
