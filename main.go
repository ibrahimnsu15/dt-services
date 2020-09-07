package main

import (
	"strconv"

	"dt-services/config"
	"dt-services/services"
	"dt-services/stores"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	config.InitConf()
	stores.InitDbs()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	initRoutes(e)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(services.PORT)))
}

func initRoutes(e *echo.Echo) {
	e.GET("/health", services.HealthHandler)

	e.GET("/github/orgs", services.GetGithubOrganizations)
	e.GET("/github/repos", services.GetGithubRepos)
	e.GET("/github/repos/:name", services.GetGithubRepoByName)
	e.GET("/github/repos/:name/commits", services.GetGithubRepoCommitsByName)
	e.GET("/github/repos/:name/commits/:sha", services.GetGithubRepoCommitBySha)
	e.GET("/github/ingest", services.Ingest)

	e.GET("/repos", services.GetReposHandler)
	e.GET("/repos/:id", services.GetRepoByIdHandler)
	e.GET("/repos/:id/commits", services.GetCommitsHandler)
	e.GET("/repos/:id/commits/:commitId", services.GetCommitByIdHandler)

	e.POST("/projects", services.PostProjectHandler)
	e.POST("/projects/:id/repos", services.PostProjectReposHandler)
	e.GET("/projects", services.GetProjectHandler)
	e.GET("/projects/:id", services.GetProjectByIdHandler)
	e.GET("/projects/:id/statss", services.GetProjectStatsByIdHandler)
	e.DELETE("/projects/:id", services.DeleteProjectByIdHandler)

	e.GET("/stats", services.GetStatsHandler)
	e.GET("/projects/:id/stats", services.GetStatsByProjectHandler)
	e.GET("/projects/:id/repos/:repoId/stats", services.GetStatsByRepoHandler)
}
