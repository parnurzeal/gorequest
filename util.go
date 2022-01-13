package gorequest

import (
	"net/http"
	"unsafe"
)

func cloneMapArray(old map[string][]string) map[string][]string {
	newMap := make(map[string][]string, len(old))
	for k, vals := range old {
		newMap[k] = make([]string, len(vals))
		for i := range vals {
			newMap[k][i] = vals[i]
		}
	}
	return newMap
}

func shallowCopyData(old map[string]interface{}) map[string]interface{} {
	if old == nil {
		return nil
	}
	newData := make(map[string]interface{})
	for k, val := range old {
		newData[k] = val
	}
	return newData
}

func shallowCopyDataSlice(old []interface{}) []interface{} {
	if old == nil {
		return nil
	}
	newData := make([]interface{}, len(old))
	copy(newData, old)
	return newData
}

func shallowCopyFileArray(old []File) []File {
	if old == nil {
		return nil
	}
	newData := make([]File, len(old))
	copy(newData, old)
	return newData
}

func shallowCopyCookies(old []*http.Cookie) []*http.Cookie {
	if old == nil {
		return nil
	}
	newData := make([]*http.Cookie, len(old))
	copy(newData, old)
	return newData
}

func shallowCopyErrors(old []error) []error {
	if old == nil {
		return nil
	}
	newData := make([]error, len(old))
	copy(newData, old)
	return newData
}

func statusesContains(statuses []int, respStatus int) bool {
	for _, status := range statuses {
		if status == respStatus {
			return true
		}
	}
	return false
}

// ===========================================================

// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ===========================================================
// copy from gin/util.go of https://github.com/gin-gonic/gin
// MIT License
func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}
