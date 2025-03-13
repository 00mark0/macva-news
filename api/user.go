package api

import (
	"log"
	"net/http"
	"time"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	"github.com/00mark0/macva-news/token"
	"github.com/00mark0/macva-news/utils"

	"os"

	"github.com/labstack/echo/v4"
)

type userResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Pfp      string `json:"pfp"`
	Role     string `json:"role"`
}

func newUserResponse(payload *token.Payload) userResponse {
	return userResponse{
		UserID:   payload.UserID,
		Username: payload.Username,
		Email:    payload.Email,
		Pfp:      payload.Pfp,
		Role:     payload.Role,
	}
}

type loginUserReq struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type loginUserRes struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) login(ctx echo.Context) error {
	var req loginUserReq
	var loginErr components.LoginErr
	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in login:", err)
		return err
	}

	if req.Email == "" {
		loginErr = "Email obavezan"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	if req.Password == "" {
		loginErr = "Lozinka obavezna"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	user, err := server.store.GetUserByEmail(ctx.Request().Context(), req.Email)
	if err != nil {
		loginErr = "Nevažeći podaci za prijavu"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	if user.Banned.Bool {
		loginErr = "Nevažecí podaci za prijavu"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		loginErr = "Nevažecí podaci za prijavu"

		return Render(ctx, http.StatusOK, components.LoginForm(loginErr))
	}

	durationStr := os.Getenv("ACCESS_TOKEN_DURATION")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Println("Error parsing duration in login:", err)
		return err
	}

	accessToken, payload, err := server.tokenMaker.CreateToken(
		user.UserID.String(),
		user.Username,
		user.Email,
		user.Pfp,
		user.Role,
		duration,
	)
	if err != nil {
		log.Println("Error creating token in login:", err)
		return err
	}

	res := loginUserRes{
		AccessToken: accessToken,
		User:        newUserResponse(payload),
	}

	if acceptHeader := ctx.Request().Header.Get("Accept"); acceptHeader == "application/json" {
		return ctx.JSON(http.StatusOK, res)
	}

	// Set token as a secure, HTTP-only cookie
	ctx.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(duration),
		Path:     "/",
		HttpOnly: true,
	})

	ctx.Response().Header().Set("HX-Redirect", "/")
	return ctx.NoContent(http.StatusOK)
}

func (server *Server) logOut(ctx echo.Context) error {
	ctx.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   "",
		MaxAge:  -1,
		Path:    "/",
		Expires: time.Now(),
	})

	ctx.Response().Header().Set("HX-Redirect", "/")
	return ctx.NoContent(http.StatusOK)
}

type ListUsersReq struct {
	Limit int32 `query:"limit"`
}

type SearchUserReq struct {
	SearchTerm string `query:"search_term" validate:"required"`
	Limit      int32  `query:"limit"`
}

func (server *Server) listActiveUsers(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listActiveUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetActiveUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listActiveUsers:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/active?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listActiveUsersOldest(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listActiveUsersOldest:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetActiveUsersOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listActiveUsersOldest:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/active/oldest?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listActiveUsersTitle(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listActiveUsersTitle:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetActiveUsersTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listActiveUsersTitle:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/active/title?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listBannedUsers(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listBannedUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetBannedUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listBannedUsers:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/banned?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listBannedUsersOldest(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listBannedUsersOldest:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetBannedUsersOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listBannedUsersOldest:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/banned/oldest?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listBannedUsersTitle(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listBannedUsersTitle:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetBannedUsersTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listBannedUsersTitle:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/banned/title?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listDeletedUsers(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDeletedUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetDeletedUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listDeletedUsers:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/deleted?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listDeletedUsersOldest(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDeletedUsersOldest:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetDeletedUsersOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listDeletedUsersOldest:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/deleted/oldest?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) listDeletedUsersTitle(ctx echo.Context) error {
	var req ListUsersReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listDeletedUsersTitle:", err)
		return err
	}

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetDeletedUsersTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in listDeletedUsersTitle:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/deleted/title?limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) searchActiveUsers(ctx echo.Context) error {
	var req SearchUserReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in searchActiveUsers:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in searchActiveUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchActiveUsersParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	var users []components.UsersRes

	data, err := server.store.SearchActiveUsers(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching users in searchActiveUsers:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/active/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) searchBannedUsers(ctx echo.Context) error {
	var req SearchUserReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in searchBannedUsers:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in searchBannedUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchBannedUsersParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	var users []components.UsersRes

	data, err := server.store.SearchBannedUsers(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching users in searchBannedUsers:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/banned/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}

func (server *Server) searchArchivedUsers(ctx echo.Context) error {
	var req SearchUserReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in searchArchivedUsers:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in searchArchivedUsers:", err)
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchDeletedUsersParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	var users []components.UsersRes

	data, err := server.store.SearchDeletedUsers(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching users in searchArchivedUsers:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/deleted/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.Users(int(nextLimit), users, url))
}
