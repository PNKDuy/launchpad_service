package controller

import (
	"github.com/labstack/echo/v4"
	"launchpad_service/model/token"
	"net/http"
)

// Create
// @Summary Create new token
// @Tags launchpad
// @Param model-value body object true "model-value"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {HTTPError} HTTPError
// @Router /launchpad/create [post]
func Create(c echo.Context) error {
	var tkn token.Token
	if err := c.Bind(&tkn); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tkn, err := token.Create(tkn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSONPretty(http.StatusOK, tkn, " ")
}

// Get
// @Summary Get activated token
// @Tags launchpad
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {HTTPError} HTTPError
// @Router /launchpad/get [get]
func Get(c echo.Context) error {
	tokens, err := token.Get()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSONPretty(http.StatusOK, tokens, " ")
}

// GetById
// @Summary Get token by id
// @Tags launchpad
// @Param id path string true "token-id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {HTTPError} HTTPError
// @Router /launchpad/get-by-id/{id} [get]
func GetById(c echo.Context) error {
	id := c.Param("id")
	tkn, err := token.GetById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSONPretty(http.StatusOK, tkn, " ")
}

// Update
// @Summary Update token
// @Tags launchpad
// @Param id path string true "token-id"
// @Param model_value body object true "model_value"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {HTTPError} HTTPError
// @Router /launchpad/update/{id} [put]
func Update(c echo.Context) error {
	id := c.Param("id")

	tkn, err := token.GetById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Bind(&tkn); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tkn, err = tkn.Update()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSONPretty(http.StatusOK, tkn, " ")
}

// DeactivateToken
// @Summary deactivate active token
// @Tags launchpad
// @Param id path string true "token-id"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {HTTPError} HTTPError
// @Router /launchpad/deactivate-token/{id} [put]
func DeactivateToken(c echo.Context) error {
	id := c.Param("id")

	tkn, err := token.GetById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Bind(&tkn); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	deactiveSuccess, err := tkn.DeactivateToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if deactiveSuccess == false {
		return c.String(200, "Delete failed")
	}

	return c.String(200, "Delete success")
}
