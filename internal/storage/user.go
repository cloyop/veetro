package storage

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Append Max Open Offers if customer && max Open Application if employee

type User struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password,omitempty"`
	Customer       bool   `json:"isCustomer" bson:"is_customer"`
	Employee       bool   `json:"isEmployee" bson:"is_employee"`
	ResumeUploaded bool   `json:"resumeUploaded" bson:"resume_uploaded"`
	Location       string `json:"location"`
	Created        int64  `json:"created"`
	// MarkDown string format
	Resume string `json:"resume"`
}

func NewUser(name, email, password, location string, isEmployee, isCustomer bool) *User {
	return &User{
		Id:       uuid.NewString(),
		Created:  time.Now().Unix(),
		Password: password,
		Email:    email,
		Name:     name,
		Location: location,
		Customer: isCustomer,
		Employee: isEmployee,
	}
}
func (u *User) Validate() (errs []string) {
	if !u.validName() {
		errs = append(errs, "Invalid name length")
	}
	if !u.validPassword() {
		errs = append(errs, "Invalid password length")
	}
	if !u.validEmail() {
		errs = append(errs, "Invalid Mail")
	}
	if !validLocation(u.Location) {
		errs = append(errs, "Invalid location")
	}
	if !u.Employee && !u.Customer {
		errs = append(errs, "user need to be employee or customer")
	}
	if u.Employee && u.Customer {
		errs = append(errs, "user cant be employee and customer")
	}
	if u.Resume != "" {
		if u.validResume() {
			u.ResumeUploaded = true
		} else {
			errs = append(errs, "Invalid Resume")
		}
	}
	return
}

func (u *User) ParseUpdateFields(us *User) (*map[string]interface{}, *[]string) {
	errs := []string{}
	edit := map[string]interface{}{}
	if u.Name != "" {
		if !u.validName() {
			errs = append(errs, "Invalid name ")
		} else {
			edit["name"] = u.Name
		}
	}
	if u.Password != "" {
		if u.validPassword() && us.CheckPassword(u.Password) != nil {
			u.HashPassword()
			edit["password"] = u.Password
		} else {
			errs = append(errs, "Invalid password")
		}
	}
	if u.Email != "" {
		if !u.validEmail() {
			errs = append(errs, "Invalid Mail")
		} else {
			edit["email"] = u.Email
		}
	}
	if u.Location != "" {
		if !validLocation(u.Location) {
			errs = append(errs, "Invalid location")
		} else {
			edit["location"] = u.Location
		}
	}
	if u.Resume != "" {
		if len(u.Resume) > 100 {
			edit["resume"] = u.Resume
			edit["resume_uploaded"] = true
		} else {
			errs = append(errs, "Invalid Resume")
		}
	}
	return &edit, &errs
}
func (u *User) validName() bool {
	return len(u.Name) >= 8 && len(u.Name) <= 50
}
func (u *User) validPassword() bool {
	return len(u.Password) >= 12 && len(u.Password) <= 70
}
func (u User) validEmail() bool {
	return len(u.Email) >= 6 && len(u.Email) <= 255 && emailPattern.MatchString(u.Email)
}
func (u User) validResume() bool {
	return len(u.Resume) > 300 && strings.Contains(u.Resume, u.Email)
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
func (u *User) CheckPassword(p string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
}
