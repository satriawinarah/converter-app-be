package io

import (
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strconv"
)

func WriteFileToTempFolder(file multipart.File) (string, string, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	fileName := uuid.New().String()
	newFile, err := os.OpenFile("/tmp/"+fileName+".jpg", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return "", strconv.Itoa(len(data)), err
	}
	defer file.Close()

	_, err = newFile.Write(data)
	if err != nil {
		log.Println(err)
		return "", strconv.Itoa(len(data)), err
	}

	return fileName, strconv.Itoa(len(data)), nil
}
