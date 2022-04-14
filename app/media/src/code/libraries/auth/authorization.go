package auth

import (
	"errors"
	"github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/dtos"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func AuthorizeRequest(webRequest *http.Request) (*dtos.User, error) {
	authHeader := webRequest.Header.Get("Authorization")

	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return nil, errors.New("authorization token malformed")
	}

	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return nil, errors.New("authorization token malformed")
	}

	userToken, err := uuid.Parse(token)

	if err != nil {
		return nil, errors.New("uuid malformed")
	}

	users := dtos.Users{}
	user, err := users.GetByUuid(userToken)

	if err != nil {
		return nil, errors.New("error: " + err.Error())
	}

	if user.Status != "Active" {
		return nil, errors.New("user is not active")
	}

	return user, nil
}
