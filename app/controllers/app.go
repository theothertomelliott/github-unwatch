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

func getWatchedByUserAndLogin(client *github.Client, user string) (map[string][]github.Repository, error) {

	var err error
	var watched []github.Repository
	var watchedByLogin map[string][]github.Repository = make(map[string][]github.Repository)

	watched, err = getAllWatchedByUser(client, "")
	if err != nil {
		return watchedByLogin, err
	}
	for _, repo := range watched {
		login := repo.Owner.Login
		if arr, ok := watchedByLogin[*login]; ok {
			watchedByLogin[*login] = append(arr, repo)
		} else {
			watchedByLogin[*login] = []github.Repository{repo}
		}
	}

	return watchedByLogin, nil
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

	var watchedByLogin map[string][]github.Repository = make(map[string][]github.Repository)
	watchedByLogin, err = getWatchedByUserAndLogin(client, "")
	if err != nil {
		revel.ERROR.Println(err)
	}
	return c.Render(me, watchedByLogin)

}

func (c Application) Unsubscribe(login string) revel.Result {

	if login == "" {
		return c.Redirect(Application.Index)
	}

	u := c.GetAuthenticatedUser()
	if u == nil {
		return c.Redirect(Auth.Index)
	}

	var err error

	client := GithubClientForUser(u)

	// TODO: Cache these lists
	var watchedByLogin map[string][]github.Repository = make(map[string][]github.Repository)
	watchedByLogin, err = getWatchedByUserAndLogin(client, "")
	if err != nil {
		revel.ERROR.Println(err)
		return c.Redirect(Application.Index)
	}

	if repos, ok := watchedByLogin[login]; ok {
		for _, repo := range repos {
			_, err = client.Activity.DeleteRepositorySubscription(*repo.Owner.Login, *repo.Name)
			if err != nil {
				revel.ERROR.Println(err)
			}
		}
	} else {
		revel.INFO.Println("No repos found")
	}

	return c.Redirect(Application.Index)
}
