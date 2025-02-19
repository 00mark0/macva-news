package db

import (
	"log"
	"testing"

	"github.com/00mark0/macva-news/utils"

	"context"

	"github.com/stretchr/testify/require"
)

func createCategoryInteractive(name string) Category {
	category, err := testQueries.CreateCategory(context.Background(), name)
	if err != nil {
		log.Fatal(err)
	}
	return category
}

func createCategories(t *testing.T) []Category {
	categoryStrings := utils.RandomCategoryList(10)
	var categories []Category
	for i := 0; i < 10; i++ {
		category := createCategoryInteractive(categoryStrings[i])
		require.NotEmpty(t, category)
		require.Equal(t, categoryStrings[i], category.CategoryName)

		categories = append(categories, category)
	}
	return categories
}

func TestCreateCategory(t *testing.T) {
	categories := createCategories(t)
	require.NotEmpty(t, categories)
	require.Len(t, categories, 10)

	for _, category := range categories {
		require.NotEmpty(t, category)
		require.NotEmpty(t, category.CategoryID)
		require.NotEmpty(t, category.CategoryName)
	}
}

func TestGetCategory(t *testing.T) {
	category1, err := testQueries.CreateCategory(context.Background(), utils.RandomCategory())
	require.NoError(t, err)

	category2, err := testQueries.GetCategory(context.Background(), category1.CategoryID)
	require.NoError(t, err)
	require.Equal(t, category1.CategoryID, category2.CategoryID)
	require.Equal(t, category1.CategoryName, category2.CategoryName)
}

func TestGetCategoryByName(t *testing.T) {
	category1, err := testQueries.CreateCategory(context.Background(), utils.RandomString(10))
	require.NoError(t, err)

	category2, err := testQueries.GetCategoryByName(context.Background(), category1.CategoryName)
	require.NoError(t, err)
	require.Equal(t, category1.CategoryID, category2.CategoryID)
	require.Equal(t, category1.CategoryName, category2.CategoryName)
}

func TestGetCategoryCount(t *testing.T) {
	count, err := testQueries.GetCategoryCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)

	_ = createCategoryInteractive(utils.RandomString(10))

	count2, err := testQueries.GetCategoryCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, count2, count+1)
}

func TestListCategories(t *testing.T) {
	arg := ListCategoriesParams{
		Limit:  5,
		Offset: 0,
	}

	categories, err := testQueries.ListCategories(context.Background(), arg)
	require.NoError(t, err)
	require.LessOrEqual(t, len(categories), 5)

	for _, category := range categories {
		require.NotEmpty(t, category)
		require.NotEmpty(t, category.CategoryID)
		require.NotEmpty(t, category.CategoryName)
	}
}

func TestUpdateCategory(t *testing.T) {
	category1, err := testQueries.CreateCategory(context.Background(), utils.RandomString(10))
	require.NoError(t, err)

	arg := UpdateCategoryParams{
		CategoryID:   category1.CategoryID,
		CategoryName: utils.RandomString(10),
	}

	category2, err := testQueries.UpdateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, category1.CategoryID, category2.CategoryID)
	require.Equal(t, arg.CategoryName, category2.CategoryName)
}

func TestDeleteCategory(t *testing.T) {
	category1, err := testQueries.CreateCategory(context.Background(), utils.RandomString(10))

	count1, err := testQueries.GetCategoryCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count1)

	deleted, err := testQueries.DeleteCategory(context.Background(), category1.CategoryID)
	require.NoError(t, err)
	require.Equal(t, category1.CategoryID, deleted.CategoryID)

	count2, err := testQueries.GetCategoryCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, count2, count1-1)
}
