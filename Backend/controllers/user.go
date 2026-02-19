package controllers

import (
	"Backend/models"

	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

// Summary returns user summary information
// @Title Get User Summary
// @Description Get current user's summary (user role only)
// @Success 200 {object} models.UserSummary
// @Failure 401 Unauthorized
// @Failure 403 Forbidden
// @router /summary [get]
func (c *UserController) Summary() {
	userId := c.Ctx.Input.GetData("userId").(int)

	user, err := models.GetUserById(userId)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "User not found"}
		c.ServeJSON()
		return
	}

	var lastLogin interface{}
	if user.LastLogin != nil {
		lastLogin = user.LastLogin.Format("2006-01-02T15:04:05Z07:00")
	} else {
		lastLogin = nil
	}

	c.Data["json"] = map[string]interface{}{
		"last_login": lastLogin,
		"message":    "Welcome, " + user.Name + "!",
	}
	c.ServeJSON()
}

