package handlers

import "github.com/gorilla/sessions"


// SetStore sets the session store for handlers
func SetStore(s *sessions.CookieStore) {
	store = s
}
