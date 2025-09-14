// Filename: cmd/api/comments.go
package main

import (
	"fmt"
	"net/http"
	"github.com/aidenpineloUB/qod/internal/data"
	"github.com/aidenpineloUB/qod/internal/validator"
)

func (appInstance *applicationDependencies) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold a comment
	// we use struct tags to make the names display in lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	// perform the decoding
	err := appInstance.readJSON(w, r, &incomingData)
	if err != nil {
		appInstance.badRequestResponse(w, r, err)
		return
	}

	// Create a Comment struct using the data package
	// Copy the values from incomingData to a new Comment struct
	comment := &data.Comment{
		Content: incomingData.Content,
		Author:  incomingData.Author,
	}

	// Initialize a Validator instance and validate BEFORE processing
	v := validator.New()
	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		appInstance.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Add the comment to the database table
	err = appInstance.commentModel.Insert(comment)
	if err != nil {
		appInstance.serverErrorResponse(w, r, err)
		return
	}

	// Set a Location header. The path to the newly created comment
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/comments/%d", comment.ID))

	// Send a JSON response with 201 (new resource created) status code
	envelope := map[string]interface{}{
		"comment": comment,
	}
	err = appInstance.writeJSON(w, http.StatusCreated, envelope, headers)
	if err != nil {
		appInstance.serverErrorResponse(w, r, err)
		return
	}
}