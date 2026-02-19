// Package routers provides API route definitions
package routers

import (
	"devops-api/controllers"

	"devops-api/middleware"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// API namespace
	ns := beego.NewNamespace("/api",
		// Auth endpoints (public)
		beego.NSNamespace("/auth",
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
			beego.NSRouter("/refresh", &controllers.AuthController{}, "post:Refresh"),
			beego.NSRouter("/me", &controllers.AuthController{}, "get:Me"),
		),

		// Admin endpoints (requires admin role)
		beego.NSNamespace("/admin",
			beego.NSBefore(middleware.AuthMiddleware, middleware.AdminMiddleware),
			beego.NSRouter("/stats", &controllers.AdminController{}, "get:Stats"),
		),

		// User endpoints (requires user role)
		beego.NSNamespace("/user",
			beego.NSBefore(middleware.AuthMiddleware, middleware.UserMiddleware),
			beego.NSRouter("/summary", &controllers.UserController{}, "get:Summary"),
		),
	)

	beego.AddNamespace(ns)
}
