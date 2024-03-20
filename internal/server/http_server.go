/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package server

import (
	"github.com/gin-gonic/gin"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"net/http"
	"time"
)

var globalDriverService *service.DriverService

type HttpServer struct {
	//sd *service.DriverService
}

func NewHttpService(sd *service.DriverService) *HttpServer {
	globalDriverService = sd
	return &HttpServer{
		//sd: sd,
	}
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	auth := r.Group("/api/v1")
	auth.Use(authHandler())
	auth.POST("/device/online/:deviceId/:productId", online)
	auth.POST("/device/sub/online/:deviceId/:productId", subOnline)
	auth.POST("/device/offline/:deviceId/:productId", offline)
	auth.POST("/device/sub/offline/:deviceId/:productId", subOffline)
	auth.POST("/device/thing/property/post/:deviceId/:productId", devicePropertyReport)
	auth.POST("/device/thing/event/post/:deviceId/:productId", deviceEventReport)

	return r
}

func (c *HttpServer) Start() *http.Server {
	route := setupRouter()
	timeout := time.Millisecond * time.Duration(5000)
	server := &http.Server{
		Addr:         "0.0.0.0:8090",
		Handler:      route,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			globalDriverService.GetLogger().Errorf("Web server start failed: %v", err)
		}
	}()

	globalDriverService.GetLogger().Infof("Web server start successful [::%d]", 8090)
	return server
}
