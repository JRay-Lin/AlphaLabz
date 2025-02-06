package casbin

import (
	"alphalabz/pkg/pocketbase"
	"encoding/json"
	"fmt"
	"net/http"
)

// FetchPermissions fetch the latest permissons settings from the database and save as a local csv
func FetchPermissions(pbClient *pocketbase.PocketBaseClient) []RolePermission {
	url := fmt.Sprintf("%s/api/collections/roles/records?perPage=99", pbClient.BaseURL)

	// Construct request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	req.Header.Set("Authorization", "Bearer"+pbClient.SuperToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request successful")
	} else {
		fmt.Println("Request failed with status code:", resp.StatusCode)
	}

	var respData PermissionRespond
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return nil
	}

	// fmt.Println(respData.Items)

	return respData.Items
}
