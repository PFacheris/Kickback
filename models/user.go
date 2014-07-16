package user

import (
  // External Packages
  "database/sql"
  "time"
  "github.com/codegangsta/martini-contrib/binding"

  // Application Specific Imports
  . "github.com/pfacheris/kickback/db/db"
)

type User struct {
  Id           int64 `json:"id"`
  Email        string `json:"email"  binding:"required" sql:"size:255;not null;unique"`
  AuthToken    string `json:"-"  binding:"required" sql:"size:255;not null"`
  CreatedAt    time.Time
  UpdatedAt    time.Time
  DeletedAt    time.Time
}

// This method implements binding.Validator and is executed by the binding.Validate middleware
func (user User) Validate(errors *binding.Errors, req *http.Request) {
    if len(user.Email) < 3 {
        errors.Fields["email"] = "Too short; minimum 3 characters"
    } else if len(user.Email) > 120 {
        errors.Fields["email"] = "Too long; maximum 255 characters"
    }
}

// Methods
func (u *User) save() *User, error {
  if err := DB.save(u).Error; err != nil {
    return u, err
  }
  return u, nil
}
