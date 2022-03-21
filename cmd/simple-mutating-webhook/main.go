package main

import (
	"fmt"
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
	"time"
)

func main()  {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	route.InitDeploymentRoute(r)

	return r
}

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			funcSlices := strings.Split(f.Function, "/")

			lineInfo := fmt.Sprintf("%s:%d", filename, f.Line)
			fileInfo := fmt.Sprintf("%s()", funcSlices[len(funcSlices)-1])

			return fmt.Sprintf("%20s] %25s]", lineInfo, fileInfo), time.Now().Format("0102 15:04:05.000000")
		},
	})
}