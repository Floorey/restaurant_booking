package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"restaurant_booking/internal/model"
	"strconv"
)

const pageSize = 50

func RegisterAdmin(r *gin.RouterGroup, db *sqlx.DB) {
	r.Use(RequireAuth())
	r.GET("/", adminList(db))
	r.POST("/cancel/:id", adminCancel(db))
}

func adminList(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.DefaultQuery("status", "all")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}

		args := []interface{}{(page - 1) * pageSize, pageSize}
		query := `SELECT * FROM bookings `
		if status != "all" {
			query += `WHERE status = $3 `
			args = append([]interface{}{status}, args...)
		}
		query += `ORDER BY booking_time DESC LIMIT $2 OFFSET $1`

		var list []model.Booking
		if err := db.Select(&list, query, args...); err != nil {
			c.String(http.StatusInternalServerError, "DB-Error")
			return
		}

		c.HTML(http.StatusOK, "admin.tmpl", gin.H{
			"list":   list,
			"status": status,
			"page":   page,
		})
	}
}

func adminCancel(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		db.Exec(`UPDATE bookings SET status='canceled' WHERE id=$1`, id)
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}
