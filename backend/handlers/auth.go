package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kcsu/store/auth"
	"github.com/labstack/echo/v4"
)

const cookieName = "_token"

// TODO: Rework entirely

func (h *Handler) GetUser(c echo.Context) error {
	userId := h.Auth.GetUserId(c)
	user, err := h.Users.Find(userId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

// OAuth2 callback route handler
func (h *Handler) AuthCallback(c echo.Context) error {
	// Fetch the OAuth2 user data
	authUser, err := h.Auth.CompleteUserAuth(c)
	if err != nil {
		return err
	}

	// Create or fetch the user in the database
	user, err := h.Users.FindOrCreate(&authUser)
	if err != nil {
		// Ensure there is no email address conflict
		exists, exerr := h.Users.Exists(authUser.Email)
		if exerr != nil {
			return exerr
		}
		if exists {
			return echo.NewHTTPError(http.StatusConflict, "email is taken")
		}
		return err
	}

	// Create JWT for login
	claims := &auth.JwtClaims{
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Subject:   strconv.Itoa(int(user.ID)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response
	t, err := token.SignedString([]byte(h.Config.JwtSecret))
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = t
	// TODO: short-lived, refresh tokens
	cookie.Expires = time.Now().Add(24 * time.Hour)
	// Cookie is secure in production
	cookie.Secure = !h.Config.Debug
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	// TODO: Redirect instead?
	return c.Redirect(http.StatusTemporaryRedirect, h.Config.OauthRedirectUrl)
}

// This function needs tests:

// Redirect to the (google) OAuth2 provider
func (h *Handler) AuthRedirect(c echo.Context) error {
	url, err := h.Auth.GetAuthUrl(c)
	if err != nil {
		return echo.ErrBadRequest
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     cookieName,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Secure:   !h.Config.Debug,
		MaxAge:   -1,
	})
	return c.NoContent(http.StatusOK)
}
