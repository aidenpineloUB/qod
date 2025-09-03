// Filename: internal/data/comments.go
package data

import (
  "time"
)

// each name begins with uppercase so that they are exportable/public
type Comment struct {
    ID int64                     
    Content  string             
    Author  string              
    CreatedAt  time.Time        
    Version int32                 
}                   
    
