package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

func WebDAVPerUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetUint("userId")
		if userId == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Create WebDAV handler for the user with userId
		webdavRoot := fmt.Sprintf("./uploads/fileserver/%d", userId)
		handler := &webdav.Handler{
			Prefix:     "/webdav/",
			FileSystem: webdav.Dir(webdavRoot),
			LockSystem: webdav.NewMemLS(),
		}

		http.Handler(handler).ServeHTTP(c.Writer, c.Request)
	}
}
