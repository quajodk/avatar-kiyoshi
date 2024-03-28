package media

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"wp-media-core/cors"
	gcstorage "wp-media-core/gc-storage"

	"github.com/gabriel-vasile/mimetype"
)

const mediaBasePath = "media"

func SetupRoutes(apiBasePath string) {

	handleMedia := http.HandlerFunc(mediaHandler)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, mediaBasePath), cors.Middleware(handleMedia))
}

func generateMediaName(originalFileName string) string {
	now := time.Now().UnixNano() // Get timestamp in nanoseconds
	originalName := strings.ReplaceAll(originalFileName, "-", "_")
	originalName = strings.ReplaceAll(originalName, " ", "_")
	return fmt.Sprintf("%d_%s", now, strings.ToLower(originalName))
}

func compressMediaWithFfmpeg(input multipart.File, w *io.PipeWriter) error {
	cmd := exec.Command("ffmpeg", "-i", "-", "-vcodec", "libx265", "-crf", "28", "pipe:1")
	cmd.Stdin = input
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func mediaHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		j, err := json.Marshal("ok")
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "application/json")
		res.Write(j)
	case http.MethodPost:
		req.ParseMultipartForm(50 << 20) // 500mb
		_, handler, err := req.FormFile("media")
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Fatal(err)
			return
		}
		file, err := handler.Open()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Fatal(err)
			return
		}
		defer file.Close()
		r, w := io.Pipe()
		mtype, err := mimetype.DetectReader(r)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Fatal(err)
			return
		}
		contentType := mtype.String()
		fileName := generateMediaName(handler.Filename)
		compressFile := filepath.Join(mediaDirectory, fmt.Sprintf("%s_%s", "comp", strings.ToLower(fileName)))
		go func() {
			defer w.Close()
			err = compressMediaWithFfmpeg(file, w)
		}()

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		link, err := gcstorage.Uploader.UploadFile(compressFile, contentType, r)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		// data, err := saveMedia(link)
		j, err := json.Marshal(link)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Header().Set("Content-Type", "application/json")
		res.Write(j)
	case http.MethodOptions:
		return
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
