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

// for now display the result
   fmt.Fprintf(w, "+%v\n", incomingData)



	

	// Create a Comment struct using the data package
// Copy the values from incomingData to a new Comment struct
// At this point in our code the JSON is well-formed JSON so now
// we will validate it using the Validator which expects a Comment
comment := &data.Comment {
    Content: incomingData.Content,
    Author: incomingData.Author,
}
// Initialize a Validator instance
 v := validator.New()
// Do the validation
data.ValidateComment(v, comment)
if !v.IsEmpty() {
    appInstance.failedValidationResponse(w, r, v.Errors)  // implemented later
    return
}

fmt.Fprintf(w, "%+v\n", incomingData)
	// for now display the result
	fmt.Fprintf(w, "Comment created: %+v\n", comment)
}