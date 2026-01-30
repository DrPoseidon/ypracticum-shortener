package redirect

import "net/http"

type UrlGetter interface {
	GetUrl(alias string) (url string, exists bool)
}

func New(urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("id")

		// originalURL, exists := urlStorage[hash]
		originalURL, exists := urlGetter.GetUrl(alias)
		if !exists {
			http.Error(w, "[ERROR] URL не найден", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	}
}
