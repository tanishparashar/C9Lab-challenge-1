package controllers

import (
	"Backend/models"
	"Backend/utils"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	beego.Controller
}

// Login handles user authentication
// @Title Login
// @Description User login
// @Param	body		body 	models.LoginRequest	true	"Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 Invalid request
// @Failure 401 Invalid credentials
// @router /login [post]
func (c *AuthController) Login() {
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &loginReq)
	if err != nil || loginReq.Email == "" || loginReq.Password == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Email and password are required"}
		c.ServeJSON()
		return
	}

	// Get user from database
	user, err := models.GetUserByEmail(loginReq.Email)
	if err != nil {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Invalid email or password"}
		c.ServeJSON()
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Account is inactive"}
		c.ServeJSON()
		return
	}

	// Verify password
	if !user.CheckPassword(loginReq.Password) {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Invalid email or password"}
		c.ServeJSON()
		return
	}

	// Update last login
	models.UpdateLastLogin(user.Id)

	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateToken(user.Id, user.Email, user.Role)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Failed to generate tokens"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"access":  accessToken,
		"refresh": refreshToken,
		"user":    user.ToDTO(),
	}
	c.ServeJSON()
}

// Refresh handles token refresh
// @Title Refresh Token
// @Description Refresh access token
// @Param	body		body 	models.RefreshRequest	true	"Refresh token"
// @Success 200 {object} models.RefreshResponse
// @Failure 400 Invalid request
// @Failure 401 Invalid token
// @router /refresh [post]
func (c *AuthController) Refresh() {
	var refreshReq struct {
		Refresh string `json:"refresh"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &refreshReq)
	if err != nil || refreshReq.Refresh == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Refresh token is required"}
		c.ServeJSON()
		return
	}

	// Parse and validate refresh token
	claims, err := utils.ParseToken(refreshReq.Refresh)
	if err != nil {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Invalid or expired refresh token"}
		c.ServeJSON()
		return
	}

	// Generate new access token
	accessToken, _, err := utils.GenerateToken(claims.UserId, claims.Email, claims.Role)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Failed to generate access token"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{
		"access": accessToken,
	}
	c.ServeJSON()
}

// Me returns the current user's information
// @Title Get Current User
// @Description Get authenticated user info
// @Success 200 {object} models.UserDTO
// @Failure 401 Unauthorized
// @router /me [get]
func (c *AuthController) Me() {
	userId := c.Ctx.Input.GetData("userId").(int)

	user, err := models.GetUserById(userId)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "User not found"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = user.ToDTO()
	c.ServeJSON()
}
