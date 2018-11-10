package main

import (
	"github.com/gin-gonic/gin"
	"gogsc/util"
	"fmt"
	"gogsc/controller"
)

func setRoute(r *gin.Engine){
	r.GET("/songci/index/:id/:open_id", controller.HandleIndex)
	r.GET("/songci/query/:q/:page/:open_id", controller.HandleQuery)
	r.GET("/user/auth/:code", controller.Code2Session)
	r.GET("/user/like/:open_id/:gsc_id", controller.SetUserLike)
	r.GET("/user/dislike/:open_id/:gsc_id", controller.SetUserDisLike)
	r.GET("/songci/mylike/:open_id", controller.QueryMyLike)

}
func main()  {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	setRoute(r)
	listenAddr := fmt.Sprintf("%v", util.GetConf("listenAddr"))
	r.Run(listenAddr)
}
