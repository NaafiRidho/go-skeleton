package controllers

import (
	"net/http"
	"user-service/common/response"
	"user-service/domain/dto"
	"user-service/services"
	"user-service/common/error"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	userService services.IServiceRegistry
}

type IUserController interface {
	Login(*gin.Context)
	Register(*gin.Context)
	Update(*gin.Context)
	GetUserLogin(*gin.Context)
	GetUserByUUID(*gin.Context)
}

func NewUserController(userService services.IServiceRegistry) IUserController {
	return &UserController{userService: userService}
}

func (c *UserController) Login(ctx *gin.Context) {
	request := &dto.LoginRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)

	if err != nil {
		errMessages := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := error.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessages,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	user, err := c.userService.GetUser().Login(ctx, request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: user.User,
		Token: &user.Token,
		Gin:  ctx,
	})
}

func (c *UserController) Register(ctx *gin.Context) {
	request := &dto.RegisterRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)

	if err != nil {
		errMessages := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := error.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessages,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	user, err := c.userService.GetUser().Register(ctx, request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: user.User,
		Gin:  ctx,
	})
}

func (c *UserController) Update(ctx *gin.Context) {
	request := &dto.UpdateUserRequest{}
	uuid := ctx.Param("uuid")

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)

	if err != nil {
		errMessages := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := error.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessages,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	user, err := c.userService.GetUser().Update(ctx, request, uuid)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (c *UserController) GetUserLogin(ctx *gin.Context) {
	user, err := c.userService.GetUser().GetUserLogin(ctx.Request.Context())
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (c *UserController) GetUserByUUID(ctx *gin.Context) {
	user, err := c.userService.GetUser().GetUserByUUID(ctx.Request.Context(), ctx.Param("uuid"))
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}
