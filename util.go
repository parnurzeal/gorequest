package gorequest

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
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

// just need to change the array pointer?
func copyRetryable(old superAgentRetryable) superAgentRetryable {
	newRetryable := old
	newRetryable.RetryableStatus = make([]int, len(old.RetryableStatus))
	for i := range old.RetryableStatus {
		newRetryable.RetryableStatus[i] = old.RetryableStatus[i]
	}
	return newRetryable
}

func copyStats(old Stats) Stats {
	newStats := old
	return newStats
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

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name and file name.
func CreateFormFile(w *multipart.Writer, fieldname, filename string, contenttype string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contenttype)
	return w.CreatePart(h)
}
