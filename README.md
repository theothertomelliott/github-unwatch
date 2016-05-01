# GitHub Bulk Unwatcher


If you've joined an organization and discovered that you're getting e-mails from all of their repos, unwatching them all can be a chore. This tool allows quick and easy bulk unwatching of repositories.

## Usage

Go to [https://github-unwatch.herokuapp.com/](https://github-unwatch.herokuapp.com/) and sign in using your GitHub account.

You will be presented with a list of all the Git repositories you're currently watching.

Filter this list by entering a query in the search field, your list will be whittled down as you type, until you've found exactly which repos you want to unwatch.

Press "Unwatch" to remove these repositories from your watchlist. Just like that, no more floods of e-mails!

## Implementation

GitHub Bulk Unwatcher is built using:

* [The Revel Framework](https://revel.github.io)
* [Vue.js](https://vuejs.org)
* [Bootstrap](http://getbootstrap.com)
* [FontAwesome](https://fortawesome.github.io/Font-Awesome/)

I can highly recommend these for your next quick web project!

GitHub integration uses:

* [The GitHub API](https://developer.github.com/v3/)
* [go-github](https://github.com/google/go-github)

## Deployment

The app is deployed to [Heroku](https://dashboard.heroku.com/), using a [modified Revel Buildpack](https://github.com/theothertomelliott/heroku-buildpack-go-revel).

A [Dockerfile](https://www.docker.com/) is also provided, with a Makefile for running the app locally if desired:

* `make container` Builds the image
* `make up` Launches a container with the built image
* `make down` Stops and deletes running containter
* `make clean` Deletes the image

This was mostly an experiment to get used to Docker and isn't used to deploy.