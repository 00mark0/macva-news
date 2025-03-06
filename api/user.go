package api

import (
	"net/http"
	"time"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/token"
	"github.com/00mark0/macva-news/utils"

	"os"

	"github.com/labstack/echo/v4"
)

type userResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pfp      string `json:"pfp"`
	Role     string `json:"role"`
}

func newUserResponse(payload *token.Payload) userResponse {
	return userResponse{
		Username: payload.Username,
		Email:    payload.Email,
		Pfp:      payload.Pfp,
		Role:     payload.Role,
	}
}

type loginUserReq struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type loginUserRes struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) login(ctx echo.Context) error {
	var req loginUserReq
	var loginErr components.LoginErr
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	if req.Email == "" {
		loginErr = "Email obavezan"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	if req.Password == "" {
		loginErr = "Lozinka obavezna"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	user, err := server.store.GetUserByEmail(ctx.Request().Context(), req.Email)
	if err != nil {
		loginErr = "Nevažeći podaci za prijavu"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	if user.Banned.Bool {
		loginErr = "Nevažecí podaci za prijavu"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		loginErr = "Nevažecí podaci za prijavu"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	durationStr := os.Getenv("ACCESS_TOKEN_DURATION")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to parse duration", err))
		return err
	}

	accessToken, payload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Email,
		user.Pfp,
		user.Role,
		duration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to create access token", err))
		return err
	}

	res := loginUserRes{
		AccessToken: accessToken,
		User:        newUserResponse(payload),
	}

	if acceptHeader := ctx.Request().Header.Get("Accept"); acceptHeader == "application/json" {
		return ctx.JSON(http.StatusOK, res)
	}

	// Set token as a secure, HTTP-only cookie
	ctx.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(duration),
		Path:     "/",
		HttpOnly: true,
	})

	ctx.Response().Header().Set("HX-Redirect", "/")
	return ctx.NoContent(http.StatusOK)
}

func (server *Server) logOut(ctx echo.Context) error {
	ctx.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   "",
		MaxAge:  -1,
		Path:    "/",
		Expires: time.Now(),
	})

	ctx.Response().Header().Set("HX-Redirect", "/")
	return ctx.NoContent(http.StatusOK)
}
