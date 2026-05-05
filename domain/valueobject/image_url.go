package valueobject

import (
	"strings"

	"github.com/Josey34/goshop/domain/errors"
)

type ImageURL struct {
	url string
}

func NewImageURL(url string) (ImageURL, error) {
	if url == "" || (!strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://")) {
		return ImageURL{}, errors.NewValidation("image_url", map[string]string{"url": "must be a valid http or https url"})
	}

	return ImageURL{
		url: url,
	}, nil
}

func (i ImageURL) Value() string {
	return i.url
}
