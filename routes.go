package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func initializeRoutes() {
	router.PUT("/buildstatus/:id", verifyClient(saveBuildStatus))
	router.PUT("/lights/config/:id", verifyClient(putLightConfig))
	router.GET("/lights/status", verifyClient(getLightStatus))

	router.GET("/lights/config/:id", verifyClient(getLight))
	router.GET("/lights/config", verifyClient(getLightConfigs))
}

func verifyClient(handler gin.HandlerFunc) gin.HandlerFunc {
	// TODO should implement better user management... Use API gateway to authenticate and throttle?
	return gin.HandlerFunc(func(ctx *gin.Context) {

		// TODO Validate Client ID is a registered ID
		clientId := ctx.GetHeader("X_CLIENT_ID")

		if clientId == "" {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		handler(ctx)
	})
}
