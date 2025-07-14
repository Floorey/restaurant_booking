package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func RegisterAdmin(r *gin.RouterGroup, db *sqlx.DB) {
	r.GET("/", adminList(db))
	r.POST("/canel/:id", adminCancel(db))
}

func adminList(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bookings []struct {
			ID        string `db:"id"`
			GuestMail string `db:"guest_email"`
		}
		db.Select(&bookings,
			`SELECT id, guest_email FROM bookings ORDER BY booking_time`)
		c.HTML(http.StatusOK, "admin_list.tmpl",
			gin.H{"bookings": bookings})
	}
}

func adminCancel(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _ = db.Exec(`UPDATE bookings
						SET status='canceled'
						WHERE id=$1`, c.Param("id"))
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}
