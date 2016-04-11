package controllers

import (
	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/revel/revel"
	"github.com/theothertomelliott/github-watchlists/app/models"
)

type Application struct {
	*revel.Controller
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
	var me *github.User
	var watched []github.Repository
	var watchedByLogin map[string][]github.Repository = make(map[string][]github.Repository)
	if u != nil {

		client := GithubClientForUser(u)

		var err error
		me, _, err = client.Users.Get("")
		if err != nil {
			revel.ERROR.Println(err)
		}

		watched, err = getAllWatchedByUser(client, "")
		if err != nil {
			revel.ERROR.Println(err)
		} else {
			for _, repo := range watched {
				login := repo.Owner.Login
				if arr, ok := watchedByLogin[*login]; ok {
					watchedByLogin[*login] = append(arr, repo)
				} else {
					watchedByLogin[*login] = []github.Repository{repo}
				}
			}
		}
		return c.Render(me, watched, watchedByLogin)
	}

	return c.Redirect(Auth.Index)
}

func (c Application) Unsubscribe() revel.Result {
	// TODO: Unsubscribe from specified repositories
	return c.Redirect(Application.Index)
}

func init() {
	RegisterAuth(&Application{})
}

func (c Application) GetAuthenticatedUser() *models.User {
	return c.RenderArgs["user"].(*models.User)
}
