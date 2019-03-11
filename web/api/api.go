package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/builder"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
)

func getServerSpec() spec.ServerSpec {
	return spec.ServerSpec{
		Name:    "Test server",
		Version: "0.0.1",
	}
}

func buildSpec(msgSpecs []spec.MessageSpec) builder.AsyncAPISpec {
	b := builder.SpecBuilder{}
	for _, msg := range msgSpecs {
		b.AddMessage(msg)
	}

	b.AddServerInfo(getServerSpec())

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

// Start the api server
func Start() {
	r := chi.NewRouter()
	r.Get("/asyncapi", handleAsyncAPI)

	log.Printf("Running api on localhost:5000")
	http.ListenAndServe(":5000", r)
}
