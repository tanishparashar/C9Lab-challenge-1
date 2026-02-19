package middleware

import (
	"Backend/utils"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// AuthMiddleware checks for valid JWT token
func AuthMiddleware(ctx *context.Context) {
	// Allow preflight requests through — CORS filter handles them
	if ctx.Input.Method() == "OPTIONS" {
		return
	}

	authHeader := ctx.Input.Header("Authorization")
	if authHeader == "" {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]string{"error": "Authorization header required"}, false, false)
		return
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]string{"error": "Invalid authorization header format"}, false, false)
		return
	}

	tokenString := parts[1]
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]string{"error": "Invalid or expired token"}, false, false)
		return
	}

	// Store user info in context
	ctx.Input.SetData("userId", claims.UserId)
	ctx.Input.SetData("email", claims.Email)
	ctx.Input.SetData("role", claims.Role)
}

// AdminMiddleware checks if user has admin role
func AdminMiddleware(ctx *context.Context) {
	if ctx.Input.Method() == "OPTIONS" {
		return
	}
	role := ctx.Input.GetData("role")
	if role == nil || role.(string) != "admin" {
		ctx.Output.SetStatus(403)
		ctx.Output.JSON(map[string]string{"error": "Admin access required"}, false, false)
		return
	}
}

// UserMiddleware checks if user has user role
func UserMiddleware(ctx *context.Context) {
	if ctx.Input.Method() == "OPTIONS" {
		return
	}
	role := ctx.Input.GetData("role")
	if role == nil || role.(string) != "user" {
		ctx.Output.SetStatus(403)
		ctx.Output.JSON(map[string]string{"error": "User access required"}, false, false)
		return
	}
}
