package authentication

import "net"

type AuthHandler struct {
	Connected map[string]net.Conn
}

func NewHandler() AuthHandler {
	return AuthHandler{Connected: make(map[string]net.Conn)}
}

func (a *AuthHandler) LoginUser(conn net.Conn, username string) bool {
	if _, has := a.Connected[username]; has {
		return false
	}
	for k, v := range a.Connected {
		if v == conn {
			delete(a.Connected, k)
		}
	}
	a.Connected[username] = conn
	return true
}