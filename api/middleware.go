package api

import (
	"fmt"
	"github.com/00mark0/macva-news/token"
	"strings"

	"net/http"

	"github.com/labstack/echo/v4"
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
					// No cookie found; redirect to login page
					return ctx.Redirect(http.StatusTemporaryRedirect, "/")
				}

				accessToken := cookie.Value
				payload, err := tokenMaker.VerifyToken(accessToken)
				if err != nil {
					// Invalid token; redirect to login page
					return ctx.Redirect(http.StatusTemporaryRedirect, "/")
				}

				ctx.Set(authorizationPayloadKey, payload)
			}

			return next(ctx)

		}
	}
}

func adminMiddleware(tokenMaker token.Maker) echo.MiddlewareFunc {
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

				if payload.Role != "admin" {
					ctx.JSON(http.StatusUnauthorized, errorResponse("not an admin", err))
					return err
				}

				ctx.Set(authorizationPayloadKey, payload)
			} else {
				// For HTML requests, get the token from the cookie "access_token"
				cookie, err := ctx.Cookie("access_token")
				if err != nil {
					// No cookie found; redirect to login page
					return ctx.Redirect(http.StatusTemporaryRedirect, "/")
				}

				accessToken := cookie.Value
				payload, err := tokenMaker.VerifyToken(accessToken)
				if err != nil {
					// Invalid token; redirect to login page
					return ctx.Redirect(http.StatusTemporaryRedirect, "/")
				}

				if payload.Role != "admin" {
					return ctx.Redirect(http.StatusSeeOther, "/")
				}

				ctx.Set(authorizationPayloadKey, payload)
			}

			return next(ctx)

		}
	}
}
