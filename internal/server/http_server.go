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
	"fmt"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"net/http"
)

type HttpServer struct {
	sd *service.DriverService
}

func NewHttpService(sd *service.DriverService) *HttpServer {
	return &HttpServer{
		sd: sd,
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello hummingbird driver http server!")
}

func (c *HttpServer) Start() {

	http.HandleFunc("/", index)
	// 启动web服务，用户可以根据项目需要添加相关路由和修改监听端口
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		c.sd.GetLogger().Error("ListenAndServe: ", err)
	}
}
