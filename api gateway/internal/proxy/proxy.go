package proxy

import (
	"api-gateway/internal/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	droneServiceURL  string
	policeServiceURL string
	mapServiceURL    string
}

func NewProxyHandler() *ProxyHandler {
	cfg := config.New()
	return &ProxyHandler{
		droneServiceURL:  cfg.DroneServiceURL,
		policeServiceURL: cfg.PoliceServiceURL,
		mapServiceURL:    cfg.MapServiceURL,
	}
}

func (p *ProxyHandler) ProxyToDroneService(c *gin.Context) {
	p.proxyRequest(c, p.droneServiceURL)
}

func (p *ProxyHandler) ProxyToPoliceService(c *gin.Context) {
	p.proxyRequest(c, p.policeServiceURL)
}

func (p *ProxyHandler) ProxyToMapService(c *gin.Context) {
	p.proxyRequest(c, p.mapServiceURL)
}

func (p *ProxyHandler) proxyRequest(c *gin.Context, targetURL string) {
	userID, _ := c.Get("user_id")
	roleID, _ := c.Get("role_id")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}

	var requestBody map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &requestBody); err != nil {
			requestBody = make(map[string]interface{})
		}
	} else {
		requestBody = make(map[string]interface{})
	}

	requestBody["user_id"] = userID
	requestBody["user_role"] = roleID

	modifiedBody, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error modifying request body"})
		return
	}

	url := targetURL + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		url += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, url, bytes.NewBuffer(modifiedBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating proxy request"})
		return
	}

	for key, values := range c.Request.Header {
		if key != "Authorization" {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
