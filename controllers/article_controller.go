package controllers

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown" // Markdown 渲染库
)

func ArticlePage(c *gin.Context) {
	// 获取文章文件名
	articleName := c.Param("name")
	filePath := "ThinkTankCentral/shared_folder/" + articleName

	// 检查文件是否存在
	if filepath.Ext(filePath) != ".md" {
		c.String(http.StatusBadRequest, "Invalid file type")
		return
	}

	// 读取文章内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.String(http.StatusNotFound, "Article not found")
		return
	}

	// 渲染 Markdown 为 HTML
	htmlContent := markdown.ToHTML(content, nil, nil)

	// 渲染文章详情模板
	c.HTML(http.StatusOK, "templates/article.html", gin.H{
		"content": string(htmlContent),
	})
}
