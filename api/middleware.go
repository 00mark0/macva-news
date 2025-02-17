package api

import (
	"fmt"
	"github.com/00mark0/macva-news/token"
	"strings"

	"net/http"

	"github.com/labstack/echo/v4"
	"os"
	"time"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			acceptHeader := ctx.Request().Header.Get("Accept")

			if strings.Contains(acceptHeader, "application/json") {
				// For JSON requests, get the token from the header "authorization"

				authorizationHeader := ctx.Request().Header.Get(authorizationHeaderKey)

				if len(authorizationHeader) == 0 {
					err := echo.ErrUnauthorized
					ctx.JSON(http.StatusUnauthorized, errorResponse("authorization header is not provided", err))
					return err
				}

				fields := strings.Fields(authorizationHeader)
				if len(fields) < 2 {
					err := echo.ErrUnauthorized
					ctx.JSON(http.StatusUnauthorized, errorResponse("invalid authorization header format", err))
					return err
				}

				authorizationType := strings.ToLower(fields[0])
				if authorizationType != authorizationTypeBearer {
					err := fmt.Errorf("Type: %s", authorizationType)
					ctx.JSON(http.StatusUnauthorized, errorResponse("unsupported authorization type", err))
					return err
				}

				accessToken := fields[1]
				payload, err := tokenMaker.VerifyToken(accessToken)
				if err != nil {
					ctx.JSON(http.StatusUnauthorized, errorResponse("invalid access token", err))
					return err
				}

				ctx.Set(authorizationPayloadKey, payload)
			} else {
				// For HTML requests, get the token from the cookie "access_token"
				cookie, err := ctx.Cookie("access_token")
				if err != nil {
					refreshCookie, err := ctx.Cookie("refresh_token")
					fmt.Println("Refresh cookie: ", refreshCookie)
					if err != nil {
						// For other requests, return an error
						err := echo.ErrUnauthorized
						return ctx.JSON(http.StatusUnauthorized, errorResponse("missing refresh token", err))
					}
					refreshToken := refreshCookie.Value

					refreshPayload, err := tokenMaker.VerifyToken(refreshToken)
					if err != nil {
						return ctx.JSON(http.StatusUnauthorized, errorResponse("invalid refresh token", err))
					}

					accessTokenDurationStr := os.Getenv("ACCESS_TOKEN_DURATION")
					accessTokenDuration, err := time.ParseDuration(accessTokenDurationStr)
					if err != nil {
						return ctx.JSON(http.StatusInternalServerError, errorResponse("failed to parse duration", err))
					}

					accessToken, accessTokenPayload, err := tokenMaker.CreateToken(
						refreshPayload.Username,
						accessTokenDuration,
					)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, errorResponse("failed to create access token", err))
						return err
					}

					ctx.SetCookie(&http.Cookie{
						Name:     "access_token",
						Value:    accessToken,
						Path:     "/",
						HttpOnly: true,
						Secure:   false,
						Expires:  time.Now().Add(accessTokenDuration),
					})

					ctx.Set(authorizationPayloadKey, accessTokenPayload)
				} else {
					accessToken := cookie.Value

					payload, err := tokenMaker.VerifyToken(accessToken)
					if err != nil {
						return ctx.JSON(http.StatusUnauthorized, errorResponse("invalid access token", err))
					}

					ctx.Set(authorizationPayloadKey, payload)
				}
			}

			return next(ctx)

		}
	}
}
