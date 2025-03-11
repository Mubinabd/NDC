package handlers

import (
	"net/http"
	"posts/internal/pkg/config"
	pb "posts/internal/pkg/genproto"
	"strconv"

	"context"
	"github.com/gin-gonic/gin"
	"posts/internal/pkg/token"
)

// Login godoc
// @Summary Login a user
// @Description Authenticate user with email and password
// @Tags login
// @Accept json
// @Produce json
// @Param credentials body pb.LoginRequest true "User login credentials"
// @Success 200 {object} token.Tokens "JWT tokens"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Invalid email or password"
// @Router /api/v1/login [post]
func (h *Handler) Login(c *gin.Context) {
	req := pb.LoginRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := h.Clients.User.Login(c, &pb.LoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or incorrect password"})
		return
	}

	tokens := token.GenerateJWTToken(user.Id, user.Email, user.Token)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// CreateUser creates a new user
// @Summary Create User
// @Description Create a new user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body pb.UserCreateRequest true "User data"
// @Success 200 {string} string pb.UserCreateResponse "User created successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/create [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req pb.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.ERROR.Println("Failed to bind request:", err)
		c.JSON(400, "Invalid request: "+err.Error())
		return
	}

	res, err := h.Clients.User.Create(c, &req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to create user:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, res)
}

// GetUserDetail retrieves a user by ID
// @Summary Get User
// @Description Retrieve a user by its ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} pb.UserGetResponse "User details"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/{id} [get]
func (h *Handler) GetUserDetail(c *gin.Context) {
	userID := c.Param("id")

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		h.Logger.ERROR.Println("Invalid log ID:", err)
		c.JSON(400, "Invalid log ID")
		return
	}
	req := &pb.ById{Id: id}
	res, err := h.Clients.User.GetDetail(c, req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to get user:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, res)
}

// UpdateUser updates an existing user
// @Summary Update User
// @Description Update an existing user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body pb.UserUpdateRequest true "User data"
// @Success 200 {string} string "User updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/update [patch]
func (h *Handler) UpdateUser(c *gin.Context) {
	var body pb.UserUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.Logger.ERROR.Println("Failed to bind request:", err)
		c.JSON(400, "Invalid request: "+err.Error())
		return
	}

	req := pb.UserUpdateRequest{
		Id:          body.Id,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		Username:    body.Username,
		PhoneNumber: body.PhoneNumber,
		Gender:      body.Gender,
		Role:        body.Role,
	}

	_, err := h.Clients.User.Update(c, &req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to update user:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, "User updated successfully")
}

// DeleteUser deletes a user by ID
// @Summary Delete User
// @Description Delete a user by its ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		h.Logger.ERROR.Println("Invalid log ID:", err)
		c.JSON(400, "Invalid log ID")
		return
	}

	req := &pb.ById{Id: id}
	_, err = h.Clients.User.Delete(c, req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to delete user:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, "User deleted successfully")
}

// GetUserList retrieves a list of users with pagination
// @Summary Get Users
// @Description Retrieve a list of users with pagination
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param page query int false "Page"
// @Param role query string false "Role"
// @Param username query string false "Username"
// @Param firstname query string false "FirstName"
// @Param lastname query string false "LastName"
// @Param gender query string false "Gender"
// @Success 200 {object} pb.UserGetAll "List of users"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/list [get]
func (h *Handler) GetUserList(c *gin.Context) {
	var filter pb.FilterUser

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = int64(l)
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = int64(o)
		}
	}
	if firstname := c.Query("first_name"); firstname != "" {
		filter.Firstname = firstname
	}
	if lastname := c.Query("last_name"); lastname != "" {
		filter.Lastname = lastname
	}
	if username := c.Query("username"); username != "" {
		filter.Username = username
	}
	if gender := c.Query("gender"); gender != "" {
		filter.Gender = gender
	}

	res, err := h.Clients.User.GetList(context.Background(), &filter)
	if err != nil {
		h.Logger.ERROR.Println("Failed to list users:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, res)
}

// ChangeUserPassword godoc
// @Summary Change password
// @Description Updates the password to new one
// @Tags User
// @Accept json
// @Produce json
// @Param request body pb.UserRecoverPasswordRequest true "Change Password Request"
// @Success 200 {object} string "Password successfully updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Incorrect verification code"
// @Failure 404 {object} string "Verification code expired or email not found"
// @Failure 500 {object} string "Error updating password"
// @Security BearerAuth
// @Router /api/v1/users/user-password [put]
func (h *Handler) ChangeUserPassword(c *gin.Context) {
	var req pb.UserRecoverPasswordRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	if req.Email == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and new password are required fields."})
		return
	}

	if err := config.IsValidPassword(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := config.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't hash your password", "details": err.Error()})
		return
	}

	_, err = h.Clients.User.ChangeUserPassword(c, &pb.UserRecoverPasswordRequest{Email: req.Email, NewPassword: hashedPassword})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully updated"})
}
