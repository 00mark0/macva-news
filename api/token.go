package api

/*
this will renew the session when called from the client on clients that use json

import (
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type renewAccessRequest struct {
	RefreshToken string `json:"refresh_token" form:"username" validate:"required"`
}

type renewAccessResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccess(ctx echo.Context) error {
	var req renewAccessRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse("invalid refresh token", err))
		return err
	}

	session, err := server.store.GetSession(ctx.Request().Context(), pgtype.UUID{Bytes: refreshPayload.ID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse("failed to get user", err))
		return err
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse("blocked session", err))
		return err
	}

	if session.Username != refreshPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse("invalid refresh token", err))
		return err
	}

	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, errorResponse("invalid refresh token", err))
		return err
	}

	if time.Now().After(session.ExpiresAt.Time) {
		ctx.JSON(http.StatusUnauthorized, errorResponse("expired refresh token", err))
		return err
	}

	accessTokenDurationStr := os.Getenv("ACCESS_TOKEN_DURATION")

	accessTokenDuration, err := time.ParseDuration(accessTokenDurationStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to parse duration", err))
		return err
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		accessTokenDuration,
	)

	res := renewAccessResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	return ctx.JSON(http.StatusOK, res)

}*/
