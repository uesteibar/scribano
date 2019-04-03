package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/uesteibar/scribano/asyncapi/builder"
	"github.com/uesteibar/scribano/asyncapi/repos/messages_repo"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/storage/db"
)

func getServerSpec(name string) spec.ServerSpec {
	return spec.ServerSpec{Name: name}
}

func buildSpec(msgSpecs []spec.MessageSpec, exchange string) builder.AsyncAPISpec {
	b := builder.SpecBuilder{}
	for _, msg := range msgSpecs {
		b.AddMessage(msg)
	}

	b.AddServerInfo(getServerSpec(exchange))

	return b.Build()
}

func handleAsyncAPI(w http.ResponseWriter, r *http.Request) {
	repo := messages_repo.New(db.DB{})

	if msgSpecs, err := repo.FindAll(); err == nil {
		spec := buildSpec(msgSpecs, "")
		json, _ := json.Marshal(spec)

		w.Write([]byte(json))
	} else {
		http.Error(w, http.StatusText(500), 500)
	}
}

func handleAsyncAPIForExchange(w http.ResponseWriter, r *http.Request) {
	repo := messages_repo.New(db.DB{})
	exchange := chi.URLParam(r, "exchange")

	if exchange == "" {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if msgSpecs, err := repo.FindByExchange(exchange); err == nil {
		spec := buildSpec(msgSpecs, exchange)
		json, _ := json.Marshal(spec)

		w.Write([]byte(json))
	} else {
		http.Error(w, http.StatusText(500), 500)
	}
}

// Start the api server
func Start() {
	r := chi.NewRouter()
	r.Get("/asyncapi", handleAsyncAPI)
	r.Get("/asyncapi/", handleAsyncAPI)
	r.Get("/asyncapi/{exchange}", handleAsyncAPIForExchange)

	log.Printf("Running api on http://localhost:5000")
	http.ListenAndServe(":5000", r)
}
