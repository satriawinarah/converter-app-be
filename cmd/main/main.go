package main

import (
	"converter-app-be/internal/convert"
	"converter-app-be/internal/io"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

		fileName, fileLength, err := io.WriteFileToTempFolder(file)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error while processing file",
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

	router.Run(":8080")
}
