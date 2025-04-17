package api

import (
	"context"
	"log"
	"time"

	"github.com/00mark0/macva-news/db/redis"
	"github.com/00mark0/macva-news/db/services"
	"github.com/jackc/pgx/v5/pgtype"
	redisClient "github.com/redis/go-redis/v9"
)

func (server *Server) getCommentsWithCache(ctx context.Context, contentID pgtype.UUID, limit int32) ([]db.ListContentCommentsRow, error) {
	// Generate the cache key
	cacheKey := redis.GenerateKey("comments", contentID, limit)

	// Try to get comments from cache
	var comments []db.ListContentCommentsRow
	cacheHit, err := server.cacheService.Get(ctx, cacheKey, &comments)
	if err != nil {
		log.Printf("Error fetching comments from cache: %v", err)
	}

	if cacheHit {
		// Cache hit, return cached comments
		log.Printf("Cache hit for comments: %s", cacheKey)
		return comments, nil
	}

	// Cache miss, fetch from database
	log.Printf("Cache miss for comments: %s", cacheKey)
	arg := db.ListContentCommentsParams{
		ContentID: contentID,
		Limit:     limit,
	}
	comments, err = server.store.ListContentComments(ctx, arg)
	if err != nil {
		return nil, err
	}

	// Store in cache for future use
	err = server.cacheService.Set(ctx, cacheKey, comments, 10*time.Minute)
	if err != nil {
		log.Printf("Error caching comments: %v", err)
	}

	return comments, nil
}

func (server *Server) getReplyCountAndAdminPfp(ctx context.Context, parentCommentID pgtype.UUID) (int64, string, error) {
	// Generate cache keys
	checkedCacheKey := redis.GenerateKey("checked_admin_replies", parentCommentID)
	adminPfpCacheKey := redis.GenerateKey("admin_pfp", parentCommentID)
	countCacheKey := redis.GenerateKey("reply_count", parentCommentID)

	// Check if we have the reply count cached
	var replyCount int64
	countCacheHit, err := server.cacheService.Get(ctx, countCacheKey, &replyCount)
	if err != nil && err != redisClient.Nil {
		log.Printf("Error checking cached reply count for ParentCommentID: %s: %v", parentCommentID, err)
	}

	// If count not in cache, fetch from database
	if !countCacheHit {
		log.Printf("Cache miss for reply count for ParentCommentID: %s, fetching from DB", parentCommentID)
		replyCount, err = server.store.GetReplyCount(ctx, parentCommentID)
		if err != nil {
			log.Printf("Error fetching reply count for ParentCommentID: %s: %v", parentCommentID, err)
			return 0, "", err
		}

		// Cache the reply count
		err = server.cacheService.Set(ctx, countCacheKey, replyCount, 10*time.Minute)
		if err != nil {
			log.Printf("Error caching reply count for ParentCommentID: %s: %v", parentCommentID, err)
		}
		log.Printf("Cached reply count %d for ParentCommentID: %s", replyCount, parentCommentID)
	} else {
		log.Printf("Cache hit for reply count for ParentCommentID: %s: %d", parentCommentID, replyCount)
	}

	// Check if we've already checked this parent comment for admin replies
	var hasAdminReply bool
	adminCheckCacheHit, err := server.cacheService.Get(ctx, checkedCacheKey, &hasAdminReply)
	if err != nil && err != redisClient.Nil {
		log.Printf("Error checking admin reply status for ParentCommentID: %s: %v", parentCommentID, err)
	}

	var adminPfp string

	if adminCheckCacheHit {
		// We already know if this comment has admin replies
		log.Printf("Cache hit for checked_admin_replies for ParentCommentID: %s", parentCommentID)
		if hasAdminReply {
			// Only try to get the cached pfp if we know there's an admin reply
			pfpCacheHit, err := server.cacheService.Get(ctx, adminPfpCacheKey, &adminPfp)
			if err != nil && err != redisClient.Nil {
				log.Printf("Error fetching cached admin pfp for ParentCommentID: %s: %v", parentCommentID, err)
			}

			if pfpCacheHit && adminPfp != "" {
				// Cache hit: use cached admin pfp
				log.Printf("Cache hit for admin pfp for ParentCommentID: %s", parentCommentID)
			} else {
				// Cache miss: admin pfp needs to be rescanned
				log.Printf("Admin pfp cache miss for ParentCommentID: %s, rescanning replies", parentCommentID)
				adminPfp = scanForAdminPfp(ctx, server, parentCommentID)
			}
		} else {
			// No admin reply found, skip rescan
			log.Printf("No admin reply in cache for ParentCommentID: %s, skipping rescan", parentCommentID)
		}
	} else {
		// If no cache hit, need to scan replies for admin replies
		log.Printf("Cache miss for checked_admin_replies for ParentCommentID: %s, scanning for admin", parentCommentID)
		adminPfp = scanForAdminPfp(ctx, server, parentCommentID)

		// Cache the checked status
		hasAdminReply = adminPfp != ""
		err = server.cacheService.Set(ctx, checkedCacheKey, hasAdminReply, 10*time.Minute)
		if err != nil {
			log.Printf("Error caching admin reply status for ParentCommentID: %s: %v", parentCommentID, err)
		}

		if hasAdminReply {
			log.Printf("Admin reply found, set checked status for ParentCommentID: %s", parentCommentID)
		} else {
			log.Printf("No admin reply found, set checked status for ParentCommentID: %s", parentCommentID)
		}
	}

	// Final log statement before returning
	log.Printf("Returning reply count %d and admin pfp for ParentCommentID: %s", replyCount, parentCommentID)

	return replyCount, adminPfp, nil
}

// Helper function to scan for admin pfp to reduce code duplication
func scanForAdminPfp(ctx context.Context, server *Server, parentCommentID pgtype.UUID) string {
	// Fetch all replies (using a reasonable high limit)
	allReplies, err := server.store.ListCommentReplies(ctx, db.ListCommentRepliesParams{
		ParentCommentID: parentCommentID,
		Limit:           10000, // Reasonable high limit
	})
	if err != nil {
		log.Printf("Error fetching all replies for ParentCommentID: %s: %v", parentCommentID, err)
		return ""
	}

	// Look for admin replies
	adminPfpCacheKey := redis.GenerateKey("admin_pfp", parentCommentID)
	for _, reply := range allReplies {
		if reply.Role == "admin" {
			adminPfp := reply.Pfp

			// Cache the admin pfp
			err := server.cacheService.Set(ctx, adminPfpCacheKey, adminPfp, 10*time.Minute)
			if err != nil {
				log.Printf("Error caching admin pfp for ParentCommentID: %s: %v", parentCommentID, err)
			} else {
				log.Printf("Found admin reply, cached admin pfp for ParentCommentID: %s", parentCommentID)
			}
			return adminPfp
		}
	}

	return ""
}
