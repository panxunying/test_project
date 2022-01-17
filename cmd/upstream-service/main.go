package main

import (
	"github.com/gin-gonic/gin"

	"github.com/jamsman94/microservice-for-test/pkg"
)

func main() {
	r := gin.Default()
	r.GET("/upstreamAPI/1", UpstreamHandler1)
	r.Run(":8080")
}

func UpstreamHandler1(c *gin.Context) {
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
	if req.Expectation == "upstream_error" {
		c.JSON(400, gin.H{
			"code":    2,
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
