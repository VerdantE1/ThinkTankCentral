package controllers

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type Article struct {
	Title string // 文章标题（文件名）
	Path  string // 文件路径
}

// 主页：读取共享文件夹中的所有文章文件
func HomePage(c *gin.Context) {
	// 定义共享文件夹路径
	sharedFolder := "ThinkTankCentral/shared_folder"

	// 获取文件夹中的所有 Markdown 文件
	files, err := ioutil.ReadDir(sharedFolder)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading shared folder")
		return
	}

	// 遍历文件，生成文章列表
	var articles []Article
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" { // 只读取 Markdown 文件
			articles = append(articles, Article{
				Title: file.Name(),                // 使用文件名作为标题
				Path:  "/articles/" + file.Name(), // 跳转到文章详情页的路径
			})
		}
	}

	// 渲染主页模板
	c.HTML(http.StatusOK, "index.html", gin.H{
		"articles": articles,
	})
}
