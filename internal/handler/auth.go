package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"restaurant_booking/internal/auth"
)

const sessionKey = "user_id"

func RegisterAuth(r *gin.Engine, db *sqlx.DB) {
	r.GET("/login", loginFrom)
	r.POST("/login", doLogin(db))
	r.GET("/logout", doLogout)
}

// --- Middleware ---

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		if sess.Get(sessionKey) == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func loginFrom(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{"error": ""})
}

func doLogin(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds struct {
			Email string `from:"email" binding:"required,email"`
			Pw    string `from:"password" binding:"required"`
		}
		if c.ShouldBind(&creds) != nil {
			c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"error": "invalid input"})
			return
		}

		var u auth.User
		if err := db.Get(&u, `SELECT * FROM admin_users WHERE email=$1`,
			creds.Email); err != nil {
			c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"error": "User not found"})
			return
		}
		if !auth.Verify(u.PassHash, creds.Pw) {
			c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"error": "Invalid password"})
			return
		}

		sess := sessions.Default(c)
		sess.Set(sessionKey, u.ID)
		_ = sess.Save()

		c.Redirect(http.StatusSeeOther, "/admin")
	}
}

func doLogout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	_ = sess.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}
