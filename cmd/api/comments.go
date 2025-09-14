// Filename: cmd/api/comments.go
package main

import (
	"fmt"
	"errors"
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

func (a *applicationDependencies) displayCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL /v1/comments/:id so that we
	// can use it to query the comments table. We will 
	// implement the readIDParam() function later
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	// Call Get() to retrieve the comment with the specified id
	comment, err := a.commentModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	// display the comment
	envelope := map[string]interface{}{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Call Get() to retrieve the comment with the specified id
	comment, err := a.commentModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Use our temporary incomingData struct to hold the data
	// Note: I have changed the types to pointer to differentiate
	// between the client leaving a field empty intentionally
	// and the field not needing to be updated
	var incomingData struct {
		Content *string `json:"content"`
		Author  *string `json:"author"`
	}

	// perform the decoding
	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	// We need to now check the fields to see which ones need updating
	// if incomingData.Content is nil, no update was provided
	if incomingData.Content != nil {
		comment.Content = *incomingData.Content
	}
	// if incomingData.Author is nil, no update was provided
	if incomingData.Author != nil {
		comment.Author = *incomingData.Author
	}

	// Before we write the updates to the DB let's validate
	v := validator.New()
	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	// perform the update
	err = a.commentModel.Update(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	envelope := map[string]interface{}{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.commentModel.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// display the comment
	envelope := map[string]interface{}{
		"message": "comment successfully deleted",
	}
	err = a.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Create a struct to hold the query parameters
	// Later on we will add fields for pagination and sorting (filters)
	var queryParametersData struct {
		Content string
		Author  string
	}
	// get the query parameters from the URL
	queryParameters := r.URL.Query()

	// Load the query parameters into our struct
	queryParametersData.Content = a.getSingleQueryParameter(
		queryParameters,
		"content",
		"")

	queryParametersData.Author = a.getSingleQueryParameter(
		queryParameters,
		"author",
		"")

	comments, err := a.commentModel.GetAllFiltered(
		queryParametersData.Content,
		queryParametersData.Author)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	envelope := map[string]interface{}{
		"comments": comments,
	}
	err = a.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}