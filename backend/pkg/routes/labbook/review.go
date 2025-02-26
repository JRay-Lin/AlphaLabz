package labbook

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"alphalabz/pkg/tools"
	"encoding/json"
	"net/http"
)

type labbookReviewRequest struct {
	LabbookId string `json:"labbook_id"`
	Status    string `json:"status"`
	Comment   string `json:"comment"`
}

func HandleLabBookReview(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient, ce *casbin.CasbinEnforcer) {
	// Check if the request method is PATCH
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get request authorization header
	rawToken, err := tools.TokenExtractor(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Verify user permission using Casbin enforcer
	hasPermission, err := ce.VerifyJWTPermission(pbClient, rawToken, casbin.PermissionConfig{
		Resources: "lab_books",
		Actions:   "update",
		Scopes:    "status"})
	if err != nil {
		http.Error(w, "Failed to verify permission", http.StatusInternalServerError)
		return
	}

	// Check if the user has permission to update lab books
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Decode the JSON request body into a LabBookReviewRequest struct
	var reviewRequest labbookReviewRequest
	err = json.NewDecoder(r.Body).Decode(&reviewRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the labbook ID and status fields in the request
	if reviewRequest.LabbookId == "" || (reviewRequest.Status != "approved" && reviewRequest.Status != "Rejected") {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Retrieve lab book information from PocketBase database
	labbook, err := pbClient.ViewLabbook(reviewRequest.LabbookId, []string{"id", "reviewer", "status"})
	if err != nil {
		http.Error(w, "Failed to retrieve lab book", http.StatusInternalServerError)
		return
	}

	if labbook.ReviewStatus != "pending" {
		http.Error(w, "The lab book aleardy been verified", http.StatusConflict)
		return
	}

	// Extract userId from JWT token
	userId, err := tools.GetUserIdFromJWT(rawToken)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	// Verify that the reviewer is the same as the one in the request body
	if labbook.Reviewer != userId {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Update the lab book review status and comments
	if err = pbClient.UpdateLabbook(reviewRequest.LabbookId, map[string]interface{}{
		"review_status":  reviewRequest.Status,
		"review_comment": reviewRequest.Comment,
	}); err != nil {
		http.Error(w, "Failed to update lab book", http.StatusInternalServerError)
		return
	}

	// Encode response with success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Lab book review updated successfully"})
}
