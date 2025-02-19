package db

import (
	"log"
	"testing"

	"context"

	//"github.com/00mark0/macva-news/utils"
	"github.com/go-loremipsum/loremipsum"
	"github.com/stretchr/testify/require"
)

var loremipsumgen = loremipsum.New()

func createRandomContent(t *testing.T) Content {
	users, err := testQueries.GetAdminUsers(context.Background())
	require.NoError(t, err)

	category, err := testQueries.ListCategories(context.Background(), ListCategoriesParams{Limit: 10, Offset: 0})

	arg := CreateContentParams{
		UserID:              users[0].UserID,
		CategoryID:          category[0].CategoryID,
		Title:               loremipsumgen.Sentence(),
		ContentDescription:  loremipsumgen.Paragraphs(10),
		CommentsEnabled:     true,
		ViewCountEnabled:    true,
		LikeCountEnabled:    true,
		DislikeCountEnabled: false,
	}

	content, err := testQueries.CreateContent(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.UserID, content.UserID)
	require.Equal(t, arg.CategoryID, content.CategoryID)
	require.Equal(t, arg.Title, content.Title)
	require.Equal(t, arg.ContentDescription, content.ContentDescription)

	return content
}

func TestCreateContent(t *testing.T) {
	var contents []Content

	for i := 0; i < 10; i++ {
		contents = append(contents, createRandomContent(t))
	}

	for _, content := range contents {
		require.NotEmpty(t, content)
		require.NotEmpty(t, content.ContentID)
		require.NotEmpty(t, content.UserID)
		require.NotEmpty(t, content.CategoryID)
		require.NotEmpty(t, content.Title)
		require.NotEmpty(t, content.ContentDescription)
	}
}

func TestUpdateContent(t *testing.T) {
	content1 := createRandomContent(t)

	arg := UpdateContentParams{
		ContentID:           content1.ContentID,
		Title:               loremipsumgen.Sentence(),
		ContentDescription:  loremipsumgen.Paragraphs(8),
		CommentsEnabled:     true,
		ViewCountEnabled:    true,
		LikeCountEnabled:    true,
		DislikeCountEnabled: false,
	}

	content2, err := testQueries.UpdateContent(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.ContentID, content2.ContentID)
	require.Equal(t, arg.Title, content2.Title)
	require.Equal(t, arg.ContentDescription, content2.ContentDescription)
	require.Equal(t, arg.CommentsEnabled, content2.CommentsEnabled)
	require.Equal(t, arg.ViewCountEnabled, content2.ViewCountEnabled)
	require.Equal(t, arg.LikeCountEnabled, content2.LikeCountEnabled)
	require.Equal(t, arg.DislikeCountEnabled, content2.DislikeCountEnabled)
}

func TestPublishContent(t *testing.T) {
	content1 := createRandomContent(t)

	content2, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)

	require.Equal(t, content2.Status, "published")
	require.Equal(t, content2.PublishedAt.Valid, true)
	require.NotNil(t, content2.PublishedAt.Time)
	require.NotEqual(t, content2.Status, content1.Status)
}

func TestSoftDeleteContent(t *testing.T) {
	content1 := createRandomContent(t)

	content2, err := testQueries.SoftDeleteContent(context.Background(), content1.ContentID)
	require.NoError(t, err)

	require.Equal(t, content1.IsDeleted.Bool, false)
	require.Equal(t, content2.IsDeleted.Bool, true)
}

func TestHardDeleteContent(t *testing.T) {
	content1 := createRandomContent(t)

	content2, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1.ContentID, content2.ContentID)

	count1, err := testQueries.GetPublishedContentCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count1)

	content3, err := testQueries.HardDeleteContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1.ContentID, content3.ContentID)

	count2, err := testQueries.GetPublishedContentCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, count2, count1-1)
}

func TestGetContentDetails(t *testing.T) {
	content1 := createRandomContent(t)

	content2, err := testQueries.GetContentDetails(context.Background(), content1.ContentID)
	require.NoError(t, err)

	require.Equal(t, content1.ContentID, content2.ContentID)
	require.Equal(t, content1.UserID, content2.UserID)
	require.Equal(t, content1.CategoryID, content2.CategoryID)
	require.Equal(t, content1.Title, content2.Title)
	require.Equal(t, content1.ContentDescription, content2.ContentDescription)
	require.Equal(t, content1.CommentsEnabled, content2.CommentsEnabled)
	require.Equal(t, content1.ViewCountEnabled, content2.ViewCountEnabled)
	require.Equal(t, content1.LikeCountEnabled, content2.LikeCountEnabled)
	require.Equal(t, content1.DislikeCountEnabled, content2.DislikeCountEnabled)
	require.Equal(t, content1.Status, content2.Status)
	require.Equal(t, content1.ViewCount, content2.ViewCount)
	require.Equal(t, content1.LikeCount, content2.LikeCount)
	require.Equal(t, content1.DislikeCount, content2.DislikeCount)
	require.Equal(t, content1.CommentCount, content2.CommentCount)
	require.Equal(t, content1.CreatedAt, content2.CreatedAt)
	require.Equal(t, content1.UpdatedAt, content2.UpdatedAt)
	require.Equal(t, content1.PublishedAt, content2.PublishedAt)
	require.Equal(t, content1.IsDeleted, content2.IsDeleted)
	log.Println(content2.AuthorUsername)
	log.Println(content2.CategoryName)
	log.Println(content2.ReactionCount)
	log.Println(content2.CommentCountSync)
}

func TestGetPublishedContentCount(t *testing.T) {
	count1, err := testQueries.GetPublishedContentCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count1)

	content1 := createRandomContent(t)

	content2, err := testQueries.PublishContent(context.Background(), content1.ContentID)
	require.NoError(t, err)
	require.Equal(t, content1.ContentID, content2.ContentID)

	count2, err := testQueries.GetPublishedContentCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, count2, count1+1)
}

func TestListPublishedContent(t *testing.T) {
	content, err := testQueries.ListPublishedContent(context.Background(), ListPublishedContentParams{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, content)
	require.LessOrEqual(t, len(content), 10)

	for _, content := range content {
		require.NotEmpty(t, content)
		require.NotEmpty(t, content.ContentID)
		require.NotEmpty(t, content.UserID)
		require.NotEmpty(t, content.CategoryID)
		require.NotEmpty(t, content.Title)
		require.NotEmpty(t, content.ContentDescription)
	}
}
