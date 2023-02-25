package main

import (
	"converter-app-be/internal/convert"
	"converter-app-be/internal/download"
	"converter-app-be/internal/io"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(errorHandler)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I am alive",
		})
	})

	router.POST("/convert/jpgToPng", func(c *gin.Context) {
		file, header, _ := c.Request.FormFile("file")

		fileName, fileLength, _ := io.WriteFileToTempFolder(file)
		convert.JpgToPngConversion(fileName)

		c.Header("Content-Length", fileLength)
		c.Header("Content-Disposition", "attachment; filename="+fileName+".png")
		c.Header("Content-Type", "application/octet-stream")
		c.File("/tmp/" + fileName + ".png")
		c.JSON(http.StatusOK, gin.H{
			"message":  "File converted successfully",
			"filename": header.Filename,
			"filesize": header.Size,
		})
	})

	router.GET("/download/youtube", func(c *gin.Context) {
		header := c.Request.Header
		videoInfo := download.GetYoutubeVideoInfo(header.Get("url"))

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, videoInfo)
	})

	router.POST("/download/youtube", func(c *gin.Context) {
		header := c.Request.Header
		url := header.Get("url")
		formatId, _ := strconv.Atoi(header.Get("formatId"))

		fileName, _ := download.FromYoutube(url, formatId)

		//c.Header("Content-Length", fileLength)
		c.Header("Content-Disposition", "attachment; filename="+fileName+".mp4")
		c.Header("Content-Type", "application/octet-stream")
		c.File("/tmp/" + fileName + ".mp4")
		c.JSON(http.StatusOK, gin.H{
			"message": "File downloaded successfully",
		})
	})

	router.Run(":8080")
}

func errorHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("exception caught: ", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
	}()
	c.Next()
}
