package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

type result struct {
	srcImagePath   string
	thumbnailImage *image.NRGBA
	err            error
}

// Image processing - sequential
// Input - directory with images.
// output - thumbnail images
func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to send directory path of images")
	}
	start := time.Now()

	err := setupPipeline(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Time taken: %s\n", time.Since(start))
}

func setupPipeline(root string) error {
	done := make(chan struct{})
	defer close(done)

	pathChan, errChan := walkFiles(done, root)
	resChan := processImage(done, pathChan)
	for res := range resChan {
		if res.err != nil {
			return res.err
		}
		saveThumbnail(res.srcImagePath, res.thumbnailImage)
	}
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

// walfiles - take diretory path as input
// does the file walk
// generates thumbnail images
// saves the image to thumbnail directory.
func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	pathChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(pathChan)
		defer close(errChan)
		errChan <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

			// filter out error
			if err != nil {
				return err
			}

			// check if it is file
			if !info.Mode().IsRegular() {
				return nil
			}

			// check if it is image/jpeg
			contentType, _ := getFileContentType(path)
			if contentType != "image/jpeg" {
				return nil
			}
			select {
			case pathChan <- path:
			case <-done:
				return fmt.Errorf("walk canceled")
			}

			return nil
		})
	}()
	return pathChan, errChan
}

// processImage - takes image file as input
// return pointer to thumbnail image in memory.
func processImage(done <-chan struct{}, pathChan <-chan string) <-chan *result {
	resChan := make(chan *result)
	var wg sync.WaitGroup
	const num = 5
	thumbnailer := func() {
		for path := range pathChan {
			srcImage, err := imaging.Open(path)
			if err != nil {
				select {
				case resChan <- &result{path, nil, err}:
				case <-done:
					return
				}
			}
			thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)
			res := &result{
				srcImagePath:   path,
				thumbnailImage: thumbnailImage,
				err:            err,
			}
			select {
			case resChan <- res:
			case <-done:
				return
			}
		}
	}
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			thumbnailer()
			wg.Done()
		}()
	}
	// load the image from file

	// scale the image to 100px * 100px
	go func() {
		wg.Wait()
		close(resChan)
	}()
	return resChan
}

// saveThumbnail - save the thumnail image to folder
func saveThumbnail(srcImagePath string, thumbnailImage *image.NRGBA) error {
	filename := filepath.Base(srcImagePath)
	dstImagePath := "thumbnail/" + filename

	// save the image in the thumbnail folder.
	err := imaging.Save(thumbnailImage, dstImagePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s -> %s\n", srcImagePath, dstImagePath)
	return nil
}

// getFileContentType - return content type and error status
func getFileContentType(file string) (string, error) {

	out, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
