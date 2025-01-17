package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gohub/internal/api/user"
	"gohub/internal/api/user/repository/dao"
	"gohub/internal/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func Init() *gin.Engine {
	server := gin.Default()
	database, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/gohub"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(database)
	if err != nil {
		panic(err)
	}

	// cors
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		// AllowAllOrigins:  true,
		// AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "https://gohub.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// set session and cookie
	cookieStore := cookie.NewStore([]byte("gohub-secret"))
	server.Use(sessions.Sessions("gohub-session", cookieStore))

	// check sign in status
	server.Use(
		middleware.
			InitAuthMiddlewareBuilder().
			AppendIgnorePath("/v1/api/user/sign-in").
			AppendIgnorePath("/v1/api/user/sign-up").
			Builder(),
	)

	v1Server := server.Group("/v1/api")

	user.Init(database, v1Server)

	return server
}
