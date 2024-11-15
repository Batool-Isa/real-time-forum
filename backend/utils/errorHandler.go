package utils

import (
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	//templateFile := "error.html"
	data := struct {
		StatusCode  int
		StatusText  string
		Description string
	}{
		StatusCode:  status,
		StatusText:  http.StatusText(status),
		Description: getErrorDescription(status),
	}
	RenderTemplate(w, r, "error.html", data)

}

func getErrorDescription(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "Sorry, the request is invalid or incomplete."
	case http.StatusNotFound:
		return "Sorry, the page you are looking for might be missing or the URL is incorrect."
	case http.StatusInternalServerError:
		return "Sorry, something went wrong on our end. We are working to fix the issue."
	case http.StatusUnauthorized: // Error 401
		return "Sorry, you are not authorized to access this page. Please log in to continue."
	case http.StatusForbidden: //Error 403
		return "Sorry, you do not have permission to perform this action. Please check your access rights."
	default:
		return "An unexpected error occurred."
	}
}
