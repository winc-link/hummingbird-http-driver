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
	"github.com/gin-gonic/gin"
	constants "github.com/winc-link/hummingbird-http-driver/constant"
	"github.com/winc-link/hummingbird-http-driver/dtos"
	"github.com/winc-link/hummingbird-http-driver/internal/pkg/tool"
)

func authHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Request.RequestURI)
		deviceId := c.Param(UrlParamDeviceId)
		productId := c.Param(UrlParamProductId)

		if deviceId == "" || productId == "" {
			encode(dtos.Response{
				Success:      false,
				Code:         int(constants.InvalidParameterErrorCode),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.InvalidParameterErrorCode]),
			}, c.Writer)
			c.Abort()
			return
		}

		device, ok := globalDriverService.GetDeviceById(deviceId)
		if !ok {
			encode(dtos.Response{
				Success:      false,
				Code:         int(constants.DeviceNotFound),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.DeviceNotFound]),
			}, c.Writer)
			c.Abort()
			return
		}

		product, ok := globalDriverService.GetProductById(productId)

		if !ok {
			encode(dtos.Response{
				Success:      false,
				Code:         int(constants.ProductNotFound),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.ProductNotFound]),
			}, c.Writer)
			c.Abort()
			return
		}

		if device.ProductId != product.Id {
			encode(dtos.Response{
				Success:      false,
				Code:         int(constants.InvalidParameterErrorCode),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.InvalidParameterErrorCode]),
			}, c.Writer)
			c.Abort()
			return
		}

		token := c.Request.Header.Get("token")
		if token == "" {
			encode(dtos.Response{
				Success:      false,
				Code:         int(constants.AuthPermissionDenyErrorCode),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.AuthPermissionDenyErrorCode]),
			}, c.Writer)
			c.Abort()
			return
		}
		if token != tool.HmacMd5(device.Secret, device.Id+"&"+product.Key) {
			encode(dtos.Response{
				Success:      false,
				Code:         int(constants.AuthPermissionDenyErrorCode),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.AuthPermissionDenyErrorCode]),
			}, c.Writer)
			c.Abort()
			return
		}
		c.Next()
	}
}
