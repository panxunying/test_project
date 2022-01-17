package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jamsman94/microservice-for-test/pkg"
)

func main() {
	r := gin.Default()
	r.GET("/withoutUpstreamSvc", DownstreamHandlerWithoutUpstream)
	r.GET("/withUpstreamAPI", DownstreamHandler1)
	r.Run(":8081")
}

func DownstreamHandlerWithoutUpstream(c *gin.Context) {
	req := new(pkg.TestRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    2,
			"message": "failed_to_decode_request",
			"data":    "error: " + err.Error(),
		})
		return
	}
	if req.Expectation == "downstream_error" {
		c.JSON(400, gin.H{
			"code":    3,
			"message": "giving_back_an_expected_error",
			"data":    "error: all according to the plan",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    "some_meaningful_data",
	})
}

func DownstreamHandler1(c *gin.Context) {
	req := new(pkg.TestRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    1,
			"message": "failed_to_decode_request",
			"data":    "error: " + err.Error(),
		})
		return
	}
	if req.Expectation == "downstream_error" {
		c.JSON(400, gin.H{
			"code":    2,
			"message": "giving_back_an_expected_error",
			"data":    "error: all according to the plan",
		})
		return
	}
	reqBody, _ := json.Marshal(req)
	outboundReq, err := http.NewRequest("GET", "http://localhost:8080/upstreamAPI/1", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(400, gin.H{
			"code":    3,
			"message": "failed_to_create_outbound_request",
			"data":    "error: " + err.Error(),
		})
		return
	}
	outboundReq.Header = c.Request.Header
	client := &http.Client{}
	resp, err := client.Do(outboundReq)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    4,
			"message": "failed_to_call_upstream_service",
			"data":    "error: " + err.Error(),
		})
		return
	}

	res := map[string]interface{}{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &res)
	if resp.StatusCode == http.StatusBadRequest {
		c.JSON(400, gin.H{
			"code":    4,
			"message": "upstream_service_gives_an_error",
			"data":    "error: " + res["message"].(string),
		})
		return
	}
	printBody := ""
	for key, val := range res {
		printBody = printBody + fmt.Sprintf("[%s : %s] ", key, val)
	}
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    "upstream response:" + printBody,
	})
}
