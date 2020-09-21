package services

import (
	"dt-services/models"
	"github.com/labstack/echo"
	"net/http"
)

func GetReposHandler(c echo.Context) error {
	repos, err := models.GetRepos()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, repos)
}

func GetRepoByIdHandler(c echo.Context) error {
	id := c.Param("id")
	repo, err := models.GetRepoById(id, true)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, repo)
}
