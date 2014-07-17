package user

import (
  // External Packages
  "net/http"
  "time"
  "github.com/martini-contrib/binding"
)

type User struct {
  Id             int64     `json:"id"`
  Email          string    `json:"email" binding:"required" sql:"size:255;not null;unique"`
  AccessToken    string    `json:"-" binding:"required" sql:"size:255;not null"`
  RefreshToken   string    `json:"-" binding:"required" sql:"size:255;not null"`
  ExpireTokenAt  time.Time `json:"-"`
  CreatedAt      time.Time
  UpdatedAt      time.Time
  DeletedAt      time.Time `json:"-"`
}

// This method implements binding.Validator and is executed by the binding.Validate middleware
func (user User) Validate(errors binding.Errors, req *http.Request) binding.Errors {
    if len(user.Email) < 3 {
        errors = append(errors, binding.Error{
            FieldNames:     []string{"email"},
            Message:        "Too short; minimum 3 characters",
        })
    } else if len(user.Email) > 120 {
      errors = append(errors, binding.Error{
          FieldNames:     []string{"email"},
          Message:        "Too long; maximum 255 characters",
      })
    }
    return errors
}
