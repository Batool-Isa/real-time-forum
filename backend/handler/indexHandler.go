package handler

import (
	"net/http"
	"real-time-forum/backend/utils"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	// Pass nil since there's no data to send
	if err := utils.RenderTemplate(w, r, "index.html", nil); err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		return
	}
}
