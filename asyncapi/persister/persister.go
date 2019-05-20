package persister

import (
	"log"
	"reflect"

	"github.com/uesteibar/scribano/asyncapi/repos/messagesrepo"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/storage/db"
)

// Persister persists analyzed messages and pipes the result through
type Persister struct {
	ChIn         <-chan spec.MessageSpec
	ChOut        chan<- spec.MessageSpec
	messagesRepo *messagesrepo.MessagesRepo
}

// New creates a new Persister
func New(chIn <-chan spec.MessageSpec, chOut chan<- spec.MessageSpec, database db.Database) *Persister {
	return &Persister{
		ChIn:         chIn,
		ChOut:        chOut,
		messagesRepo: messagesrepo.New(database),
	}
}

func isDifferent(msg, newMsg spec.MessageSpec) bool {
	return !reflect.DeepEqual(msg, newMsg)
}

func (p *Persister) persist(msg spec.MessageSpec) error {
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

// Watch for incoming messages and persist them
func (p *Persister) Watch() {
	for msg := range p.ChIn {
		if err := p.persist(msg); err != nil {
			log.Printf("ERROR %s", err)
		} else {
			p.ChOut <- msg
		}
	}
}
