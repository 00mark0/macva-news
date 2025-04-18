package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/redis"
	"github.com/00mark0/macva-news/db/services"
	"github.com/00mark0/macva-news/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (server *Server) listContentComments(ctx echo.Context) error {
	var req ListAdsReq
	var userData db.GetUserByIDRow
	var userReactions map[string]string

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

	var userID pgtype.UUID
	cookie, err := ctx.Cookie("refresh_token")
	if err == nil {
		payload, err := server.tokenMaker.VerifyToken(cookie.Value)
		if err == nil {
			parsedUserID, err := utils.ParseUUID(payload.UserID, "userID")
			if err == nil {
				userID = parsedUserID

				user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
				if err == nil {
					userData = user

					userReactions, err = server.getUserReactionsForContentWithCache(ctx.Request().Context(), contentID, userID)
					if err != nil {
						log.Println("Error fetching user reactions:", err)
					}
				} else {
					log.Println("Error getting user in listContentComments:", err)
				}
			} else {
				log.Println("Error parsing user_id in listContentComments:", err)
			}
		} else {
			log.Println("Error verifying token in listContentComments:", err)
		}
	} else {
		log.Println("User is not logged in.")
		userReactions = make(map[string]string) // still need an empty map for rendering
	}

	nextLimit := req.Limit + 10

	comments, err := server.getCommentsWithCache(ctx.Request().Context(), contentID, nextLimit)
	if err != nil {
		log.Println("Error listing comments in listContentComments:", err)
		return err
	}

	url := fmt.Sprintf("/api/content/comments/%s?limit=", contentIDStr)

	commentCount, err := server.getCommentCountWithCache(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting comment count in listContentComments:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleComments(
		contentIDStr,
		comments,
		userData,
		userReactions,
		int(nextLimit),
		url,
		int(commentCount),
	))
}

type CreateCommentReq struct {
	CommentText string `form:"comment_text" validate:"required,min=3,max=10000"`
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

	// Invalidate cache for this content's comments
	pattern := redis.GenerateKey("comments", contentID, "*")
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), pattern)
	if err != nil {
		log.Printf("Error invalidating comments cache: %v", err)
		// Continue despite cache invalidation error
	} else {
		log.Printf("Invalidated cache for content ID: %s", contentID)
	}

	commentCountKey := redis.GenerateKey("comment_count", contentID)
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), commentCountKey)
	if err != nil {
		log.Printf("Error invalidating comment count cache: %v", err)
		// Continue despite cache invalidation error
	}

	return server.listContentComments(ctx)
}

// Handler for upvoting a comment
func (server *Server) handleUpvoteComment(ctx echo.Context) error {
	var userData db.GetUserByIDRow
	commentIDStr := ctx.Param("id")

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		log.Println("User is not logged in.")
	}

	payload, err := server.tokenMaker.VerifyToken(cookie.Value)
	if err != nil {
		log.Println("Invalid token:", err)
		return err
	}

	userID, err := utils.ParseUUID(payload.UserID, "userID")
	if err != nil {
		log.Println("Error parsing user_id in handleDownvoteComment:", err)
		return err
	}

	user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error getting user in handleDownvoteComment:", err)
		return err
	}
	userData = user

	commentID, err := utils.ParseUUID(commentIDStr, "comment ID")
	if err != nil {
		log.Println("Invalid comment ID format in handleUpvoteComment:", err)
		return err
	}

	// Check if the user already has a reaction
	userReaction, err := server.store.GetUserCommentReaction(ctx.Request().Context(), db.GetUserCommentReactionParams{
		CommentID: commentID,
		UserID:    userData.UserID,
	})

	// Handle reaction logic based on whether we found a reaction and what it was
	if err == nil {
		// User has an existing reaction
		if userReaction.Reaction == "like" {
			// If already liked, remove the reaction
			_, err = server.store.DeleteCommentReaction(ctx.Request().Context(), db.DeleteCommentReactionParams{
				CommentID: commentID,
				UserID:    userData.UserID,
			})
			if err != nil {
				log.Println("Error deleting comment reaction from like to remove like:", err)
				return err
			}
		} else if userReaction.Reaction == "dislike" {
			// If disliked, change to like
			_, err := server.store.InsertOrUpdateCommentReaction(ctx.Request().Context(), db.InsertOrUpdateCommentReactionParams{
				CommentID: commentID,
				UserID:    userData.UserID,
				Reaction:  "like",
			})
			if err != nil {
				log.Println("Error changing reaction from dislike to like:", err)
				return err
			}
		}
	} else {
		// No reaction yet, add a like
		_, err := server.store.InsertOrUpdateCommentReaction(ctx.Request().Context(), db.InsertOrUpdateCommentReactionParams{
			CommentID: commentID,
			UserID:    userData.UserID,
			Reaction:  "like",
		})
		if err != nil {
			log.Println("Error adding new like reaction:", err)
			return err
		}
	}

	// Update the comment's score
	updatedComment, err := server.store.UpdateCommentScore(ctx.Request().Context(), commentID)
	if err != nil {
		log.Println("Error updating comment score:", err)
		return err
	}

	// Get the updated user reaction for the response
	reactionStatus := ""
	updatedUserReaction, err := server.store.GetUserCommentReaction(ctx.Request().Context(), db.GetUserCommentReactionParams{
		CommentID: commentID,
		UserID:    userData.UserID,
	})

	if err == nil {
		reactionStatus = updatedUserReaction.Reaction
	}

	// Invalidate cache for this content's comments
	pattern := redis.GenerateKey("comments", updatedComment.ContentID, "*")
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), pattern)
	if err != nil {
		log.Printf("Error invalidating comments cache: %v", err)
		// Continue despite cache invalidation error
	} else {
		log.Printf("Invalidated cache for content ID: %s", updatedComment.ContentID)
	}

	// Invalidate cache for this content's reactions
	reactionsKey := redis.GenerateKey("user_reactions", updatedComment.ContentID, userData.UserID)
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), reactionsKey)
	if err != nil {
		log.Printf("Error invalidating reactions cache: %v", err)
		// Continue despite cache invalidation error
	} else {
		log.Printf("Invalidated cache for content ID: %s", updatedComment.ContentID)
	}

	// Invalidate cache for fetched replies
	if updatedComment.ParentCommentID.Valid {
		repliesKey := redis.GenerateKey("comment_replies", updatedComment.ParentCommentID, "*")
		err = server.cacheService.DeleteByPattern(ctx.Request().Context(), repliesKey)
		if err != nil {
			log.Printf("Error invalidating replies cache: %v", err)
			// Continue despite cache invalidation error
		}
	}

	// Render just the comment actions part
	return Render(ctx, http.StatusOK, components.CommentActions(updatedComment, userData, reactionStatus))
}

// Handler for downvoting a comment
func (server *Server) handleDownvoteComment(ctx echo.Context) error {
	var userData db.GetUserByIDRow
	commentIDStr := ctx.Param("id")

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		log.Println("User is not logged in.")
	}

	payload, err := server.tokenMaker.VerifyToken(cookie.Value)
	if err != nil {
		log.Println("Invalid token:", err)
		return err
	}

	userID, err := utils.ParseUUID(payload.UserID, "userID")
	if err != nil {
		log.Println("Error parsing user_id in handleDownvoteComment:", err)
		return err
	}

	user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error getting user in handleDownvoteComment:", err)
		return err
	}
	userData = user

	commentID, err := utils.ParseUUID(commentIDStr, "comment ID")
	if err != nil {
		log.Println("Invalid comment ID format in handleDownvoteComment:", err)
		return err
	}

	// Check if the user already has a reaction
	userReaction, err := server.store.GetUserCommentReaction(ctx.Request().Context(), db.GetUserCommentReactionParams{
		CommentID: commentID,
		UserID:    userData.UserID,
	})

	if err == nil {
		// User has an existing reaction
		if userReaction.Reaction == "dislike" {
			// If already disliked, remove the reaction
			_, err = server.store.DeleteCommentReaction(ctx.Request().Context(), db.DeleteCommentReactionParams{
				CommentID: commentID,
				UserID:    userData.UserID,
			})
			if err != nil {
				log.Println("Error deleting comment reaction from dislike to remove dislike:", err)
				return err
			}
		} else if userReaction.Reaction == "like" {
			// If liked, change to dislike
			_, err = server.store.InsertOrUpdateCommentReaction(ctx.Request().Context(), db.InsertOrUpdateCommentReactionParams{
				CommentID: commentID,
				UserID:    userData.UserID,
				Reaction:  "dislike",
			})
			if err != nil {
				log.Println("Error changing reaction from like to dislike:", err)
				return err
			}
		}
	} else {
		// No reaction yet, add a dislike
		_, err = server.store.InsertOrUpdateCommentReaction(ctx.Request().Context(), db.InsertOrUpdateCommentReactionParams{
			CommentID: commentID,
			UserID:    userData.UserID,
			Reaction:  "dislike",
		})
		if err != nil {
			log.Println("Error adding new dislike reaction:", err)
			return err
		}
	}

	// Update the comment's score
	updatedComment, err := server.store.UpdateCommentScore(ctx.Request().Context(), commentID)
	if err != nil {
		log.Println("Error updating comment score:", err)
		return err
	}

	// Get the updated user reaction for the response
	reactionStatus := ""
	updatedUserReaction, err := server.store.GetUserCommentReaction(ctx.Request().Context(), db.GetUserCommentReactionParams{
		CommentID: commentID,
		UserID:    userData.UserID,
	})

	if err == nil {
		reactionStatus = updatedUserReaction.Reaction
	}

	// Invalidate cache for this content's comments
	pattern := redis.GenerateKey("comments", updatedComment.ContentID, "*")
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), pattern)
	if err != nil {
		log.Printf("Error invalidating comments cache: %v", err)
		// Continue despite cache invalidation error
	} else {
		log.Printf("Invalidated cache for content ID: %s", updatedComment.ContentID)
	}

	// Invalidate cache for this content's reactions
	reactionsKey := redis.GenerateKey("user_reactions", updatedComment.ContentID, userData.UserID)
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), reactionsKey)
	if err != nil {
		log.Printf("Error invalidating reactions cache: %v", err)
		// Continue despite cache invalidation error
	} else {
		log.Printf("Invalidated cache for content ID: %s", updatedComment.ContentID)
	}

	// Invalidate cache for fetched replies
	if updatedComment.ParentCommentID.Valid {
		repliesKey := redis.GenerateKey("comment_replies", updatedComment.ParentCommentID, "*")
		err = server.cacheService.DeleteByPattern(ctx.Request().Context(), repliesKey)
		if err != nil {
			log.Printf("Error invalidating replies cache: %v", err)
			// Continue despite cache invalidation error
		}
	}

	// Render just the comment actions part
	return Render(ctx, http.StatusOK, components.CommentActions(updatedComment, userData, reactionStatus))
}

type CreateReplyReq struct {
	ReplyText string `form:"reply_text" validate:"required,min=3,max=10000"`
}

func (server *Server) createReply(ctx echo.Context) error {
	var req CreateReplyReq
	var userData db.GetUserByIDRow

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createReply:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in createReply:", err)
		return err
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in createReply:", err)
		return err
	}

	parentCommentIDStr := ctx.Param("id")
	parentCommentID, err := utils.ParseUUID(parentCommentIDStr, "parent comment ID")
	if err != nil {
		log.Println("Invalid parent comment ID format in createReply:", err)
		return err
	}

	parentComment, err := server.store.GetCommentByID(ctx.Request().Context(), parentCommentID)
	if err != nil {
		log.Println("Error getting parent comment in createReply:", err)
		return err
	}

	arg := db.CreateReplyParams{
		ParentCommentID: parentCommentID,
		UserID:          userData.UserID,
		ContentID:       parentComment.ContentID,
		CommentText:     req.ReplyText,
	}

	comment, err := server.store.CreateReply(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error creating reply:", err)
		return err
	}

	// After creating the reply, invalidate relevant caches
	replyCountCacheKey := redis.GenerateKey("reply_count", parentCommentID)
	checkedCacheKey := redis.GenerateKey("checked_admin_replies", parentCommentID)
	adminPfpCacheKey := redis.GenerateKey("admin_pfp", parentCommentID)

	// Invalidate caches related to the parent comment
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), replyCountCacheKey)
	if err != nil {
		log.Printf("Error invalidating caches for ParentCommentID: %s: %v", parentCommentID, err)
	}
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), checkedCacheKey)
	if err != nil {
		log.Printf("Error invalidating caches for ParentCommentID: %s: %v", parentCommentID, err)
	}
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), adminPfpCacheKey)
	if err != nil {
		log.Printf("Error invalidating caches for ParentCommentID: %s: %v", parentCommentID, err)
	}

	commentCountKey := redis.GenerateKey("comment_count", comment.ContentID)
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), commentCountKey)
	if err != nil {
		log.Printf("Error invalidating comment count cache: %v", err)
		// Continue despite cache invalidation error
	}

	commentRepliesKey := redis.GenerateKey("comment_replies", parentCommentID, "*")
	err = server.cacheService.DeleteByPattern(ctx.Request().Context(), commentRepliesKey)
	if err != nil {
		log.Printf("Error invalidating comment replies cache: %v", err)
		// Continue despite cache invalidation error
	}

	// Log the invalidation
	log.Printf("Invalidated caches for reply count, admin check, and admin pfp for ParentCommentID: %s", parentCommentID)

	convertedComment := db.ListContentCommentsRow{
		CommentID:       comment.CommentID,
		ContentID:       comment.ContentID,
		UserID:          comment.UserID,
		CommentText:     comment.CommentText,
		Score:           comment.Score,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		IsDeleted:       comment.IsDeleted,
		ParentCommentID: comment.ParentCommentID,
		Username:        userData.Username,
		Pfp:             userData.Pfp,
		Role:            userData.Role,
	}

	return Render(ctx, http.StatusOK, components.CommentReplyItem(convertedComment, userData, ""))
}

func (server *Server) listRepliesInfo(ctx echo.Context) error {
	commentIDStr := ctx.Param("id")

	commentID, err := utils.ParseUUID(commentIDStr, "comment ID")
	if err != nil {
		log.Println("Invalid comment ID format in listRepliesInfo:", err)
		return err
	}

	replyCount, adminPfp, err := server.getReplyCountAndAdminPfp(ctx.Request().Context(), commentID)
	if err != nil {
		log.Println("Error getting reply count and admin pfp:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.CommentReplyInfo(int(replyCount), adminPfp, commentIDStr))
}

func (server *Server) listCommentReplies(ctx echo.Context) error {
	var req ListAdsReq

	parentCommentIDStr := ctx.Param("id")
	parentCommentID, err := utils.ParseUUID(parentCommentIDStr, "parent comment ID")
	if err != nil {
		log.Println("Invalid parent comment ID format in listCommentReplies:", err)
		return err
	}

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listCommentReplies:", err)
		return err
	}

	nextLimit := req.Limit + 10

	replies, err := server.listCommentRepliesWithCache(ctx.Request().Context(), parentCommentID, nextLimit)
	if err != nil {
		log.Println("Error listing comment replies:", err)
		return err
	}

	var convertedReplies []db.ListContentCommentsRow
	for _, reply := range replies {
		convertedReplies = append(convertedReplies, db.ListContentCommentsRow{
			CommentID:       reply.CommentID,
			ContentID:       reply.ContentID,
			UserID:          reply.UserID,
			CommentText:     reply.CommentText,
			Score:           reply.Score,
			CreatedAt:       reply.CreatedAt,
			UpdatedAt:       reply.UpdatedAt,
			IsDeleted:       reply.IsDeleted,
			ParentCommentID: reply.ParentCommentID,
			Username:        reply.Username,
			Pfp:             reply.Pfp,
			Role:            reply.Role,
		})
	}

	userData, err := server.getUserFromCacheOrDb(ctx, "refresh_token")
	if err != nil {
		log.Println("Error getting user in listCommentReplies:", err)
		return err
	}

	userReactions, err := server.getUserReactionsForContentWithCache(ctx.Request().Context(), replies[0].ContentID, userData.UserID)
	if err != nil {
		log.Println("Error fetching user reactions:", err)
	}

	url := fmt.Sprintf("/api/comments/%s/more-replies?limit=", parentCommentIDStr)

	return Render(ctx, http.StatusOK, components.CommentReplyList(convertedReplies, userData, userReactions, int(nextLimit), url))
}
