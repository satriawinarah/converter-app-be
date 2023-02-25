package main

import (
	"converter-app-be/internal/convert"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
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

	router.POST("/convert/jpgToPng", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing file parameter",
			})
			return
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading file",
			})
			return
		}

		fileName := uuid.New().String()
		newFile, err := os.OpenFile("/tmp/"+fileName+".jpg", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error opening temp directory",
			})
			return
		}
		defer file.Close()

		_, err = newFile.Write(data)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error writing data to temp directory",
			})
			return
		}

		err = convert.JpgToPngConversion(fileName)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error converting file",
			})
			return
		}

		c.Header("Content-Disposition", "attachment; filename="+fileName+".png")
		c.Header("Content-Type", "application/octet-stream")
		c.File("/tmp/" + fileName + ".png")
		c.JSON(http.StatusOK, gin.H{
			"message":  "File converted successfully",
			"filename": header.Filename,
			"filesize": header.Size,
		})
	})

	router.Run(":8080")
}
