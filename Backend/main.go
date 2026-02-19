package main

import (
	"Backend/database"
	"Backend/middleware"
	_ "Backend/routers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func main() {
	// Initialize database connection
	database.InitDB()
	defer database.CloseDB()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// Enable CORS for development — allow Vite default dev origin (localhost:5173)
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
		origin := ctx.Input.Header("Origin")
		// Allowed dev origins
		allowed := map[string]bool{
			"http://localhost:5173": true,
			"http://127.0.0.1:5173": true,
		}

		if origin != "" && (allowed[origin] || beego.BConfig.RunMode == "dev") {
			// Echo the origin (required when Access-Control-Allow-Credentials=true)
			ctx.Output.Header("Access-Control-Allow-Origin", origin)
			ctx.Output.Header("Vary", "Origin")
			ctx.Output.Header("Access-Control-Allow-Credentials", "true")
		} else if beego.BConfig.RunMode == "dev" {
			// fallback for other dev origins
			ctx.Output.Header("Access-Control-Allow-Origin", "*")
		}

		ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Output.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		ctx.Output.Header("Access-Control-Expose-Headers", "Authorization")
		// fast response for preflight — write directly to the ResponseWriter so Beego routing doesn't override it
		if ctx.Input.Method() == "OPTIONS" {
			ctx.ResponseWriter.Header().Set("Content-Length", "0")
			ctx.ResponseWriter.WriteHeader(204)
			return
		}
	})

	// Protect /api/auth/me — registered AFTER CORS filter so preflight is never blocked
	beego.InsertFilter("/api/auth/me", beego.BeforeRouter, middleware.AuthMiddleware)

	beego.Run()
}
