package controllers

import (
	"fmt"
	"github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/dtos"
	"github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/auth"
	"github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/helper"
	"github.com/google/uuid"
	"gopkg.in/gographics/imagick.v2/imagick"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Student struct {

	// defining struct fields
	Name  string
	Marks int
	Id    string
}

const MaxUploadSize = 1024 * 1024 * 10 // 1MB
const MaxImageDimension = 1024 * 2     // 2048
const ImageThumbnailSize = 150

func UploadImageControllerHandle(responseWriter http.ResponseWriter, webRequest *http.Request) {

	user, err := auth.AuthorizeRequest(webRequest)

	if err != nil {
		http.Error(responseWriter, "You are not authorized to make this request", http.StatusUnauthorized)
		return
	}

	webRequest.Body = http.MaxBytesReader(responseWriter, webRequest.Body, MaxUploadSize)
	if err := webRequest.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(responseWriter, "The uploaded file is too big. Please choose an file that's less than 10MB in size", http.StatusBadRequest)
		return
	}

	userId := webRequest.FormValue("user_id")
	entityId := webRequest.FormValue("entity_id")
	entityName := webRequest.FormValue("entity_name")
	imageClass := webRequest.FormValue("image_class")
	parentId := webRequest.FormValue("parent_id")

	if userId == "" || entityId == "" || entityName == "" || imageClass == "" {
		http.Error(responseWriter, "There was an error parsing the request: userId="+userId+"; entityId="+entityId+"; entityName="+entityName+"; imageClass="+imageClass+";", http.StatusBadRequest)
		return
	}

	targetDirectory := "/media/images/" + entityName + "s/" + entityId
	targetThumbDirectory := "/media/images/" + entityName + "s/" + entityId + "/thumb"

	file, fileHeader, err := webRequest.FormFile("file")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	err = os.MkdirAll("/media/images", os.ModePerm)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll(targetDirectory, os.ModePerm)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll(targetThumbDirectory, os.ModePerm)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	uuid, err := uuid.NewUUID()

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	fileName := fmt.Sprintf(targetDirectory+"/%s%s", uuid.String(), filepath.Ext(fileHeader.Filename))
	fileThumbnail := fmt.Sprintf(targetThumbDirectory+"/%s%s", "thumb_"+uuid.String(), filepath.Ext(fileHeader.Filename))

	dst, err := os.Create(fileName)

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	localFile, err := os.Open(dst.Name())
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	defer localFile.Close()

	contentType, err := getFileContentType(localFile)

	if err != nil {
		http.Error(responseWriter, err.Error()+": "+contentType, http.StatusInternalServerError)
		return
	}

	if fileIsNotAnImage(contentType) {
		http.Error(responseWriter, "fileIsNotAnImage: "+contentType, http.StatusInternalServerError)
		os.Remove(dst.Name())
		return
	}

	err = resizeFileIfLargerThanMax(dst.Name())

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	imageCndName := buildCndNameFromFileName(fileName)
	imageCndThumbnail := buildCndNameFromFileName(fileThumbnail)

	createImageThumbnail(dst.Name(), fileThumbnail)
	image, err := insertImageRecord(user, fileName, imageCndName, imageCndThumbnail, entityId, entityName, contentType, imageClass, parentId)

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	type ReturnImage struct {
		Image     string `json:"image"`
		Thumb     string `json:"thumb"`
		ImageId   int    `json:"image_id"`
		Type      string `json:"type"`
		CompanyId string `json:"company_id"`
	}

	imageData := ReturnImage{Image: imageCndName, Thumb: imageCndThumbnail, ImageId: image.ImageId, Type: contentType, CompanyId: fmt.Sprint(user.CompanyId)}

	healthCheck := helper.TransactionResult{Success: true, Message: "image registered", Data: imageData}
	helper.JsonReturn(healthCheck, responseWriter)
}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func fileIsNotAnImage(fileType string) bool {
	switch fileType {
	case "image/jpeg", "image/png", "image/gif":
		return false
	default:
		return true
	}
}

func createImageThumbnail(localName string, fileThumbnail string) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	err := mw.ReadImage(localName)

	if err != nil {
		return err
	}

	originalWidth := mw.GetImageWidth()
	originalHeight := mw.GetImageHeight()

	err = resizeImage(mw, originalWidth, originalHeight, ImageThumbnailSize, false)

	if err = mw.WriteImage(fileThumbnail); err != nil {
		return err
	}

	mwthumb := imagick.NewMagickWand()
	err = mwthumb.ReadImage(fileThumbnail)

	if err != nil {
		return err
	}

	thumbWidth := mwthumb.GetImageWidth()
	thumbHeight := mwthumb.GetImageHeight()

	xcrop := processCropValue(thumbWidth, ImageThumbnailSize)
	ycrop := processCropValue(thumbHeight, ImageThumbnailSize)

	mwthumb.CropImage(ImageThumbnailSize, ImageThumbnailSize, xcrop, ycrop)

	if err = mwthumb.WriteImage(fileThumbnail); err != nil {
		return err
	}

	return nil
}

func processCropValue(originalSize uint, cropSize uint) int {
	if originalSize-cropSize > 0 {
		return int((originalSize - cropSize) / 2)
	}

	return 0
}

func resizeFileIfLargerThanMax(localName string) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	err := mw.ReadImage(localName)

	if err != nil {
		return err
	}

	originalWidth := mw.GetImageWidth()
	originalHeight := mw.GetImageHeight()

	if originalWidth > MaxImageDimension || originalHeight > MaxImageDimension {
		err := resizeImage(mw, originalWidth, originalHeight, MaxImageDimension, true)

		os.Remove(localName)

		if err = mw.WriteImage(localName); err != nil {
			return err
		}
	}

	return nil
}

func resizeImage(mw *imagick.MagickWand, originalWidth uint, originalHeight uint, imageDimension uint, contain bool) error {

	var err error

	if contain {
		if (originalWidth > imageDimension && originalHeight > imageDimension && originalWidth > originalHeight) || (originalWidth > imageDimension && originalWidth > originalHeight) {
			err = mw.ResizeImage(imageDimension, calculateResizedDimension(originalWidth, originalHeight, imageDimension), imagick.FILTER_LANCZOS, 1)
		} else if (originalWidth > imageDimension && originalHeight > imageDimension && originalWidth < originalHeight) || (originalHeight > imageDimension && originalWidth < originalHeight) {
			err = mw.ResizeImage(calculateResizedDimension(originalHeight, originalWidth, imageDimension), imageDimension, imagick.FILTER_LANCZOS, 1)
		}
	} else {
		if (originalWidth > imageDimension && originalHeight > imageDimension && originalWidth > originalHeight) || (originalWidth > imageDimension && originalWidth > originalHeight) {
			err = mw.ResizeImage(calculateResizedDimension(originalHeight, originalWidth, imageDimension), imageDimension, imagick.FILTER_LANCZOS, 1)
		} else if (originalWidth > imageDimension && originalHeight > imageDimension && originalWidth < originalHeight) || (originalHeight > imageDimension && originalWidth < originalHeight) {
			err = mw.ResizeImage(imageDimension, calculateResizedDimension(originalWidth, originalHeight, imageDimension), imagick.FILTER_LANCZOS, 1)
		}
	}

	if err != nil {
		return err
	}

	// Set the compression quality to 95 (high quality = low compression)
	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		return err
	}

	return nil
}

func calculateResizedDimension(sideA uint, sideB uint, imageDimension uint) uint {
	return uint((float32(imageDimension) / float32(sideA)) * float32(sideB))
}

func buildCndNameFromFileName(name string) string {
	return strings.Replace(name, "/media", "/cdn", -1)
}

func insertImageRecord(user *dtos.User, fileName string, cndName string, cndThumbnail string, entityId string, entityName string, contentType string, imageClass string, parentId string) (*dtos.Image, error) {

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	err := mw.ReadImage(fileName)

	if err != nil {
		return nil, err
	}

	originalWidth := mw.GetImageWidth()
	originalHeight := mw.GetImageHeight()

	images := dtos.Images{}

	newEntityId, _ := strconv.Atoi(entityId)
	companyId := strconv.Itoa(user.CompanyId)
	newParentIntId, _ := strconv.Atoi(parentId)
	newParentId := helper.NullInt{Value: newParentIntId, Valid: true}

	imageModel := dtos.Image{}
	imageModel.UserId = user.UserId
	imageModel.ParentId = newParentId
	imageModel.CompanyId = companyId
	imageModel.EntityId = newEntityId
	imageModel.EntityName = entityName
	imageModel.ImageClass = imageClass
	imageModel.Type = contentType
	imageModel.Title = fileName
	imageModel.Url = cndName
	imageModel.Thumb = cndThumbnail
	imageModel.Width = int(originalWidth)
	imageModel.Height = int(originalHeight)
	imageModel.CreatedBy = user.UserId
	imageModel.CreatedOn = helper.NullTime{Value: time.Now(), Valid: true}
	imageModel.UpdatedBy = user.UserId
	imageModel.LastUpdated = helper.NullTime{Value: time.Now(), Valid: true}

	return images.CreateNew(imageModel)
}
