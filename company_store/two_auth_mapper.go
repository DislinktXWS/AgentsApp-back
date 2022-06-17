package company_store

import (
	"github.com/sec51/twofactor"
)

func TwoAuthMapper(username string, totp *twofactor.Totp) TwoFactorAuth {
	bytes, _ := totp.ToBytes()
	auth := TwoFactorAuth{
		Username: username,
		Totp:     bytes,
	}
	return auth
}
