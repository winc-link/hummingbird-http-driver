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
	"net/http"
)

const (
	UrlParamDeviceId  = "deviceId"
	UrlParamProductId = "productId"
)

const (
	ContentType     = "Content-Type"
	ContentTypeCBOR = "application/cbor"
	ContentTypeJSON = "application/json"
	ContentTypeYAML = "application/x-yaml"
	ContentTypeText = "text/plain"
	ContentTypeXML  = "application/xml"
	AcceptLanguage  = "Accept-Language"
)

func encode(i interface{}, w http.ResponseWriter) {
	w.Header().Add(ContentType, ContentTypeJSON)
	enc := json.NewEncoder(w)
	err := enc.Encode(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
