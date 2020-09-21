package services

import (
	"context"
	"fmt"
	"net/http"

	"dt-services/config"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v32/github"
	"github.com/labstack/echo"
)

func GetGithubOrganizations(c echo.Context) error {
	organizations, err := fetchOrganizations(GITHUB_ORG)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, organizations)
}

func GetGithubRepos(c echo.Context) error {
	repos, err := fetchRepos(GITHUB_ORG)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, repos)
}

func GetGithubRepoByName(c echo.Context) error {
	name := c.Param("name")
	repo, err := fetchRepoByName(GITHUB_ORG, name)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, repo)
}

func GetGithubRepoCommitsByName(c echo.Context) error {
	name := c.Param("name")
	commits, err := fetchRepoCommitsByName(GITHUB_ORG, name)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, commits)
}

func GetGithubRepoCommitBySha(c echo.Context) error {
	name := c.Param("name")
	sha := c.Param("sha")
	commit, err := fetchRepoCommitBySha(GITHUB_ORG, name, sha)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, commit)
}

func fetchOrganizations(username string) ([]*github.Organization, error) {
	client, err := newGithubClient()
	if err != nil {
		return nil, err
	}

	orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func fetchRepos(username string) ([]*github.Repository, error) {
	client, err := newGithubClient()
	if err != nil {
		return nil, err
	}

	var allRepos []*github.Repository
	opts := &github.RepositoryListByOrgOptions{
		Type: "all",
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), username, opts)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func fetchRepoByName(username, name string) (*github.Repository, error) {
	client, err := newGithubClient()
	if err != nil {
		return nil, err
	}

	repo, _, err := client.Repositories.Get(context.Background(), username, name)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func fetchRepoCommitsByName(username, name string) ([]*github.RepositoryCommit, error) {
	client, err := newGithubClient()
	if err != nil {
		return nil, err
	}

	var allCommits []*github.RepositoryCommit
	//opts := &github.CommitsListOptions{
	//	Since: time.Now().Add((time.Hour * 24 * 7) * -1),
	//}
	opts := &github.CommitsListOptions{}
	for {
		commits, resp, err := client.Repositories.ListCommits(
			context.Background(),
			username,
			name,
			opts,
		)
		if err != nil {
			fmt.Println(err.Error())
			//return nil, err
		}
		allCommits = append(allCommits, commits...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allCommits, nil
}

func fetchRepoCommitBySha(username, name, sha string) (*github.RepositoryCommit, error) {
	client, err := newGithubClient()
	if err != nil {
		return nil, err
	}

	// todo:files do paginate on a length == 300
	// var allFiles []*github.RepositoryCommit
	commit, resp, err := client.Repositories.GetCommit(
		context.Background(),
		username,
		name,
		sha,
	)

	fmt.Println(resp.NextPage, len(commit.Files))

	if err != nil {
		fmt.Println(err.Error())
		//return nil, err
	}

	return commit, nil
}

func newGithubClient() (*github.Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport,
		57726,
		7369504,
		config.Conf.GithubAppKey,
	)

	if err != nil {
		return nil, err
	}
	client := github.NewClient(&http.Client{Transport: itr})

	return client, nil
}
