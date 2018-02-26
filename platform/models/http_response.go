package models

import "net/http"

type Success struct {
	Code   int `json:"code"`
	Result struct {
		Data interface{} `json:"data"`
	} `json:"result"`
}

func NewSuccessResponse(response interface{}) Success {
	success := Success{Code: http.StatusOK}
	success.Result.Data = response
	return success
}

func EmptySuccessResponse() Success {
	success := Success{Code: http.StatusOK}
	success.Result.Data = ""
	return success
}
