package middleware

import (
	"database/sql"
	"fmt"
	"golang-final/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func authorizeWithRole(role db.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetHeader("username")
		password := ctx.GetHeader("password")
		if username == "" || password == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No Authorization headers found",
			})
			return
		}
		user_info, err := db.AuthoriseUser(username, password)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "No user with this credintails",
				})
			} else {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Internal error: " + err.Error(),
				})
			}
			return
		}
		if role != db.ANY_ROLE && role != user_info.Role {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User must have role: " + role,
			})
			return
		}
		fmt.Printf("Login with username: %s and password %s\n", username, password)
		ctx.Set("user_info", user_info)
	}
}

func Authorise() gin.HandlerFunc {
	return authorizeWithRole(db.ANY_ROLE)
}

func AuthoriseAdmin() gin.HandlerFunc {
	return authorizeWithRole(db.ADMIN)
}

func AuthoriseSeller() gin.HandlerFunc {
	return authorizeWithRole(db.SELLER)
}

func AuthoriseBuyer() gin.HandlerFunc {
	return authorizeWithRole(db.BUYER)
}
