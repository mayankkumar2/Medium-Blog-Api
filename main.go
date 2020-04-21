package main

import (
	"github.com/gin-gonic/gin"
	"medium-feed-2/utils"
	"github.com/gin-contrib/cors"
)

func getArticles(e *gin.Context) {
	orgid  := e.Request.URL.Query()["orgid"]
	if orgid == nil {
		e.JSON(400, map[string] interface{} {
			"status": "failed",
			"error": "orgid not supplied. Need",
		})
	} else {
		res, statuscode := utils.GetArticlesforOrganizationV2(orgid[0])
		e.JSON(statuscode, res)
	}
}

func getArticlesV1(e *gin.Context){
	orgid := e.Request.URL.Query()["orgid"]
	userid := e.Request.URL.Query()["userid"]
	if orgid == nil && userid == nil {
		e.JSON(400, map[string] interface{} {
			"status": "failed",
			"error": "orgid OR userid not supplied. Need",
		})
	} else if orgid != nil {
		resp,statuscode := utils.GetArticles(orgid[0]);
		e.Header("Content-Type","application/json")
		e.String(statuscode,"%s",(resp))
	} else if userid != nil {
		resp,statuscode := utils.GetArticles("@"+userid[0]);
		e.Header("Content-Type","application/json")
		e.String(statuscode,"%s",(resp))
	}
}
func main() {
	router := gin.Default()
	router.Use(cors.Default())
	api := router.Group("/api")
	{
		v2 := api.Group("/v2")
		{
			v2.GET("/articles",getArticles)
		}
		v1 := api.Group("/v1")
		{
			v1.GET("/articles",getArticlesV1)
		}

	}
	router.Run()
}