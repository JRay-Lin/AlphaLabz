package pocketbase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/patrickmn/go-cache"
)

type UserRole struct {
	UserId string `json:"id"`
	Name   string `json:"name"`
	RoleId string `json:"role"`
}

func (pbClient *PocketBaseClient) FetchUserRole(userJwt string) (UserRole, error) {
	roles, found := pbClient.UserCache.Get(userJwt)
	if found {
		fmt.Println("Cache hit!")
		return roles.(UserRole), nil
	} else {
		fmt.Println("uncache jwt")
		url := fmt.Sprintf("%s/api/role/%s", pbClient.BaseURL, userJwt)
		// Create request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return UserRole{}, err
		}
		req.Header.Set("Authorization", "Bearer "+pbClient.SuperToken)

		// Send request
		resp, err := pbClient.HTTPClient.Do(req)
		if err != nil {
			return UserRole{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return UserRole{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var userRoleResponse UserRole
		if err := json.NewDecoder(resp.Body).Decode(&userRoleResponse); err != nil {
			return UserRole{}, err
		}

		// Save to the Cache for future use (30 minutes expiry time). This will prevent us from hitting our API again until 24 hours have passed and we'll need to refresh it in that case.)
		pbClient.UserCache.Set(userJwt, userRoleResponse, cache.DefaultExpiration)

		return userRoleResponse, nil
	}
}
