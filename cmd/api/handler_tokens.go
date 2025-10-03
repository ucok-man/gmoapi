package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ucok-man/gmoapi/internal/data"
	"github.com/ucok-man/gmoapi/internal/validator"
)

// @Summary      Generate Authentication Token
// @Description  Generate a bearer token for authentication. The token is valid for 24 hours and must be included in the Authorization header for protected endpoints.
// @Description
// @Description  **Usage:** Include the token in subsequent requests as: `Authorization: Bearer {token}`
// @Description
// @Description  **Token Lifetime:** 24 hours
// @Tags         Tokens
// @Accept       json
// @Produce      json
// @Param        credentials  body      object{email=string, password=string}  true  "User login credentials"
// @Success      201  {object}  object{authentication_token=object{token=string, expiry=string}}  "Token generated successfully"
// @Failure      400  {object}  object{error=string}  "Bad request - malformed JSON"
// @Failure      401  {object}  object{error=string}  "Unauthorized - invalid email or password"
// @Failure      422  {object}  object{error=map[string]string}  "Unprocessable entity - validation errors"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /tokens/authentication [post]
func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the email and password from the request body.
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the email and password provided by the client.
	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Lookup the user record based on the email address.
	// If no matching user was found, send 401 Unauthorized.
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// If the passwords don't match, return Unauthorized.
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// Otherwise, generate a new token with 24-hour expiry and scope 'authentication'.
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Encode the token to JSON and send it in the response with 201 Created.
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary      Request Password Reset Token
// @Description  Request a password reset token to be sent via email. The token is valid for 45 minutes and can only be used once. The user account must be activated to receive a reset token.
// @Description
// @Description  **Email Delivery:** Token is sent to the email address registered in the system (not the one provided in request).
// @Description
// @Description  **Token Lifetime:** 45 minutes
// @Tags         Tokens
// @Accept       json
// @Produce      json
// @Param        email  body      object{email=string}  true  "User email address"
// @Success      202  {object}  object{message=string}  "Password reset email will be sent"
// @Failure      400  {object}  object{error=string}  "Bad request - malformed JSON"
// @Failure      422  {object}  object{error=map[string]string}  "Unprocessable entity - email not found or account not activated"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /tokens/password-reset [post]
func (app *application) createPasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the user's email address.
	var input struct {
		Email string `json:"email"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Try to retrieve the corresponding user record for the email address. If it can't
	// be found, return an error message to the client.
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return an error message if the user is not activated.
	if !user.Activated {
		v.AddError("email", "user account must be activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Otherwise, create a new password reset token with a 45-minute expiry time.
	token, err := app.models.Tokens.New(user.ID, 45*time.Minute, data.ScopePasswordReset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Email the user with their password reset token.
	app.background(func() {
		data := map[string]any{
			"passwordResetToken": token.Plaintext,
		}
		// Since email addresses MAY be case sensitive, notice that we are sending this
		// email using the address stored in our database for the user --- not to the
		// input.Email address provided by the client in this request.
		err = app.mailer.Send(user.Email, "token_password_reset.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})
	// Send a 202 Accepted response and confirmation message to the client.
	env := envelope{"message": "an email will be sent to you containing password reset instructions"}

	err = app.writeJSON(w, http.StatusAccepted, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary      Resend Activation Token
// @Description  Request a new activation token to be sent via email. Useful if the original token expired or was lost. The token is valid for 3 days. This endpoint cannot be used if the account is already activated.
// @Description
// @Description  **Email Delivery:** Token is sent to the email address registered in the system (not the one provided in request).
// @Description
// @Description  **Token Lifetime:** 3 days
// @Tags         Tokens
// @Accept       json
// @Produce      json
// @Param        email  body      object{email=string}  true  "User email address"
// @Success      202  {object}  object{message=string}  "Activation email will be sent"
// @Failure      400  {object}  object{error=string}  "Bad request - malformed JSON"
// @Failure      422  {object}  object{error=map[string]string}  "Unprocessable entity - email not found or already activated"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /tokens/activation [post]
func (app *application) createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the user's email address.
	var input struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Try to retrieve the corresponding user record for the email address. If it can't
	// be found, return an error message to the client.
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return an error if the user has already been activated.
	if user.Activated {
		v.AddError("email", "user has already been activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Otherwise, create a new activation token.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Email the user with their additional activation token.
	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
		}
		// Since email addresses MAY be case sensitive, notice that we are sending this
		// email using the address stored in our database for the user --- not to the
		// input.Email address provided by the client in this request.
		err = app.mailer.Send(user.Email, "token_activation.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	// Send a 202 Accepted response and confirmation message to the client.
	env := envelope{"message": "an email will be sent to you containing activation instructions"}
	err = app.writeJSON(w, http.StatusAccepted, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
