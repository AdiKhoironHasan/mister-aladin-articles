package http

import (
	"net/http"
	"os"
	"strconv"

	"github.com/AdiKhoironHasan/mister-aladin-articles/internal/services"
	servConst "github.com/AdiKhoironHasan/mister-aladin-articles/pkg/common/const"
	"github.com/AdiKhoironHasan/mister-aladin-articles/pkg/dto"
	servErrors "github.com/AdiKhoironHasan/mister-aladin-articles/pkg/errors"

	"github.com/apex/log"
	"github.com/labstack/echo"
)

type HttpHandler struct {
	service services.Services
}

func NewHttpHandler(e *echo.Echo, srv services.Services) {
	handler := &HttpHandler{
		srv,
	}
	e.GET("/ping", handler.Ping)
	e.POST("/articles", handler.CreateArticles)
	e.GET("/articles", handler.ShowArticles)
	e.GET("/articles/:id", handler.ShowArticlesByID)
	e.PUT("/articles/:id", handler.UpdateArticle)
	e.DELETE("/articles/:id", handler.DeleteArticle)
}

func (h *HttpHandler) Ping(c echo.Context) error {
	version := os.Getenv("VERSION")

	if version == "" {
		version = "pong"
	}

	data := version

	return c.JSON(http.StatusOK, data)

}

func (h *HttpHandler) CreateArticles(c echo.Context) error {
	postDTO := dto.ArticleReqDTO{}

	if err := c.Bind(&postDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}
	err := postDTO.Validate()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = h.service.CreateArticles(&postDTO)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.SaveSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) ShowArticles(c echo.Context) error {
	getDTO := dto.ArticleParamReqDTO{}

	if err := c.Bind(&getDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	err := getDTO.Validate()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	result, err := h.service.ShowArticles(&getDTO)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.GetDataSuccess,
		Data:    result,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) ShowArticlesByID(c echo.Context) error {
	var articleID int

	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	result, err := h.service.ShowArticlesByID(articleID)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.GetDataSuccess,
		Data:    result,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) UpdateArticle(c echo.Context) error {
	var articleID int
	putDTO := dto.ArticleReqDTO{}

	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Bind(&putDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	putDTO.ID = articleID

	err = putDTO.Validate()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = h.service.UpdateArticle(&putDTO)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.UpdateSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *HttpHandler) DeleteArticle(c echo.Context) error {
	var articleID int

	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	err = h.service.DeleteArticle(articleID)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.DeleteSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)

}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case servErrors.ErrInternalServerError:
		return http.StatusInternalServerError
	case servErrors.ErrNotFound:
		return http.StatusNotFound
	case servErrors.ErrConflict:
		return http.StatusConflict
	case servErrors.ErrInvalidRequest:
		return http.StatusBadRequest
	case servErrors.ErrFailAuth:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
