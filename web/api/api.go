package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/builder"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
)

func buildSpec(msgSpecs []spec.MessageSpec) builder.AsyncAPISpec {
	b := builder.SpecBuilder{}
	for _, msg := range msgSpecs {
		b.AddMessage(msg)
	}

	return b.Build()
}

func handleAsyncAPI(w http.ResponseWriter, r *http.Request) {
	repo := messages_repo.New(db.DB{})

	if msgSpecs, err := repo.FindAll(); err == nil {
		spec := buildSpec(msgSpecs)
		json, _ := json.Marshal(spec)

		w.Write([]byte(json))
	} else {
		http.Error(w, http.StatusText(500), 500)
	}
}

func Start() {
	r := chi.NewRouter()
	r.Get("/asyncapi", handleAsyncAPI)

	http.ListenAndServe(":5000", r)
}
