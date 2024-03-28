package util

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type CommonError struct {
	apiName     string
	Description string `json:"description"`
	ErrCode     int64  `json:"error_code"`
	Message     string `json:"message"`
}

func (c *CommonError) Error() string {
	return fmt.Sprintf("%s Error , errcode=%d , errmsg=%s", c.apiName, c.ErrCode, c.Message)
}

// NewCommonError 新建 CommonError 错误，对于无 errcode 和 errmsg 的返回也可以返回该通用错误
func NewCommonError(apiName string, code int64, msg string) *CommonError {
	return &CommonError{
		apiName: apiName,
		ErrCode: code,
		Message: msg,
	}
}

// DecodeWithError 将返回值按照解析
func DecodeWithError(response []byte, obj interface{}, apiName string) error {
	// fmt.Println(string(response))
	err := json.Unmarshal(response, obj)
	if err != nil {
		return fmt.Errorf("json Unmarshal Error, err=%v", err)
	}
	responseObj := reflect.ValueOf(obj)
	// fmt.Println(responseObj.Elem())
	if !responseObj.IsValid() {
		return fmt.Errorf("obj is invalid")
	}
	data := responseObj.Elem().FieldByName("Data")
	if !data.IsValid() || (data.Kind() != reflect.Struct && data.Kind() != reflect.Interface) {
		return fmt.Errorf("data is invalid or not struct %v", data.Kind())
	}
	// fmt.Println(data)
	// dataStruct := reflect.ValueOf(data)
	// if !dataStruct.IsValid() {
	// 	return fmt.Errorf("dataStruct is invalid or not struct %v", dataStruct)
	// }
	commonError := &CommonError{}
	err = mapstructure.Decode(data.Interface(), commonError)
	if err != nil {
		return fmt.Errorf("commonError is invalid or not struct %v", err)
	}
	// if !commonError.IsValid() || commonError.Kind() != reflect.Struct {
	// 	return fmt.Errorf("commonError is invalid or not struct %v", commonError)
	// }
	errCode := commonError.ErrCode
	errMsg := commonError.Message
	// if !errCode.IsValid() || !errMsg.IsValid() {
	// 	return fmt.Errorf("errcode or errmsg is invalid")
	// }
	if errCode != 0 {
		return &CommonError{
			apiName: apiName,
			ErrCode: errCode,
			Message: errMsg,
		}
	}
	return nil
}
