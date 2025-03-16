package api

import (
	"log"
	"os"
	"time"

	"github.com/00mark0/macva-news/token"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (server *Server) authMiddleware(tokenMaker token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cookie, err := ctx.Cookie("access_token")
			if err != nil {
				refreshCookie, err := ctx.Cookie("refresh_token")
				log.Println("Refresh cookie: ", refreshCookie)
				if err != nil {
					// For other requests, return an error
					log.Println("Error getting refresh cookie in adminMiddleware:", err)
					return ctx.NoContent(http.StatusNoContent)
				}
				refreshToken := refreshCookie.Value

				refreshPayload, err := tokenMaker.VerifyToken(refreshToken)
				if err != nil {
					log.Println("Error verifying refresh token:", err)
					return ctx.NoContent(http.StatusNoContent)
				}

				userIDStr := refreshPayload.UserID

				parsedUserID, err := uuid.Parse(userIDStr)
				if err != nil {
					log.Println("Invalid content ID format in archivePubContent:", err)
					return err
				}

				// Create a pgtype.UUID with the parsed UUID
				userID := pgtype.UUID{
					Bytes: parsedUserID,
					Valid: true,
				}

				user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
				if err != nil {
					log.Println("Error getting user by ID:", err)
					return ctx.NoContent(http.StatusNoContent)
				}

				if !user.EmailVerified.Bool {
					log.Println("User is not verified.")
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.Banned.Bool {
					log.Println("User is banned.")
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.IsDeleted.Bool {
					log.Println("User is deleted.")
					return ctx.NoContent(http.StatusNoContent)
				}

				accessTokenDurationStr := os.Getenv("ACCESS_TOKEN_DURATION")
				accessTokenDuration, err := time.ParseDuration(accessTokenDurationStr)
				if err != nil {
					log.Println("Error parsing access token duration:", err)
					return ctx.NoContent(http.StatusNoContent)
				}

				accessToken, accessTokenPayload, err := tokenMaker.CreateToken(
					refreshPayload.UserID,
					refreshPayload.Username,
					refreshPayload.Email,
					refreshPayload.Pfp,
					refreshPayload.Role,
					refreshPayload.EmailVerified,
					refreshPayload.Banned,
					refreshPayload.IsDeleted,
					accessTokenDuration,
				)
				if err != nil {
					log.Println("Error creating access token:", err)
					return ctx.NoContent(http.StatusNoContent)
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
					log.Println("Error verifying access token:", err)
					return ctx.NoContent(http.StatusNoContent)
				}

				userIDStr := payload.UserID

				parsedUserID, err := uuid.Parse(userIDStr)
				if err != nil {
					log.Println("Invalid content ID format in archivePubContent:", err)
					return err
				}

				// Create a pgtype.UUID with the parsed UUID
				userID := pgtype.UUID{
					Bytes: parsedUserID,
					Valid: true,
				}

				user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
				if err != nil {
					log.Println("Error getting user by ID:", err)
					return ctx.NoContent(http.StatusNoContent)
				}

				if !user.EmailVerified.Bool {
					log.Println("User is not verified.")
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.Banned.Bool {
					log.Println("User is banned.")
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.IsDeleted.Bool {
					log.Println("User is deleted.")
					return ctx.NoContent(http.StatusNoContent)
				}

				ctx.Set(authorizationPayloadKey, payload)
			}

			return next(ctx)

		}
	}
}

func (server *Server) adminMiddleware(tokenMaker token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cookie, err := ctx.Cookie("access_token")
			if err != nil {
				refreshCookie, err := ctx.Cookie("refresh_token")
				if err != nil {
					// For other requests, return an error
					log.Println("Error getting refresh cookie in adminMiddleware:", err)
					return ctx.NoContent(http.StatusNoContent)
				}
				refreshToken := refreshCookie.Value

				refreshPayload, err := tokenMaker.VerifyToken(refreshToken)
				if err != nil {
					log.Println("Error verifying refresh token:", err)
					return ctx.NoContent(http.StatusNoContent)
				}

				userIDStr := refreshPayload.UserID

				parsedUserID, err := uuid.Parse(userIDStr)
				if err != nil {
					log.Println("Invalid content ID format in archivePubContent:", err)
					return err
				}

				// Create a pgtype.UUID with the parsed UUID
				userID := pgtype.UUID{
					Bytes: parsedUserID,
					Valid: true,
				}

				user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
				if err != nil {
					log.Println("Error getting user by ID:", err)
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.Banned.Bool {
					log.Println("User is banned.")
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.IsDeleted.Bool {
					log.Println("User is deleted.")
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.Role != "admin" {
					log.Println("User is not admin.")
					return ctx.NoContent(http.StatusNoContent)
				}

				accessTokenDurationStr := os.Getenv("ACCESS_TOKEN_DURATION")
				accessTokenDuration, err := time.ParseDuration(accessTokenDurationStr)
				if err != nil {
					log.Println("Error parsing access token duration:", err)
					return ctx.NoContent(http.StatusNoContent)
				}

				accessToken, accessTokenPayload, err := tokenMaker.CreateToken(
					refreshPayload.UserID,
					refreshPayload.Username,
					refreshPayload.Email,
					refreshPayload.Pfp,
					refreshPayload.Role,
					refreshPayload.EmailVerified,
					refreshPayload.Banned,
					refreshPayload.IsDeleted,
					accessTokenDuration,
				)
				if err != nil {
					log.Println("Error creating access token:", err)
					return ctx.NoContent(http.StatusNoContent)
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
					log.Println("Error verifying access token:", err)
					// Invalid token; redirect to login page
					return ctx.NoContent(http.StatusNoContent)
				}

				userIDStr := payload.UserID

				parsedUserID, err := uuid.Parse(userIDStr)
				if err != nil {
					log.Println("Invalid content ID format in archivePubContent:", err)
					return err
				}

				// Create a pgtype.UUID with the parsed UUID
				userID := pgtype.UUID{
					Bytes: parsedUserID,
					Valid: true,
				}

				user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
				if err != nil {
					log.Println("Error getting user by ID:", err)
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.Banned.Bool {
					log.Println("User is banned.")
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.IsDeleted.Bool {
					log.Println("User is deleted.")
					return ctx.NoContent(http.StatusNoContent)
				}

				if user.Role != "admin" {
					log.Println("User is not admin.")
					return ctx.NoContent(http.StatusNoContent)
				}

				ctx.Set(authorizationPayloadKey, payload)
			}

			return next(ctx)

		}
	}
}
