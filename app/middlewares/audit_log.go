package middlewares

import (
	"apiok-admin/app/models"
	"apiok-admin/app/services"
	"apiok-admin/app/utils"
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" && c.Request.URL.Path == "/admin/log/list" {
			c.Next()
			return
		}

		username := ""
		token := c.GetHeader("auth-token")
		if token != "" {
			parsedUsername, err := utils.ParseToken(token)
			if err == nil {
				username = parsedUsername
			}
		}

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		responseBody := &bytes.Buffer{}
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}
		c.Writer = writer

		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		action := getAction(c.Request.Method, c.Request.URL.Path)
		resourceType := getResourceType(c.Request.URL.Path)
		resourceID := getResourceID(c.Request.URL.Path)

		logData := &models.Logs{
			Username:     username,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			IP:           c.ClientIP(),
			StatusCode:   c.Writer.Status(),
		}

		if len(requestBody) > 0 && len(requestBody) < 10000 {
			logData.RequestData = string(requestBody)
		}

		if responseBody.Len() > 0 && responseBody.Len() < 10000 {
			logData.ResponseData = responseBody.String()
		}

		if c.Writer.Status() >= 400 {
			var errorResp map[string]interface{}
			if err := json.Unmarshal(responseBody.Bytes(), &errorResp); err == nil {
				if msg, ok := errorResp["msg"].(string); ok {
					logData.ErrorMessage = msg
				}
			}
		}

		go func() {
			_ = services.LogAdd(logData)
		}()

		_ = duration
	}
}

func getAction(method, path string) string {
	if method == "GET" {
		if strings.Contains(path, "/list") {
			return "list"
		}
		return "info"
	}
	if method == "POST" {
		return "create"
	}
	if method == "PUT" {
		return "update"
	}
	if method == "DELETE" {
		return "delete"
	}
	return "unknown"
}

func getResourceType(path string) string {
	if strings.Contains(path, "/service") {
		return "service"
	}
	if strings.Contains(path, "/router") {
		return "router"
	}
	if strings.Contains(path, "/upstream") {
		return "upstream"
	}
	if strings.Contains(path, "/user") {
		return "user"
	}
	if strings.Contains(path, "/certificate") {
		return "certificate"
	}
	if strings.Contains(path, "/plugin") {
		return "plugin"
	}
	if strings.Contains(path, "/cluster-node") {
		return "cluster_node"
	}
	if strings.Contains(path, "/letsencrypt") {
		return "letsencrypt"
	}
	return "unknown"
}

func getResourceID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if i > 0 && (part == "update" || part == "delete" || part == "info" || part == "switch" || part == "enable" || part == "release") {
			if i+1 < len(parts) {
				return parts[i+1]
			}
		}
		if part == "res_id" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

