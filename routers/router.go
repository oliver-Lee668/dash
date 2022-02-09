package routers

import (
	"net/http"

	"github.com/cnych/dash/controllers"
	"github.com/gin-gonic/gin"
)

func InitApi(eng *gin.Engine) {
	// gin 配置使用中间件
	//eng.Use(CorsMiddleware)

	// 注册一个check health的接口
	eng.GET("/ping", controllers.Ping)

	// 接口分组
	api := eng.Group("/api/v1")
	// 获取Node列表的接口
	api.GET("nodes", controllers.GetNodeList)
	// 获取Metrics指标数据
	api.POST("metrics", controllers.GetMetrics)
	// 添加namespace
	api.POST("namespace", controllers.AddNamespace)
	// 获取namespace列表
	api.GET("namespaces", controllers.GetNamespaceList)
	// 删除namespace
	api.DELETE("namespace/:name", controllers.DeleteNamespace)
	// 添加pod
	api.POST("pod", controllers.AddPod)
	// 删除pod
	api.DELETE("pod", controllers.DeletePod)
	// 获取所有pod
	api.GET("pods", controllers.GetPodList)
	// 获取Pod列表
	api.GET("namespaces/:namespace/pods", controllers.GetKubePods)
	// 获取Pod(容器)日志
	api.GET("namespaces/:namespace/pods/:pod/logs", controllers.GetKubeLogs)
	// 执行Pod(容器)命令
	//ws://127.0.0.1:8888/api/v1/namespaces/kube-ops/pods/gitlab-f4d95db8-fj24z/shell
	api.GET("namespaces/:namespace/pods/:pod/shell", controllers.HandleTerminal)
	// 获取dynamic列表
	api.GET("dynamics", controllers.GetDynamicList)
	// 获取dynamic
	api.GET("dynamic", controllers.GetDynamic)
	api.POST("dynamic", controllers.AddDynamic)

}

// CorsMiddleware 允许跨域的一个中间件
func CorsMiddleware(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, access-control-allow-origin, Origin, X-Requested-With, Content-Type, Accept, Content-Length, Accept-Encoding, Content-Range, Content-Disposition, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Set("content-type", "application/json")
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.Next()
}
