package persister

import (
	"log"
	"reflect"

	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
)

type Persister struct {
	ChIn         <-chan spec.MessageSpec
	ChOut        chan<- spec.MessageSpec
	messagesRepo *messages_repo.MessagesRepo
}

func New(chIn <-chan spec.MessageSpec, chOut chan<- spec.MessageSpec, database db.Database) *Persister {
	return &Persister{
		ChIn:         chIn,
		ChOut:        chOut,
		messagesRepo: messages_repo.New(database),
	}
}

func isDifferent(msg, newMsg spec.MessageSpec) bool {
	return !reflect.DeepEqual(msg, newMsg)
}

func (p *Persister) Persist(msg spec.MessageSpec) error {
	ms, err := p.messagesRepo.Find(msg.Topic)

	if err == nil {
		if isDifferent(msg, ms) {
			return p.messagesRepo.Update(msg)
		}
	} else if err.Error() == "NOT_FOUND" {
		return p.messagesRepo.Create(msg)
	}

	return err
}

func (p *Persister) Watch() {
	for msg := range p.ChIn {
		if err := p.Persist(msg); err != nil {
			log.Printf("Error: %s", err)
		} else {
			p.ChOut <- msg
		}
	}
}
