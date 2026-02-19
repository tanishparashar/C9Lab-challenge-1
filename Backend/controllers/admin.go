package controllers

import (
	"Backend/models"

	beego "github.com/beego/beego/v2/server/web"
)

type AdminController struct {
	beego.Controller
}

// Stats returns admin statistics
// @Title Get Admin Stats
// @Description Get user statistics (admin only)
// @Success 200 {object} models.AdminStats
// @Failure 401 Unauthorized
// @Failure 403 Forbidden
// @router /stats [get]
func (c *AdminController) Stats() {
	totalUsers, err := models.GetTotalUsers()
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Failed to get total users"}
		c.ServeJSON()
		return
	}

	activeUsers, err := models.GetActiveUsers()
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Failed to get active users"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"total_users":  totalUsers,
		"active_users": activeUsers,
	}
	c.ServeJSON()
}
