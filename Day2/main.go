package main

import (
	"camp/Day2/controller"
	"camp/Day2/dao"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	if err := Init("./data/"); err != nil {
		fmt.Println("main init failed", err)
		os.Exit(1)
	}
	r := gin.Default()

	r.GET("/community/page/get/:id", func(c *gin.Context) {
		topicId := c.Param("id")
		data := controller.QueryPageInfo(topicId)
		c.JSON(200, data)
	})
	r.Run(":8080")
}

func Init(filePath string) error {
	if err := dao.Init(filePath); err != nil {
		fmt.Println("dao init failed", err)
		return err
	}
	return nil
}
