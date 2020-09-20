package services

import (
	"net/http"

	"dt-services/config"
	"dt-services/models"

	"github.com/labstack/echo"
)

func HealthHandler(c echo.Context) error {
	health := new(models.Health)
	health.Service = SERVICE_NAME
	health.Environment = config.Conf.Env
	health.Status = http.StatusOK

	return c.JSON(http.StatusOK, health)
}
