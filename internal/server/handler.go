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
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/winc-link/hummingbird-http-driver/config"
	constants "github.com/winc-link/hummingbird-http-driver/constant"
	"github.com/winc-link/hummingbird-http-driver/dtos"
	"github.com/winc-link/hummingbird-http-driver/internal/pkg/convert"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"reflect"
	"strconv"
	"time"
)

// offline 设备离线
func offline(c *gin.Context) {
	var response dtos.Response
	deviceId := c.Param(UrlParamDeviceId)
	err := globalDriverService.Offline(deviceId)
	if err != nil {
		response.Code = int(constants.SystemErrorCode)
		response.Success = true
		response.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.SystemErrorCode]) + err.Error()
		encode(response, c.Writer)
		return
	}
	response.Code = int(constants.DefaultSuccessCode)
	response.Success = true
	response.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.DefaultSuccessCode])
	encode(response, c.Writer)
	return
}

// subOffline 子设备离线
func subOffline(c *gin.Context) {
	var response dtos.Response
	deviceId := c.Param(UrlParamDeviceId)
	err := globalDriverService.Offline(deviceId)
	if err != nil {
		response.Code = int(constants.SystemErrorCode)
		response.Success = true
		response.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.SystemErrorCode]) + err.Error()
		encode(response, c.Writer)
		return
	}
	response.Code = int(constants.DefaultSuccessCode)
	response.Success = true
	response.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.DefaultSuccessCode])
	encode(response, c.Writer)
	return
}

// online 设备在线
func online(c *gin.Context) {
	var response dtos.Response
	deviceId := c.Param(UrlParamDeviceId)
	err := globalDriverService.Online(deviceId)
	if err != nil {
		response.Code = int(constants.SystemErrorCode)
		response.Success = true
		response.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.SystemErrorCode]) + err.Error()
		encode(response, c.Writer)
		return
	}
	response.Code = int(constants.DefaultSuccessCode)
	response.Success = true
	response.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.DefaultSuccessCode])
	encode(response, c.Writer)
	return
}

// subOnline 子设备在线
func subOnline(c *gin.Context) {
	var response dtos.Response
	deviceId := c.Param(UrlParamDeviceId)
	err := globalDriverService.Online(deviceId)
	if err != nil {
		response.Code = int(constants.SystemErrorCode)
		response.Success = true
		response.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.SystemErrorCode]) + err.Error()
		encode(response, c.Writer)
		return
	}
	response.Code = int(constants.DefaultSuccessCode)
	response.Success = true
	response.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.DefaultSuccessCode])
	encode(response, c.Writer)
	return
}

// devicePropertyReport 设备属性上报
func devicePropertyReport(c *gin.Context) {
	deviceId := c.Param(UrlParamDeviceId)
	productId := c.Param(UrlParamProductId)
	propertyPost := new(dtos.PropertyPost)
	if err := c.ShouldBind(propertyPost); err != nil {
		encode(dtos.Response{
			Success:      false,
			Code:         int(constants.FormatErrorCode),
			ErrorMessage: string(constants.ErrorCodeMsgMap[constants.FormatErrorCode]),
		}, c.Writer)
	}
	var propertyPostReply dtos.Response
	var delPropertyCode []string
	for code, param := range propertyPost.Params {
		if property, ok := globalDriverService.GetProductPropertyByCode(productId, code); !ok {
			delPropertyCode = append(delPropertyCode, code)
			continue
		} else {
			value := param.Value
			if config.GetConfig().TslParamVerify {
				//推送一条错误消息到客户端
				if verifyErrorCode, verifyErrorMsg := verifyParam(property, param); verifyErrorCode != constants.DefaultSuccessCode {
					delPropertyCode = append(delPropertyCode, code)
					propertyPostReply.Code = int(verifyErrorCode)
					propertyPostReply.Success = false
					propertyPostReply.ErrorMessage = string(verifyErrorMsg)
					encode(propertyPostReply, c.Writer)
					return
				}
			}
			if param.Time == 0 {
				propertyPost.Params[code] = model.PropertyData{
					Time:  time.Now().UnixMilli(),
					Value: value,
				}
			}
		}
	}

	filterPropertyPost := propertyPost.Params
	for _, code := range delPropertyCode {
		delete(filterPropertyPost, code)
	}
	propertyPost.Params = filterPropertyPost
	if len(propertyPost.Params) == 0 {
		propertyPostReply.Code = int(constants.PropertyCodeNotFound)
		propertyPostReply.Success = false
		propertyPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.PropertyCodeNotFound])
		encode(propertyPostReply, c.Writer)
		return
	}
	_, err := globalDriverService.PropertyReport(deviceId, model.NewPropertyReport(true, propertyPost.Params))
	if err != nil {
		propertyPostReply.Code = int(constants.SystemErrorCode)
		propertyPostReply.Success = false
		propertyPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.SystemErrorCode]) + " :" + err.Error()
	} else {
		propertyPostReply.Code = int(constants.DefaultSuccessCode)
		propertyPostReply.Success = true
		propertyPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.DefaultSuccessCode])
	}
	encode(propertyPostReply, c.Writer)
	return
}

// deviceEventReport 设备事件上报
func deviceEventReport(c *gin.Context) {
	deviceId := c.Param(UrlParamDeviceId)
	productId := c.Param(UrlParamProductId)
	eventPost := new(dtos.EventPost)
	if err := c.ShouldBind(eventPost); err != nil {
		encode(dtos.Response{
			Success:      false,
			Code:         int(constants.FormatErrorCode),
			ErrorMessage: string(constants.ErrorCodeMsgMap[constants.FormatErrorCode]),
		}, c.Writer)
	}

	var eventPostReply dtos.Response
	_, ok := globalDriverService.GetProductEventByCode(productId, eventPost.Params.EventCode)
	if !ok {
		eventPostReply.Code = int(constants.EventCodeNotFound)
		eventPostReply.Success = false
		eventPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.EventCodeNotFound]) + fmt.Sprintf(" %s is undefined", eventPost.Params.EventCode)
		encode(eventPost, c.Writer)
		return
	}

	if eventPost.Params.EventTime == 0 {
		eventPost.Params.EventTime = time.Now().UnixMilli()
	}
	_, err := globalDriverService.EventReport(deviceId, model.NewEventReport(true, eventPost.Params))
	if err != nil {
		eventPostReply.Code = int(constants.SystemErrorCode)
		eventPostReply.Success = false
		eventPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.SystemErrorCode]) + " :" + err.Error()
	} else {
		eventPostReply.Code = int(constants.DefaultSuccessCode)
		eventPostReply.Success = true
		eventPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.DefaultSuccessCode])
	}
	encode(eventPostReply, c.Writer)
	return
}

func verifyParam(property model.Property, param model.PropertyData) (constants.ErrorCode, constants.ErrorMessage) {
	switch property.TypeSpec.Type {
	case "int":
		var intOrFloatSpecs dtos.IntOrFloatSpecs
		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &intOrFloatSpecs)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		t, err := convert.GetInterfaceToInt(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an int type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		min, _ := convert.StringToInt(intOrFloatSpecs.Min)
		max, _ := convert.StringToInt(intOrFloatSpecs.Max)

		if t < min || t > max {
			return constants.ReportDataRangeErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataRangeErrorCode]
		}
	case "float":
		var intOrFloatSpecs dtos.IntOrFloatSpecs
		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &intOrFloatSpecs)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		t, err := convert.GetInterfaceToFloat64(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an float type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		min, _ := convert.StringToFloat64(intOrFloatSpecs.Min)
		max, _ := convert.StringToFloat64(intOrFloatSpecs.Max)
		if t < min || t > max {
			return constants.ReportDataRangeErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataRangeErrorCode]
		}
	case "bool":
		t, err := convert.GetInterfaceToFloat64(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an int type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		if !(t == 0 || t == 1) {
			return constants.ReportDataRangeErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataRangeErrorCode]
		}
	case "text":
		var textSpecs dtos.TextSpecs
		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &textSpecs)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		t, err := convert.GetInterfaceToString(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an string type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		length, err := strconv.Atoi(textSpecs.Length)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		if length < len(t) {
			return constants.ReportDataLengthErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataLengthErrorCode]
		}
	case "enum":
		enumSpecs := make(map[string]string)
		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &enumSpecs)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		t, err := convert.GetInterfaceToFloat64(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an int type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		if _, ok := enumSpecs[strconv.Itoa(int(t))]; !ok {
			return constants.ReportDataRangeErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataRangeErrorCode]
		}
	}
	return constants.DefaultSuccessCode,
		constants.ErrorCodeMsgMap[constants.DefaultSuccessCode]
}
