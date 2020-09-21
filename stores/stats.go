package services

import (
	"net/http"

	"dt-services/models"

	"github.com/labstack/echo"
)

func GetStatsHandler(c echo.Context) error {
	stats, err := models.ParseGlobalStats()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, stats)
}

func GetStatsByProjectHandler(c echo.Context) error {
	id := c.Param("id")
	project, err := models.GetProjectById(id, true)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	stats, err := models.ParseProjectStat(project, nil, nil)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, stats)
}

func GetStatsByRepoHandler(c echo.Context) error {
	id := c.Param("id")
	repoId := c.Param("repoId")
	project, err := models.GetProjectById(id, true)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	found := false
	for _, rId := range project.RepoIds {
		if rId.Hex() == repoId {
			found = true
			break
		}
	}
	if !found {
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, err)
	}

	repo, err := models.GetRepoById(repoId, false)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	stats, err := models.ParseRepoStat(repo, nil, nil)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, stats)
}
