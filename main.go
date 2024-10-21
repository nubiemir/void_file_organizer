package main

import (
	"bufio"
	"errors"
	"fmt"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"

	figlet "github.com/mbndr/figlet4go"
)

var FileNotFound = errors.New("File path not found")

func scann() (result string, err error) {
	reader := bufio.NewReader(os.Stdin)
	temp, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	result = strings.TrimSuffix(temp, "\n")
	return
}

func checkFilePath(filePath string) (err error) {
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return FileNotFound
	}
	return nil
}

func organizeFiles(src, dst string) (err error) {
	files, err := os.ReadDir(src)
	for _, file := range files {
		if !file.IsDir() {
			fp := path.Join(src, file.Name())
			mimeType := mime.TypeByExtension(filepath.Ext(fp))
			fileType := strings.Split(mimeType, "/")[0]
			assignFileType(dst, fileType, file.Name(), fp)
		}
	}
	return
}

func assignFileType(dst, fileType, fileName, src string) {
	switch strings.ToLower(fileType) {
	case "text", "audio", "image", "video":
		createDir(dst, fileType)
		moveFile(dst, fileType, fileName, src)
	default:
		createDir(dst, "other")
		moveFile(dst, "other", fileName, src)
	}
}

func createDir(dst, fileType string) {
	fp := path.Join(dst, fileType)
	os.MkdirAll(fp, 0777)
}

func moveFile(dst, fileType, fileName, src string) {
	fp := path.Join(dst, fileType, fileName)
	os.Rename(src, fp)
}

func main() {
	ascii := figlet.NewAsciiRender()
	options := figlet.NewRenderOptions()
	options.FontColor = []figlet.Color{
		figlet.ColorGreen,
		figlet.ColorCyan,
		figlet.ColorRed,
		figlet.ColorBlue,
	}
	options.FontName = "avatar"

	renderStr, _ := ascii.RenderOpts("Void File Organizer", options)
	fmt.Println(renderStr)

askForSrc:
	fmt.Println("Please enter the source folder: (./)")
	src, _ := scann()
	srcErr := checkFilePath(src)
	if srcErr != nil {
		fmt.Printf(srcErr.Error() + "\n")
		goto askForSrc
	}

askForDst:
	fmt.Println("Please enter the destination folder: (./)")
	dst, _ := scann()
	dstErr := checkFilePath(dst)
	if dstErr != nil {
		fmt.Printf(dstErr.Error() + "\n")
		goto askForDst
	}

	organizeFiles(src, dst)
}
