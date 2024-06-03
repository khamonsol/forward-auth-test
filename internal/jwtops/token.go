package jwtops

import "github.com/SoleaEnergy/forwardAuth/internal/policy"

var SOLEA_AUTH_HEADERS = []string{"x-access-token",
	"X-Access-Token",
	"Authorization"}

type AccessToken interface {
	Validate(policy.Policy) error
	GetShortName() string
}
