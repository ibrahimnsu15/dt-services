package services

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"dt-services/models"

	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

func PostProjectHandler(c echo.Context) error {
	type reqBody struct {
		Name string `json:"name"`
	}
	var req reqBody
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	project := models.NewProject(req.Name)
	if err := project.Save(); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, project)
}

func PostProjectReposHandler(c echo.Context) error {
	id := c.Param("id")
	project, err := models.GetProjectById(id, true)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	reqBody := map[string]interface{}{}
	if err := c.Bind(&reqBody); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var repoIds []bson.ObjectId
	for k, v := range reqBody {
		if k == "repo_ids" {
			repoIdsIs := v.([]interface{})
			for _, repoIdI := range repoIdsIs {
				repoIds = append(repoIds, bson.ObjectIdHex(repoIdI.(string)))
			}
		}
	}

	changed := false
	for _, repoId := range repoIds {
		found := false
		if len(project.RepoIds) == 0 {
			changed = true
			project.RepoIds = append(project.RepoIds, repoId)
			continue
		} else {
			for _, pRepoId := range project.RepoIds {
				if repoId == pRepoId {
					found = true
					break
				}
			}

			if !found {
				changed = true
				project.RepoIds = append(project.RepoIds, repoId)
				break
			}
		}
	}

	if changed == true {
		if err := project.Update(); err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, project)
}

func GetProjectHandler(c echo.Context) error {
	projects, err := models.GetProjects()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, projects)
}

func GetProjectByIdHandler(c echo.Context) error {
	id := c.Param("id")
	project, err := models.GetProjectById(id, true)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, project)
}

func DeleteProjectByIdHandler(c echo.Context) error {
	id := c.Param("id")
	if err := models.DeleteProjectById(id); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func GetProjectStatsByIdHandler(c echo.Context) (err error) {
	id := c.Param("id")
	start := strings.TrimSuffix(c.QueryParam("start"), "000")
	end := strings.TrimSuffix(c.QueryParam("end"), "000")
	period := c.QueryParam("period")

	var periodQuery string
	if period == PERIOD_DAILY {
		periodQuery = PERIOD_DAILY
	} else if period == PERIOD_WEEKLY {
		periodQuery = PERIOD_WEEKLY
	} else if period == PERIOD_MONTHLY {
		periodQuery = PERIOD_MONTHLY
	} else if period == "" {
		periodQuery = PERIOD_DAILY
	}

	var startTimestamp, endTimestamp *time.Time
	if start != "" {
		startI, err := strconv.Atoi(start)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		timestamp := time.Unix(int64(startI), 0)
		startTimestamp = &timestamp
	}

	if end != "" {
		endI, err := strconv.Atoi(end)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		timestamp := time.Unix(int64(endI), 0)
		endTimestamp = &timestamp
	}

	project, err := models.GetProjectById(id, true)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	projectStats, err := models.ParseProjectStat(project, startTimestamp, endTimestamp)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	type repo struct {
		Name string         `json:"name"`
		Data map[string]int `json:"data"`
	}

	type resp struct {
		Repos []repo `json:"repos"`
	}

	statsResp := new(resp)
	for _, r := range projectStats.RepoStats {
		repo := new(repo)
		repo.Name = r.Name
		var commitsS []string
		for _, c := range r.CommitStats {
			date := getDate(c.CommittedAt)
			period := getPeriod(date, periodQuery)
			commitsS = append(commitsS, period)
		}

		dataM := make(map[string]int)
		for _, c := range commitsS {
			if len(dataM) == 0 {
				dataM[c] = 1
				continue
			}

			found := false
			foundK := ""
			for k, _ := range dataM {
				if c == k {
					found = true
					foundK = k
					break
				}
			}

			if found {
				dataM[foundK] = dataM[foundK] + 1
			} else {
				dataM[c] = 1
			}
		}

		repo.Data = dataM
		statsResp.Repos = append(statsResp.Repos, *repo)
	}

	return c.JSON(http.StatusOK, statsResp)
}
