package api

import (
	"log"
	"net/http"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	"github.com/00mark0/macva-news/utils"
	"github.com/labstack/echo/v4"
)

func (server *Server) listContentComments(ctx echo.Context) error {
	var req ListAdsReq
	var userData db.GetUserByIDRow

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listContentComments:", err)
		return err
	}

	contentIDStr := ctx.Param("id")

	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in listContentComments:", err)
		return err
	}

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		log.Println("User is not logged in.")
	} else {
		payload, err := server.tokenMaker.VerifyToken(cookie.Value)

		userID, err := utils.ParseUUID(payload.UserID, "userID")
		if err != nil {
			log.Println("Error parsing user_id in homePage:", err)
			return err
		}

		user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
		if err != nil {
			log.Println("Error getting user in homePage:", err)
			return err
		}

		userData = user
	}

	nextLimit := req.Limit + 10

	arg := db.ListContentCommentsParams{
		ContentID: contentID,
		Limit:     nextLimit,
	}

	comments, err := server.store.ListContentComments(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error listing content comments:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleComments(contentIDStr, comments, userData))
}

type CreateCommentReq struct {
	CommentText string `form:"comment_text" validate:"required,min=3,max=1000"`
}

func (server *Server) createComment(ctx echo.Context) error {
	var req CreateCommentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createComment:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in createComment:", err)
		return err
	}

	var userData db.GetUserByIDRow
	contentIDStr := ctx.Param("id")

	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in createComment:", err)
		return err
	}

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		log.Println("Login to post a comment.")
		return err
	} else {
		payload, err := server.tokenMaker.VerifyToken(cookie.Value)

		userID, err := utils.ParseUUID(payload.UserID, "userID")
		if err != nil {
			log.Println("Error parsing user_id in homePage:", err)
			return err
		}

		user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
		if err != nil {
			log.Println("Error getting user in homePage:", err)
			return err
		}

		userData = user
	}

	_, err = server.store.CreateComment(ctx.Request().Context(), db.CreateCommentParams{
		ContentID:   contentID,
		UserID:      userData.UserID,
		CommentText: req.CommentText,
	})
	if err != nil {
		log.Println("Error creating comment:", err)
		return err
	}

	return server.listContentComments(ctx)
}
