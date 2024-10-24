package session 

import(
	"net/http"
	"github.com/gorilla/sessions"
	"fmt"
)

var (
    key = []byte("12345678912345678912345678912356")
    store = sessions.NewCookieStore(key)
)

func getSession(w http.ResponseWriter, r *http.Request)
{
session, err := store.Get(r,"User Session")
if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}

auth, ok := session.Value["authenticcated"].(bool)
if auth && ok {
	fmt.Fprintf(w,"Welcome")
} else {
	fmt.Fprintf(w,"Log in first")

}


}