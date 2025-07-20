package publicEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"context"
	"database/sql"
	"errors"
	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

func (r *Handler) HandleLoginCallback(c *fiber.Ctx) error {
	// * parse body
	body := new(payload.OauthCallback)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * exchange code for token
	token, err := r.Oauth2Config.Exchange(context.Background(), *body.Code)
	if err != nil {
		return gut.Err(false, "failed to exchange code for token", err)
	}

	// * parse id token from oauth2 token
	userInfo, err := r.OidcProvider.UserInfo(context.TODO(), oauth2.StaticTokenSource(token))
	if err != nil {
		return gut.Err(false, "failed to get user info", err)
	}

	// * parse user claims
	oidcClaims := new(common.OidcClaims)
	if err := userInfo.Claims(oidcClaims); err != nil {
		return gut.Err(false, "failed to parse user claims", err)
	}

	// * find user with oid
	user, err := r.database.P().UserGetByOid(c.Context(), oidcClaims.Id)
	if err != nil {
		// * if user not exist, create new user
		if errors.Is(err, sql.ErrNoRows) {
			user, err = r.database.P().UserCreate(c.Context(), &psql.UserCreateParams{
				Oid:        oidcClaims.Id,
				Firstname:  oidcClaims.FirstName,
				Lastname:   oidcClaims.Lastname,
				Email:      oidcClaims.Email,
				PictureUrl: oidcClaims.Picture,
			})
			if err != nil {
				return gut.Err(false, "failed to create user", err)
			}
		} else {
			return gut.Err(false, "failed to query user", err)
		}
	}

	// * generate jwt token
	claims := &common.LoginClaims{
		UserId: user.Id,
	}

	// * sign jwt token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJwtToken, err := jwtToken.SignedString([]byte(*r.config.Secret))
	if err != nil {
		return gut.Err(false, "failed to sign jwt token", err)
	}

	// * set cookie
	c.Cookie(&fiber.Cookie{
		Name:  "login",
		Value: signedJwtToken,
	})

	return c.JSON(response.Success(c, map[string]string{
		"token": signedJwtToken,
	}))
}
