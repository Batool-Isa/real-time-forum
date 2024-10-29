package session 

import(
	"net/http"
	"github.com/gorilla/sessions"
	"fmt"
)

var (
	//This key used to encrypt and sign session data
    key = []byte("12345678912345678912345678912356")
	 // create cookie with the key provided earlier
    store = sessions.NewCookieStore(key)
)

//function to create sesssion
func CreatSession(w http.ResponseWriter, r *http.Request)
{
	//create new session
userSession, err := store.Get(r,"User Session")
if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
//set the status of the user session to true
userSession.Values["authenticated"] = true

//save the session to tht cookie
_, err = sessions.Save(w,r);
if err != nil{
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
fmt.Fprintf(w, "Session created and user authenticated")

}

//function to retive the user session
func GetSession(w http.ResponseWriter, r *http.Request)
{
userSession, err := store.Get(r,"User Session")
if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}

auth, ok := userSession.Values["authenticcated"].(bool)
if auth && ok {
	fmt.Fprintf(w,"Welcome")
} else {
	fmt.Fprintf(w,"Log in first")

}


}