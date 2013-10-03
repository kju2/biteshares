package biteshares

import (
	"appengine"
	"appengine/datastore"
	"html/template"
	"net/http"
)

type Food struct {
	Name string
}

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles("biteshares.html"))

	http.HandleFunc("/", bitshares)
}

func bitshares(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if name := r.FormValue("food"); name != "" {
			_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Food", nil), &Food{name})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {

			for encodedKey, _ := range r.Form {
				/*fmt.Fprintln(w, key)*/
				key, err := datastore.DecodeKey(encodedKey)
				if err != nil {
					//TODO(): log error
					/*http.Error(w, err.Error(), http.StatusInternalServerError)*/
					continue
				}
				err = datastore.Delete(c, key)
				if err != nil {
					//TODO(): log error
					/*http.Error(w, err.Error(), http.StatusInternalServerError)*/
					continue
				}
			}
			/*http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)*/
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	stash := make(map[string]string, 10)
	query := datastore.NewQuery("Food")

	for t := query.Run(c); ; {
		var food Food
		key, err := t.Next(&food)

		if err == datastore.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stash[key.Encode()] = food.Name
	}

	err := templates.ExecuteTemplate(w, "biteshares.html", stash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
