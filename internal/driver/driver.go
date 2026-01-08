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

package driver

import (
	"context"
	"errors"
	"github.com/spf13/cast"
	"github.com/winc-link/hummingbird-http-driver/internal/server"
	"github.com/winc-link/hummingbird-sdk-go/commons"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"net/http"
	"time"
)

type HttpProtocolDriver struct {
	sd         *service.DriverService
	httpServer *http.Server
}

func (dr HttpProtocolDriver) HandlePropertyReportDebug(ctx context.Context, deviceId string, data model.PropertyReport) error {
	data.Time = time.Now().UnixMilli()
	newData := make(map[string]interface{})
	for k, v := range data.Data {
		newData[k] = cast.ToFloat64(v)
	}
	data.Data = newData
	resp, _ := dr.sd.PropertyReport(deviceId, data)
	if resp.Success != true {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}

func (dr HttpProtocolDriver) HandleEventReportDebug(ctx context.Context, deviceId string, data model.EventReport) error {
	data.Time = time.Now().UnixMilli()
	resp, _ := dr.sd.EventReport(deviceId, data)
	if resp.Success != true {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}

// DeviceNotify 设备添加/修改/删除通知
func (dr HttpProtocolDriver) DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, deviceId string, device model.Device) error {
	return nil
}

// ProductNotify 产品添加/修改/删除通知
func (dr HttpProtocolDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, productId string, product model.Product) error {
	return nil
}

// Stop 驱动退出通知。
func (dr HttpProtocolDriver) Stop(ctx context.Context) error {
	for _, d := range dr.sd.GetDeviceList() {
		err := dr.sd.Offline(d.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// HandlePropertySet 设备属性设置
func (dr HttpProtocolDriver) HandlePropertySet(ctx context.Context, deviceId string, data model.PropertySet) error {
	return nil
}

// HandlePropertyGet 设备属性查询
func (dr HttpProtocolDriver) HandlePropertyGet(ctx context.Context, deviceId string, data model.PropertyGet) error {
	return nil
}

// HandleServiceExecute 设备服务调用
func (dr HttpProtocolDriver) HandleServiceExecute(ctx context.Context, deviceId string, data model.ServiceExecuteRequest) error {
	return nil
}

// NewHttpProtocolDriver Http协议驱动
func NewHttpProtocolDriver(sd *service.DriverService) *HttpProtocolDriver {
	httpServer := server.NewHttpService(sd).Start()
	h := &HttpProtocolDriver{
		sd:         sd,
		httpServer: httpServer,
	}
	return h
}
