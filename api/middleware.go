package api

import (
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
