package main

import "github.com/gin-gonic/gin"

// TODO: Make sure there's no duplicates. Return UI-friendly messages.
// Error names should be consistent (e.g. use 'invalid' in place of 'bad').

var (
	ErrMissingFields = gin.H{
		"code": "420",
		"msg":  "missing required fields",
	}

	ErrInvalidEmail = gin.H{
		"code": "421",
		"msg":  "email must be a valid ufl.edu address",
	}

	ErrEmailTaken = gin.H{
		"code": "422",
		"msg":  "email address has already been registered",
	}

	ErrInvalidCaptcha = gin.H{
		"code": "423",
		"msg":  "bad or expired captcha token",
	}

	ErrDbFailure = gin.H{
		"code": "424",
		"msg":  "internal server error",
	}

	ErrEmailFailure = gin.H{
		"code": "425",
		"msg":  "internal server error",
	}

	ErrBadRespType = gin.H{
		"code": "426",
		"msg":  "invalid response type",
	}

	ErrInvalidScope = gin.H{
		"code": "427",
		"msg":  "invalid scope",
	}

	ErrInvalidState = gin.H{
		"code": "428",
		"msg":  "invalid state",
	}

	ErrDifferentRedirectURI = gin.H{
		"code": "429",
		"msg":  "redirect URI differs from client configuration",
	}

	ErrIncorrectUserPass = gin.H{
		"code": "430",
		"msg":  "incorrect username or password",
	}

	ErrNewJWT = gin.H{
		"code": "431",
		"msg":  "internal server error",
	}

	ErrMalformedURLParams = gin.H{
		"code": "432",
		"msg":  "malformed or missing required URL parameters",
	}

	ErrGrantTokenNotFound = gin.H{
		"code": "435",
		"msg":  "grant token is invalid or expired",
	}

	ErrGrantWrongClient = gin.H{
		"code": "436",
		"msg":  "client redirect_uri or client_id is incorrect",
	}

	ErrCreateAccessToken = gin.H{
		"code": "437",
		"msg":  "internal server error",
	}

	ErrCreateRefreshToken = gin.H{
		"code": "438",
		"msg":  "internal server error",
	}

	ErrInvalidRefreshToken = gin.H{
		"code": "439",
		"msg":  "invalid or revoked refresh token",
	}

	ErrRefreshWrongClient = gin.H{
		"code": "440",
		"msg":  "refresh token was not issued for the current client",
	}

	ErrInvalidRealm = gin.H{
		"code": "441",
		"msg":  "you are not authorized to access this route",
	}

	ErrInvalidClientType = gin.H{
		"code": "442",
		"msg":  "invalid client type, expected 'confidential' or 'public'",
	}

	ErrInvalidPublicScope = gin.H{
		"code": "443",
		"msg":  "invalid scope for public client",
	}

	ErrClientNameTaken = gin.H{
		"code": "444",
		"msg":  "client name has already been registered",
	}

	ErrInvalidRedirectURI = gin.H{
		"code": "445",
		"msg":  "invalid redirect_uri",
	}

	ErrRedirectURITaken = gin.H{
		"code": "446",
		"msg":  "redirect_uri has already been registered",
	}

	ErrInvalidClientInfo = gin.H{
		"code": "447",
		"msg":  "client name and description must be alphanumeric",
	}

	ErrHashError = gin.H{
		"code": "448",
		"msg":  "internal server error",
	}

	ErrVerifyEmail = gin.H{
		"code": "449",
		"msg":  "please check your email address for a verification email. Email addresses must be verified at least once per semester to maintain student status",
	}

	ErrEmailAlreadySent = gin.H{
		"code": "450",
		"msg":  "verification email already sent. Please check your inbox & spam folders",
	}

	ErrInvalidURLParam = gin.H{
		"code": "451",
		"msg":  "invalid URL parameters",
	}

	ErrNoCookie = gin.H{
		"code": "452",
		"msg":  "could not extract auth cookie",
	}

	ErrUserNotFound = gin.H{
		"code": "453",
		"msg":  "could not find user",
	}

	ErrPasswordChanged = gin.H{
		"code": "454",
		"msg":  "password changed",
	}

	ErrVerificationRequired = gin.H{
		"code": "455",
		"msg":  "email verification required. please log in again",
	}

	ErrClientNotFound = gin.H{
		"code": "456",
		"msg":  "client id not found",
	}

	ErrInvalidGrantType = gin.H{
		"code": "457",
		"msg":  "missing or invalid grant type",
	}

	ErrTokenExpired = gin.H{
		"code": "458",
		"msg":  "token has expired",
	}

	ErrTokenUnauthorized = gin.H{
		"code": "459",
		"msg":  "unauthorized or invalid auth header",
	}

	ErrWrongTokenType = gin.H{
		"code": "460",
		"msg":  "expected bearer header token type",
	}

	ErrInsufficientPermission = gin.H{
		"code": "461",
		"msg":  "this client does not have sufficient permission to access this route at this time",
	}
)