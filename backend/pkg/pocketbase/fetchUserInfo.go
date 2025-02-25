package pocketbase

import (
	"alphalabz/pkg/tools"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/patrickmn/go-cache"
)

type UserInfo struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	UserSettings string `json:"user_settings" `
	Avatar       string `json:"avatar,omitempty"`
	Gender       string `json:"gender,omitempty"`
	BirthDate    string `json:"birthdate,omitempty"`
	Expand       struct {
		UserSettings struct {
			Id       string `json:"id"`
			Language string `json:"language"`
			Theme    string `json:"theme"`
		} `json:"user_settings"`
		Role struct {
			Id          string                 `json:"id"`
			Name        string                 `json:"name"`
			Description string                 `json:"description"`
			Permission  map[string]interface{} `json:"permission"`
		} `json:"role"`
	} `json:"expand"`
}

func (pbClient *PocketBaseClient) FetchUserInfo(userJWT string) (UserInfo, error) {
	// Get userId from JWT token
	userId, err := tools.GetUserIdFromJWT(userJWT)
	if err != nil {
		return UserInfo{}, err
	}

	// Check if user info is already cached
	userInfo, found := pbClient.UserInfoCache.Get(userId)
	if found {
		return userInfo.(UserInfo), nil
	} else {
		// Fetch user info from the server and cache it
		url := fmt.Sprintf("%s/api/collections/users/records/%s?expand=role,user_settings", pbClient.BaseURL, userId)

		// Create HTTP request and send it to the server
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return UserInfo{}, err
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", pbClient.SuperToken))

		// Send HTTP request and receive response from the server
		resp, err := pbClient.HTTPClient.Do(req)
		if err != nil {
			return UserInfo{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var userResponse UserInfo
			if err = json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
				return UserInfo{}, err
			}

			// Cache the fetched user info for future use
			pbClient.UserInfoCache.Set(userId, userResponse, cache.DefaultExpiration)
			return userResponse, nil
		} else {
			return UserInfo{}, fmt.Errorf("failed to fetch user info: %s", resp.Status)
		}

	}
}
