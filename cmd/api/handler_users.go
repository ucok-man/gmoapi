// File: cmd/api/users.go
package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ucok-man/gmoapi/internal/data"
	"github.com/ucok-man/gmoapi/internal/validator"
)

// @Summary      Register New User
// @Description  Register a new user account. Upon successful registration, an activation email will be sent containing a token valid for 3 days. The account must be activated before it can be used.
// @Description
// @Description  **Validation Rules:**
// @Description  - Name: Required, max 500 characters
// @Description  - Email: Required, valid email format, must be unique
// @Description  - Password: Required, 8-72 characters
// @Description
// @Description  **Default Permissions:** New users receive `movies:read` permission by default.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      object{name=string, email=string, password=string}  true  "User registration data"
// @Success      202  {object}  object{message=string, user=object{id=int64, created_at=string, name=string, email=string, activated=bool}}  "Registration successful - activation email sent"
// @Failure      400  {object}  object{error=string}  "Bad request - malformed JSON or invalid data types"
// @Failure      422  {object}  object{error=map[string]string}  "Unprocessable entity - validation errors (e.g., duplicate email, weak password)"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /users/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Permissions.AddForUser(user.ID, "movies:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})
	envelope := envelope{
		"message": "an email will be sent to you containing activation instructions",
		"user":    user,
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary      Activate User Account
// @Description  Activate a user account using the token received via email. The token is single-use and expires after 3 days. Once activated, all activation tokens for this user are deleted.
// @Description
// @Description  **Token Format:** 26-character alphanumeric string
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        token  body      object{token=string}  true  "Activation token (26 characters)"
// @Success      200  {object}  object{user=object{id=int64, created_at=string, name=string, email=string, activated=bool}}  "Account activated successfully"
// @Failure      400  {object}  object{error=string}  "Bad request - malformed JSON"
// @Failure      409  {object}  object{error=string}  "Conflict - account has been modified during activation"
// @Failure      422  {object}  object{error=map[string]string}  "Unprocessable entity - invalid or expired token"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /users/activated [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validasi token
	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Ambil user dari token
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Update status aktivasi
	user.Activated = true

	// Simpan update ke database
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Hapus semua activation token milik user
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary      Reset User Password
// @Description  Reset a user's password using a password reset token. The token must be obtained via the `/v1/tokens/password-reset` endpoint and is valid for 45 minutes. Once used, all password reset tokens for this user are deleted.
// @Description
// @Description  **Validation Rules:**
// @Description  - Password: Required, 8-72 characters
// @Description  - Token: Required, 26-character alphanumeric string
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        data  body      object{password=string, token=string}  true  "New password and reset token"
// @Success      200  {object}  object{message=string}  "Password reset successfully"
// @Failure      400  {object}  object{error=string}  "Bad request - malformed JSON"
// @Failure      409  {object}  object{error=string}  "Conflict - account has been modified during password reset"
// @Failure      422  {object}  object{error=map[string]string}  "Unprocessable entity - invalid/expired token or weak password"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /users/password [put]
func (app *application) updateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the user's new password and password reset token.
	var input struct {
		Password       string `json:"password"`
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidatePasswordPlaintext(v, input.Password)
	data.ValidateTokenPlaintext(v, input.TokenPlaintext)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve the details of the user associated with the password reset token,
	// returning an error message if no matching record was found.
	user, err := app.models.Users.GetForToken(data.ScopePasswordReset, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired password reset token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Set the new password for the user.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the updated user record in our database, checking for any edit conflicts as
	// normal.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If everything was successful, then delete all password reset tokens for the user.
	err = app.models.Tokens.DeleteAllForUser(data.ScopePasswordReset, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the user a confirmation message.
	env := envelope{"message": "your password was successfully reset"}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
