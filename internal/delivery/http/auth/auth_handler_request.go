package auth

import "github.com/datpham/user-service-ms/internal/util"

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (r *SignupRequest) Validate() error {
	if err := util.ValidateEmail(r.Email); err != nil {
		return err
	}

	if err := util.ValidatePassword(r.Password); err != nil {
		return err
	}

	return nil
}
