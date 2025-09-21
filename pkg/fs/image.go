package fs

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/image/webp"
	"image"
	_ "image/gif" // подключаем форматы
	_ "image/jpeg"
	_ "image/png"
	"os"
	"regexp"
	"slices"
	"strings"
)

type IFileSystemImage interface {
	GetImageExtension(base64Image *string) (string, error)
	IsSupportingExtension(extension string) bool
	SafeImageToStorage(extension string, base64Image *string) (string, error)
	GetPath() string
}

var ImageExtensions = []string{"jpeg", "jpg", "png", "webp"}

type ImageConfig struct {
	Path       string `toml:"image_path"`
	extensions []string
	/*
		Добавить поддержку конфигурации расширений и размера изображения
	*/
}

func NewImageConfig() *ImageConfig {
	return &ImageConfig{
		extensions: ImageExtensions,
	}

}

type Image struct {
	config *ImageConfig
}

func NewImage(config *ImageConfig) *Image {
	return &Image{
		config: config,
	}
}
func (i *Image) GetPath() string {
	return i.config.Path
}

func (i *Image) IsSupportingExtension(extension string) bool {
	return slices.Contains(i.config.extensions, extension)
}

func (i *Image) GetImageExtension(base64Image *string) (string, error) {
	reg := regexp.MustCompile(`^data:image/([a-zA-Z0-9]+);base64,`)

	matches := reg.FindStringSubmatch(*base64Image)

	if len(matches) < 2 {
		return "", errors.New("Extension of image not found. Check provided image")
	}

	return matches[1], nil
}

func (i *Image) IsValidImage(r *strings.Reader) error {
	_, _, err := image.Decode(r)

	if err != nil {
		if _, errWebp := webp.Decode(r); errWebp != nil {
			return errors.New("Cant decode provided image. Please check data correctness")
		}
	}

	return nil
}

func (i *Image) SafeImageToStorage(extension string, base64Image *string) (string, error) {
	trimmed := strings.TrimPrefix(*base64Image, fmt.Sprintf("data:image/%s;base64,", extension))
	imageData, err := base64.StdEncoding.DecodeString(trimmed)

	if err != nil {
		return "", err
	}

	err = i.IsValidImage(strings.NewReader(string(imageData)))

	if err != nil {
		return "", err
	}

	filename := uuid.New().String() + "." + extension

	err = os.MkdirAll(i.config.Path, 0755)

	if err != nil {
		return "", err
	}

	err = os.WriteFile(i.config.Path+filename, imageData, 0644)

	if err != nil {
		return "", err
	}

	return filename, nil
}
