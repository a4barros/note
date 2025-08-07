package main

import (
	"net/http"
	"note/db"
	"note/middleware"
	"note/noteConfig"
	"note/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	noteConfig.ConfigInit()

	if noteConfig.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.MaxMultipartMemory = 256 << 20 // 256MB file max
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{noteConfig.GetDomain()},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/api/login", routes.Login)
	r.GET("/api/version", routes.Version)
	internal := r.Group("/api")

	internal.Use(middleware.AuthMiddleware())
	internal.POST("/logout", routes.Logout)
	internal.POST("/note", routes.PostNote)
	internal.PUT("/note/:id", routes.PutNote)
	internal.GET("/note", routes.GetAllNotes)
	internal.GET("/isLogged", routes.IsLogged)
	internal.GET("/category", routes.GetAllCategories)
	internal.GET("/category/hidden", routes.GetAllCategoriesWithHidden)
	internal.POST("/note/category", routes.GetNotesByCategory)
	internal.DELETE("/note/:id", routes.DeleteNote)
	internal.POST("/note/:id/attachment", routes.PostAttachment)
	internal.DELETE("/note/:id/attachment/:attachmentId", routes.DeleteAttachment)
	internal.GET("/note/:id/attachment/:attachmentId/file", routes.GetAttachmentFile)

	fileserverInternal := internal.Group("/fileserver")
	fileserverInternal.Use(middleware.GetPathMiddleware())
	fileserverInternal.GET("/ls", routes.Ls)
	fileserverInternal.POST("/mkdir", routes.Mkdir)
	fileserverInternal.PUT("/mv", routes.Mv)
	fileserverInternal.POST("/write", routes.WriteFile)
	fileserverInternal.DELETE("/rm", routes.Rm)
	fileserverInternal.GET("/read", routes.ReadFile)

	r.OPTIONS("/webdav/*path", func(c *gin.Context) {
		c.Header("Allow", "OPTIONS, PROPFIND, GET, PUT, POST, DELETE, MKCOL, COPY, MOVE, LOCK, UNLOCK")
		c.Header("DAV", "1, 2")
		c.Status(http.StatusOK)
	})

	webdavRoutes := r.Group("/webdav", middleware.WebDavAuthMiddleware())
	{
		webdavRoutes.Handle(http.MethodGet, "/*path", middleware.WebDAVPerUserHandler())
		webdavRoutes.Handle(http.MethodPut, "/*path", middleware.WebDAVPerUserHandler())
		webdavRoutes.Handle(http.MethodPost, "/*path", middleware.WebDAVPerUserHandler())
		webdavRoutes.Handle(http.MethodDelete, "/*path", middleware.WebDAVPerUserHandler())

		webdavRoutes.Handle("PROPFIND", "/*path", middleware.WebDAVPerUserHandler())
		webdavRoutes.Handle("MKCOL", "/*path", middleware.WebDAVPerUserHandler())
		webdavRoutes.Handle("PROPPATCH", "/*path", middleware.WebDAVPerUserHandler())
		webdavRoutes.Handle("COPY", "/*path", middleware.WebDAVPerUserHandler())
		webdavRoutes.Handle("MOVE", "/*path", middleware.WebDAVPerUserHandler())
		webdavRoutes.Handle("LOCK", "/*path", middleware.WebDAVPerUserHandler())
		webdavRoutes.Handle("UNLOCK", "/*path", middleware.WebDAVPerUserHandler())
	}

	r.Run(":5003")
}
