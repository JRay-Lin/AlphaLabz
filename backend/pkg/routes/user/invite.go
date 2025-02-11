package user

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"alphalabz/pkg/settings"
	"alphalabz/pkg/smtp"
	"alphalabz/pkg/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type newUserReq struct {
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
}

type inviteeData struct {
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
	RoleId   string `json:"role_id"`
}

type inviteResp struct {
	Link string `json:"link"`
}

func HandleInviteNewUser(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient, ce *casbin.CasbinEnforcer, sc *smtp.SMTPClient) {
	var permissionConfig = casbin.PermissionConfig{
		Resources: "users",
		Actions:   "create",
	}

	// Constrain request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse requester JWT
	rawJwtToken, err := tools.TokenExtractor(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var inviteData newUserReq
	if err := json.NewDecoder(r.Body).Decode(&inviteData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Obtain uesr role
	userRole, err := pbClient.FetchUserRole(rawJwtToken)
	if err != nil {
		http.Error(w, "Failed to fetch user role", http.StatusInternalServerError)
		return
	}

	// Forbidden any creation of admin account
	if strings.ToLower(inviteData.RoleName) == "admin" {
		http.Error(w, "You can't create an admin account", http.StatusForbidden)
		return
	}

	// Check request role exist
	roles, err := pbClient.GetAvailableRoles([]string{"id", "name"})
	if err != nil {
		http.Error(w, "Failed to fetch available roles", http.StatusInternalServerError)
		return
	}

	var invitee inviteeData
	roleExists := false
	for _, role := range roles {
		if strings.EqualFold(role.Name, inviteData.RoleName) {
			roleExists = true
			invitee = inviteeData{
				Email:    inviteData.Email,
				RoleName: role.Name,
				RoleId:   role.Id,
			}
			break
		}
	}

	if !roleExists {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	notAdminPermission, err := ce.VerifyPermission(userRole.RoleId, permissionConfig.Resources, permissionConfig.Actions, "notAdmin")
	if err != nil {
		http.Error(w, "Failed to verify permission", http.StatusInternalServerError)
		return
	}

	if notAdminPermission {
		sendInviteResponse(w, invitee)
		return
	}

	rolePermission, err := ce.VerifyPermission(userRole.RoleId, permissionConfig.Resources, permissionConfig.Actions, invitee.RoleName)
	if err != nil {
		http.Error(w, "Failed to verify permission", http.StatusInternalServerError)
		return
	}

	if rolePermission {
		sendInviteResponse(w, invitee)
	} else {
		http.Error(w, "You don't have permission to create this role", http.StatusForbidden)
	}
}

func sendInviteResponse(w http.ResponseWriter, invitee inviteeData) {
	inviteLink, err := generateInvitation(invitee)
	if err != nil {
		http.Error(w, "Failed to generate invitation link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(inviteResp{inviteLink}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func generateInvitation(invitee inviteeData) (inviteLink string, err error) {
	// Get JWT secret from settings
	settings, err := settings.LoadSettings("settings.yml")
	if err != nil {
		return "", fmt.Errorf("failed to load settings")
	} else {
		// Get invite JWT secret from settings
		inviteSecret := settings.JWTSecret
		if inviteSecret == "" {
			return "", fmt.Errorf("invite secret not set")
		}

		// Generate JWT token
		secretKey := []byte(inviteSecret)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"email":     invitee.Email,
				"role_id":   invitee.RoleId,
				"role_name": invitee.RoleName,
				"exp":       time.Now().Add(time.Hour * 24).Unix(),
			})

		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			return "", err
		}

		// Format the invite link
		inviteLink = fmt.Sprintf("%s/invite?token=%s", settings.AppUrl, tokenString)
		return inviteLink, nil
	}
}
