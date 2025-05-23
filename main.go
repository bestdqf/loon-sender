package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Unzip 解压指定文件到目标目录

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
func main() {
	// TODO: Implement me

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/send-file", func(c *gin.Context) {

		fileHeader, err := c.FormFile("file")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid file"})
			return
		}
		prd := c.PostForm("prd")
		fmt.Printf("prd: %s", prd)

		// 保存上传文件到临时目录
		tempFile := filepath.Join("./", fileHeader.Filename)
		if err := c.SaveUploadedFile(fileHeader, tempFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		fmt.Println("dir: ", os.TempDir())

		// // // 解压文件到指定目录
		// destDir := "/data/wwwroot/lang-robot.bj.gooki.com"
		// if err := os.MkdirAll(destDir, 0755); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create destination directory"})
		// 	return
		// }

		// files, err := Unzip(tempFile, destDir)
		// if err != nil {
		// 	// log.Fatal(err)
		// }

		// fmt.Printf("%v", files)

		// if _,err := Unzip(tempFile, destDir); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unzip file: %v", err)})
		// 	return
		// }

		// 返回成功响应
		c.JSON(200, gin.H{
			"message": "File unzipped successfully",
			// "dest":    destDir,
		})

	})

	fmt.Println("HTTP 服务器已启动, 监听端口 32002")
	router.Run(":32002")
}
