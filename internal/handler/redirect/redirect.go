package redirect

import "net/http"

type URLGetter interface {
	GetURL(alias string) (url string, exists bool)
}

func New(urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("id")

		originalURL, exists := urlGetter.GetURL(alias)
		if !exists {
			http.Error(w, "[ERROR] URL не найден", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	}
}
