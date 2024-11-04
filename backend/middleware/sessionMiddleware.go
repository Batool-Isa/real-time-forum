package middleware

import(
	"net/http"
	"real-time-forum/middleWare"
	"fmt"
)
function SessionMiddleware(next http.Handler) http.Handler{


	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get the cookie
		sessionCook , err := r.Cookie("session_id")
		if err != nil {
			fmt.Fprintf("Error geting the session cookie")
			return
		}
		//get the session info
		userSession , err = GetSession()


		//check if the session is still valid and not expired
	})
}