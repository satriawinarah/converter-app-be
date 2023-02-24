package main

import (
	"converter-app-be/internal/convert"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World",
		})
	})

	router.POST("/convert", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing file parameter",
			})
			return
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading file",
			})
			return
		}

		dir, err := os.Open("temp/")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error opening temp directory",
			})
			return
		}
		defer file.Close()

		_, err = dir.Write(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error writing data to temp directory",
			})
			return
		}

		err = convert.JpgToPngConversion(header.Filename)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error converting file",
			})
			return
		}

		c.Header("Content-Disposition", "attachment; filename="+header.Filename)
		c.Header("Content-Type", "application/octet-stream")
		c.File("temp/" + header.Filename + ".png")
		c.JSON(http.StatusOK, gin.H{
			"message":  "File uploaded successfully",
			"filename": header.Filename,
			"filesize": header.Size,
		})
	})

	router.Run(":8080")
}
