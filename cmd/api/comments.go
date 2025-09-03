// Filename: cmd/api/comments.go
package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/aidenpineloUB/qod/internal/data"
)

func (appInstance *applicationDependencies) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold a comment
	// we use struct tags[``] to make the names display in lowercase
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
	comment := &data.Comment{
		Content:   incomingData.Content,
		Author:    incomingData.Author,
		CreatedAt: time.Now(),
		Version:   1,
	}

	// for now display the result
	fmt.Fprintf(w, "Comment created: %+v\n", comment)
}