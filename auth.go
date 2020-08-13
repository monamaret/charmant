package charm

import (
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
)

// Auth is the authenticated user's charm id and jwt returned from the ssh server
type Auth struct {
	CharmID              string        `json:"charm_id"`
	JWT                  string        `json:"jwt"`
	PublicKey            string        `json:"public_key"`
	EncryptKeys          []*EncryptKey `json:"encrypt_keys"`
	claims               *jwt.StandardClaims
	encryptKeysDecrypted bool
}

func (cc *Client) Auth() (*Auth, error) {
	cc.authLock.Lock()
	defer cc.authLock.Unlock()
	if cc.auth.claims == nil || cc.auth.claims.Valid() != nil {
		auth := &Auth{}
		s, err := cc.sshSession()
		if err != nil {
			return nil, err
		}
		defer s.Close()
		b, err := s.Output("api-auth")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, auth)
		if err != nil {
			return nil, err
		}
		token, err := jwt.ParseWithClaims(auth.JWT, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return cc.jwtPublicKey, nil
		})
		if err != nil {
			return nil, err
		}
		auth.claims = token.Claims.(*jwt.StandardClaims)
		auth.encryptKeysDecrypted = false
		cc.auth = auth
		if err != nil {
			return nil, err
		}
	}
	return cc.auth, nil
}

func (cc *Client) InvalidateAuth() {
	cc.authLock.Lock()
	defer cc.authLock.Unlock()
	cc.auth.claims = nil
}