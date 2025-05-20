package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sugurta/api/handlers"
	status_http "sugurta/api/http_status"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/config"
	em "sugurta/internal/pkg/email"
	"sugurta/internal/pkg/helper"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sugurta/internal/usecase/user"
)

// Error constants
const (
	ErrEmailExists      = "email already exists"
	ErrWrongEmailOrPass = "wrong email or password"
	ErrCodeExpired      = "verification code has expired"
	ErrIncorrectCode    = "incorrect verification code"
	BadRequest          = "bad request"
	InternalServerError = "internal server error"
)

// Cache key constants
const (
	RegisterCodeKey   = "register_code_"
	ForgotPasswordKey = "forgot_password_code_"
)

// authRoutes represents the auth controller
type authRoutes struct {
	handlers.BaseHandler
	log         *zap.Logger
	cfg         *config.Config
	enforcer    *casbin.CachedEnforcer
	userUseCase user.User
}

// NewAuthRoutes creates a new auth routes controller
func NewAuthRoutes(apiV1Group *gin.RouterGroup, option *handlers.HandlerOption) {
	r := &authRoutes{
		log:         option.Logger,
		cfg:         option.Config,
		enforcer:    option.Enforcer,
		userUseCase: option.User,
	}

	r.Cache = option.Cache
	r.Config = option.Config

	authGroup := apiV1Group.Group("/auth")
	{

		authGroup.POST("/register", r.register)
		authGroup.POST("/verify", r.verify)
		authGroup.POST("/refresh-token/:token", r.refreshAccessToken)
		authGroup.POST("/login", r.login)
		authGroup.POST("/forgot-password", r.forgotPassword)
		authGroup.POST("/verify-forgot-password", r.verifyForgotPassword)
		authGroup.DELETE("/delete/:id", r.deleteAccount)
		authGroup.POST("/update-password", r.updatePassword)
		authGroup.GET("/clients/list", r.ListClients)
		authGroup.GET("/clients/:id", r.GetClientByID)
		authGroup.GET("/users/:id", r.GetUserByID)
		authGroup.GET("/users/list", r.ListUsers)
		authGroup.PUT("/clients/block", r.BlockClient)
		authGroup.PUT("/clients/pause", r.PauseClientChat)

	}
}

// @Router /auth/register [post]
// @Summary Register
// @Description Register
// @Tags AUTH
// @Accept json
// @Produce json
// @Param admin body entity.RegisterRequest true "RegisterRequest"
// @Success 200 {object} status_http.Response
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (a *authRoutes) register(c *gin.Context) {
	var (
		req entity.RegisterRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	fmt.Println("1")

	// Password validation
	isValid := helper.ValidatePassword(req.Password)
	if !isValid {
		a.handleResponse(c, status_http.BadRequest, "invalid password: must be at least 6 characters with at least one uppercase and one lowercase letter")
		return
	}
	fmt.Println("2")

	// Phone number validation
	if req.PhoneNumber != "" {
		isValid = helper.ValidatePhoneNumber(req.PhoneNumber)
		if !isValid {
			a.handleResponse(c, status_http.BadRequest, "invalid phone number: must be in format +998XXXXXXXXX")
			return
		}
	}
	fmt.Println("3", req.Email)

	// Check if email already exists
	_, err = a.userUseCase.Get(context.Background(), map[string]string{
		"email": req.Email,
	})

	fmt.Println("3.1")

	if err == nil {
		a.handleResponse(c, status_http.BadRequest, ErrEmailExists)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	fmt.Println("4")

	// Hash password
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	fmt.Println("5")

	req.Password = hashedPassword

	fmt.Println("6")

	err = a.Cache.Set(c, "email_"+req.Email, req, 5*time.Minute)
	if err != nil {
		a.handleResponse(c, status_http.InternalServerError, err)
		return
	}

	fmt.Println("7")

	// Send verification code
	go func() {
		err := a.sendVerificationCode(RegisterCodeKey, req.Email)
		if err != nil {
			a.log.Error(fmt.Sprintf("failed to send verification code: %v", err))
		}
	}()

	fmt.Println("8")

	a.handleResponse(c, status_http.OK, "Verification code send to email")
}

// sendVerificationCode sends a verification code to the email
func (a *authRoutes) sendVerificationCode(key, email string) error {
	// Generate random code
	code, err := helper.GenerateRandomCode(6)
	if err != nil {
		return err
	}

	// Store code in cache
	err = a.Cache.Set(context.Background(), key+email, code, 5*time.Minute)
	if err != nil {
		return err
	}

	// Send email
	emailCfg := &config.EmailConfig{
		From:     a.cfg.Email.From,
		Password: a.cfg.Email.Password,
		Host:     a.cfg.Email.Host,
		Port:     a.cfg.Email.Port,
	}

	fmt.Println("asdfasdfasdfasdf: ", a.cfg.Email.From, " ", a.cfg.Email.Password)

	err = em.SendEmail(emailCfg, &em.SendEmailRequest{
		To:      []string{email},
		Subject: "Verification email",
		Body: map[string]string{
			"code": code,
		},
		Type: em.VerificationEmail,
	})
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	return nil
}

// @Router /auth/verify [post]
// @Summary Verify user
// @Description Verify user
// @Tags AUTH
// @Accept json
// @Produce json
// @Param data body entity.VerifyEmail true "Data"
// @Success 201 {object} status_http.Response
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (a *authRoutes) verify(c *gin.Context) {
	var (
		req entity.VerifyEmail
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	// Get stored user data
	userData, err := a.Cache.Get(c, "email_"+req.Email)
	if err != nil {
		a.handleResponse(c, status_http.Forbidden, err.Error())
		return
	}
	fmt.Println("data:", string(userData))
	var user entity.RegisterRequest
	err = json.Unmarshal(userData, &user)
	if err != nil {
		fmt.Println("UNMARSHAL ERROR:", err)
		a.handleResponse(c, status_http.Forbidden, err.Error())
		return
	}

	// Get verification code
	code, err := a.Cache.Get(c, RegisterCodeKey+user.Email)
	if err != nil {
		a.handleResponse(c, status_http.Forbidden, ErrCodeExpired)
		return
	}

	// Verify code
	codeStr := strings.Trim(string(code), "\"")
	if req.Code != codeStr {
		a.handleResponse(c, status_http.Forbidden, ErrIncorrectCode)
		return
	}

	// Generate tokens
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	roleData := entity.Role{Name: "user"}
	// Create user in database
	userResp, err := a.userUseCase.Create(ctx, &entity.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		RoleData:  roleData,
	})
	if err != nil {
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := helper.GenerateJWT(userResp.ID, userResp.BusinessID, "user", a.cfg.JWT.Secret, 12) // 12 hours
	if err != nil {
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	userResp.AccessToken = accessToken
	userResp.RefreshToken = refreshToken

	// Success response
	a.handleResponse(c, status_http.OK, userResp)
}

// @Router /auth/refresh-token/:token [post]
// @Summary Creates new valid access-token
// @Description Refresh access-token
// @Tags AUTH
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {object} entity.AccessToken
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (a *authRoutes) refreshAccessToken(c *gin.Context) {
	token := c.Param("token")

	// Parse and validate token
	claims, err := helper.ParseToken(token, a.cfg.JWT.Secret)
	if err != nil {
		a.handleResponse(c, status_http.Forbidden, "access token expired")
		return
	}

	// Generate new access token
	accessToken, _, err := helper.GenerateJWT(claims.Sub, claims.BusinessId, claims.Role, a.cfg.JWT.Secret, 12) // 12 hours
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.AccessToken{
		AccessToken: accessToken,
	})
}

// @Router /auth/login [post]
// @Summary Login
// @Description Login by email and password
// @Tags AUTH
// @Accept json
// @Produce json
// @Param data body entity.LoginRequest true "Data"
// @Success 200 {object} status_http.Response{data=entity.User} "User data"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (a *authRoutes) login(c *gin.Context) {
	var (
		req entity.LoginRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get user by email
	user, err := a.userUseCase.Get(ctx, map[string]string{
		"email": req.Email,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.handleResponse(c, status_http.BadRequest, ErrWrongEmailOrPass)
			return
		}
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	// Check password
	err = helper.CheckPassword(req.Password, user.Password)
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, ErrWrongEmailOrPass)
		return
	}

	// Generate tokens

	accessToken, refreshToken, err := helper.GenerateJWT(user.ID, user.BusinessID, user.RoleData.Name, a.cfg.JWT.Secret, 12) // 12 hours
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	// Success response
	a.handleResponse(c, status_http.OK, user)
}

// @Router /auth/forgot-password [post]
// @Summary Forgot password
// @Description Forgot password
// @Tags AUTH
// @Accept json
// @Produce json
// @Param data body entity.ForgotPasswordRequest true "Data"
// @Success 200 {object} status_http.Response{data=string} "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (a *authRoutes) forgotPassword(c *gin.Context) {
	var (
		req entity.ForgotPasswordRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	// Check if user exists
	_, err = a.userUseCase.Get(c, map[string]string{
		"email": req.Email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.handleResponse(c, status_http.BadRequest, "user not found")
			return
		}
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	// Send verification code
	go func() {
		err := a.sendVerificationCode(ForgotPasswordKey, req.Email)
		if err != nil {
			a.log.Error(fmt.Sprintf("failed to send verification code: %v", err))
		}
	}()

	// Success response
	a.handleResponse(c, status_http.OK, "verification code sent")
}

// @Router /auth/verify-forgot-password [post]
// @Summary Verify forgot password
// @Description Verify forgot password
// @Tags AUTH
// @Accept json
// @Produce json
// @Param data body entity.VerifyRequest true "Data"
// @Success 201 {object} status_http.Response{data=string} "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (a *authRoutes) verifyForgotPassword(c *gin.Context) {
	var (
		req entity.VerifyRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	// Get verification code
	code, err := a.Cache.Get(c, ForgotPasswordKey+req.Email)
	if err != nil {
		a.handleResponse(c, status_http.Forbidden, ErrCodeExpired)
		return
	}

	// Verify code
	codeStr := strings.Trim(string(code), "\"")
	if req.Code != codeStr {
		a.handleResponse(c, status_http.Forbidden, ErrIncorrectCode)
		return
	}
	fmt.Println(1111111)
	// Hash new password
	hashedPassword, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	fmt.Println(222222222)
	// Update password
	err = a.userUseCase.UpdatePassword(c, &entity.UpdatePasswordRequest{
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		fmt.Println(3333333333)
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}
	fmt.Println(44444444444)
	// Success response
	a.handleResponse(c, status_http.OK, "Password successfully updated!")
}

// @Router /auth/update-password [post]
// @Summary Update password
// @Description Update password
// @Security BearerAuth
// @Tags AUTH
// @Accept json
// @Produce json
// @Param data body entity.UpdatePasswordRequest true "Data"
// @Success 201 {object} status_http.Response{data=string} "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (a *authRoutes) updatePassword(c *gin.Context) {
	var (
		req entity.UpdatePasswordRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	// Validate password
	isValid := helper.ValidatePassword(req.Password)
	if !isValid {
		a.handleResponse(c, status_http.BadRequest, "invalid password: must be at least 6 characters with at least one uppercase and one lowercase letter")
		return
	}

	// Hash password
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		a.handleResponse(c, status_http.BadRequest, err.Error())
		return
	}

	// Update password
	err = a.userUseCase.UpdatePassword(c, &entity.UpdatePasswordRequest{
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	// Success response
	a.handleResponse(c, status_http.OK, "Password has been updated")
}

// @Router /auth/delete/{id} [delete]
// @Summary Delete account
// @Description Delete account with id
// @Security BearerAuth
// @Tags AUTH
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} status_http.Response{data=string} "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (a *authRoutes) deleteAccount(c *gin.Context) {
	id := c.Param("id")

	// Delete account
	err := a.userUseCase.Delete(c, id)
	if err != nil {
		a.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	// Success response
	a.handleResponse(c, status_http.OK, "Successfully deleted")
}

// @Router /auth/clients/list [get]
// @Summary Get list of clients
// @Description Get a list of clients with optional filtering by name, phone, from, goal, order_status and pagination
// @Security BearerAuth
// @Tags AUTH
// @Accept json
// @Produce json
// @Param name query string false "Filter by name"
// @Param phone query string false "Filter by phone"
// @Param from query string false "Filter by source (from)"
// @Param goal query string false "Filter by goal"
// @Param client_id query int false "Filter by client id"
// @Param order_status query string false "Filter by order status"
// @Param limit query int false "Limit the number of clients" default(10)
// @Param page query int false "Page number for pagination" default(1)
// @Success 200 {array} entity.Client "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *authRoutes) ListClients(c *gin.Context) {
	filter := entity.ClientFilter{
		Name:        c.DefaultQuery("name", ""),
		Phone:       c.DefaultQuery("phone", ""),
		From:        c.DefaultQuery("from", ""),
		Goal:        c.DefaultQuery("goal", ""),
		OrderStatus: c.DefaultQuery("order_status", ""),
	}
	clientid, err := strconv.Atoi(c.DefaultQuery("client_id", "0"))
	if err != nil {
		clientid = 0
	}
	filter.ClientId=clientid
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	filter.Limit = limit
	filter.Page = page

	clients, err := r.userUseCase.ListClients(c, filter)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, clients)
}

// @Router /auth/clients/{id} [get]
// @Summary Get client by ID
// @Description Retrieve a single client by its unique identifier (GUID)
// @Security BearerAuth
// @Tags AUTH
// @Accept json
// @Produce json
// @Param id path string true "Client ID (GUID)"
// @Success 200 {object} entity.Client "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 404 {object} status_http.Response{data=string} "Client Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *authRoutes) GetClientByID(c *gin.Context) {
	id := c.Param("id")

	client, err := r.userUseCase.GetClientByID(c, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			r.handleResponse(c, status_http.NotFound, err.Error())
			return
		}
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, client)
}

// @Router /auth/clients/block [put]
// @Summary Block or unblock a client
// @Description Block or unblock a client by business ID and platform ID
// @Security BearerAuth
// @Tags AUTH
// @Accept json
// @Produce json
// @Param data body entity.BlockUser true "Block User Request"
// @Success 200 {object} status_http.Response{data=string} "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
func (r *authRoutes) BlockClient(c *gin.Context) {
	var req entity.BlockUser

	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid request payload")
		return
	}

	if req.PlatformId == "" || req.BusinessID == "" {
		r.handleResponse(c, status_http.BadRequest, "Missing required fields: platform_id or business_id")
		return
	}

	err := r.userUseCase.BlockUser(c, req)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	action := "unblocked"
	if req.Block {
		action = "blocked"
	}

	r.handleResponse(c, status_http.OK, fmt.Sprintf("Client successfully %s", action))
}

// @Router /auth/clients/pause [put]
// @Summary Pause or unpause chat for a client
// @Description Pause or unpause a client's chat by business ID, platform ID, and source type (e.g., bot or channel)
// @Security BearerAuth
// @Tags AUTH
// @Accept json
// @Produce json
// @Param data body entity.PauzeChat true "Pause Chat Request"
// @Success 200 {object} status_http.Response{data=string} "Chat pause status updated successfully"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Internal Server Error"
func (r *authRoutes) PauseClientChat(c *gin.Context) {
	var req entity.PauzeChat

	if err := c.ShouldBindJSON(&req); err != nil {
		r.handleResponse(c, status_http.BadRequest, "Invalid request payload")
		return
	}

	if req.PlatformId == "" || req.BusinessID == "" || req.Type == "" {
		r.handleResponse(c, status_http.BadRequest, "Missing required fields: platform_id, business_id or type")
		return
	}

	err := r.userUseCase.PauzChat(c, req)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	action := "unpaused"
	if req.Pauze {
		action = "paused"
	}

	r.handleResponse(c, status_http.OK, fmt.Sprintf("Client chat successfully %s", action))
}

// @Router /auth/users/list [get]
// @Summary Get list of users
// @Description Get a list of users with optional filtering by is_active, role_id, created_at, and pagination
// @Security BearerAuth
// @Tags AUTH
// @Accept json
// @Produce json
// @Param is_active query bool false "Filter by active status"
// @Param role_id query string false "Filter by role ID"
// @Param limit query int false "Limit the number of users" default(10)
// @Param page query int false "Page number for pagination" default(1)
// @Success 200 {array} entity.User "List of users"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @ID list-users
func (r *authRoutes) ListUsers(c *gin.Context) {
	var filter entity.UserFilter

	// is_active ni bool ga parse qilish
	isActiveStr := c.Query("is_active")
	if isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			r.handleResponse(c, status_http.BadRequest, "Invalid is_active value")
			return
		}
		filter.IsActive = &isActive
	}

	filter.RoleID = c.Query("role_id")

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	filter.Limit = uint64(limit)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	filter.Offset = uint64(page)

	users, err := r.userUseCase.List(c, filter)
	if err != nil {
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	r.handleResponse(c, status_http.OK, users)
}

// @Router /auth/users/{id} [get]
// @Summary Get user by ID
// @Description Retrieve a single user by their unique identifier (GUID)
// @Security BearerAuth
// @Tags AUTH
// @Accept json
// @Produce json
// @Param id path string true "User ID (GUID)"
// @Success 200 {object} entity.User "Success"
// @Failure 400 {object} status_http.Response{data=string} "Bad Request"
// @Failure 404 {object} status_http.Response{data=string} "User Not Found"
// @Failure 500 {object} status_http.Response{data=string} "Server Error"
// @ID get-user-by-id
func (r *authRoutes) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	users, err := r.userUseCase.GetByIDs(c, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			r.handleResponse(c, status_http.NotFound, err.Error())
			return
		}
		r.handleResponse(c, status_http.InternalServerError, err.Error())
		return
	}

	if len(users) == 0 {
		r.handleResponse(c, status_http.NotFound, "user not found")
		return
	}

	r.handleResponse(c, status_http.OK, users[0])
}

// handleResponse handles the HTTP response
func (h *authRoutes) handleResponse(c *gin.Context, status status_http.Status, data interface{}) {
	switch code := status.Code; {
	case code < 400:
	default:
		h.log.Error(
			"response",
			zap.Int("code", status.Code),
			zap.String("status", status.Status),
			zap.Any("description", status.Description),
			zap.Any("data", data),
			zap.Any("custom_message", status.CustomMessage),
		)
	}

	c.JSON(status.Code, status_http.Response{
		Status:        status.Status,
		Description:   status.Description,
		Data:          data,
		CustomMessage: status.CustomMessage,
	})
}
