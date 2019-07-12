package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.LoadHTMLGlob("public/html/**")

	staticBox := packr.New("staticBox", "./public/static")
	router.StaticFS("/static", http.FileSystem(staticBox))


	imgDir := "printKit-image"
	uploadBox := packr.New("uploadBox", imgDir)
	router.StaticFS("/image", http.FileSystem(uploadBox))

	// 创建目录
	exist, err := PathExists(imgDir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}
	if !exist {
		// 创建文件夹
		err := os.Mkdir(imgDir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		}
	}
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Posts",
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		response := make(map[string]interface{})

		// Source
		file, err := c.FormFile("image")
		if err != nil {
			response["code"] = -1
			response["msg"] = "param error"
			c.JSON(http.StatusOK, response)
			return
		}

		current := time.Now().Nanosecond()

		filename := imgDir + "/" + strconv.Itoa(current) + "-" + filepath.Base(file.Filename)
		url := "/image/" + strconv.Itoa(current) + "-" + filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			response["code"] = -2
			response["msg"] = "upload file err"
			c.JSON(http.StatusOK, response)
			return
		}

		response["code"] = 0
		response["msg"] = "ok"
		response["url"] = url
		c.JSON(http.StatusOK, response)
	})

	router.POST("/detail", func(c *gin.Context) {
		name := c.PostForm("name")
		image1 := c.PostForm("image1")
		image2 := c.PostForm("image2")
		image3 := c.PostForm("image3")
		image4 := c.PostForm("image4")
		image5 := c.PostForm("image5")

		t := time.Now()
		dayStr := fmt.Sprintf("%d.%d.%d", t.Year(), t.Month(), t.Day())

		if name == "" && image1 == "" && image2 == "" && image3 == "" && image4 == "" && image5 == "" {
			c.String(http.StatusOK, "Params error.")
			return
		}

		c.HTML(http.StatusOK, "detail.tmpl", map[string]interface{}{
			"Name":   name,
			"Date":   dayStr,
			"image1": image1,
			"image2": image2,
			"image3": image3,
			"image4": image4,
			"image5": image5,
		})
	})

	router.Run(":8080")
}


func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
