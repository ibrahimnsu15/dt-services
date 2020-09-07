package services

import (
	"dt-services/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetCommitsHandler(c echo.Context) error {
	id := c.Param("id")
	repo, err := models.GetRepoById(id, true)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	commits, err := models.GetCommitsByRepo(repo, nil, nil)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, commits)
}

func GetCommitByIdHandler(c echo.Context) error {
	//id := c.Param("id")
	commitId := c.Param("commitId")
	commit, err := models.GetCommitById(commitId, true)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, commit)
}
