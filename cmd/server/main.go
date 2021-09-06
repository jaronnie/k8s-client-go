package main

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/api/v1/create-credential", CreateCredential)
	router.Run(":8080")
}

func CreateCredential(c *gin.Context) {
	var (
		url        string
		ca         string
		token      string
		kubeconfig multipart.File
		err        error
	)
	url = c.PostForm("url")
	ca = c.PostForm("ca")
	token = c.PostForm("token")
	if url != "" && ca != "" && token != "" {
		c.IndentedJSON(200, gin.H{
			"url":   url,
			"ca":    ca,
			"token": token,
		})
		return
	}
	if kubeconfig, _, err = c.Request.FormFile("file"); err != nil {
		c.IndentedJSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if url == "" && ca == "" && token == "" && kubeconfig != nil {
		var buf bytes.Buffer
		io.Copy(&buf, kubeconfig)
		c.IndentedJSON(200, gin.H{
			"content": buf.String(),
		})
	}
}
