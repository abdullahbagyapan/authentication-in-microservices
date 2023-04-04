package app

import (
	"github.com/google/uuid"
	"log"
	"microservices/auth-service/internal/ports"
)

type Adapter struct {
	core      ports.CorePort
	msgBroker ports.MsgBrokerPort
	db        ports.DbPort
}

func NewAdapter(core ports.CorePort, msgBroker ports.MsgBrokerPort, db ports.DbPort) *Adapter {
	return &Adapter{msgBroker: msgBroker, db: db, core: core}
}

func (A Adapter) Login(username, password string) (string, bool, error) {

	hashedPassword := A.core.Hash(password)

	user, err := A.db.FindByUsername(username)

	if err != nil {
		log.Printf("Error querying database %v", err)
		return "", false, err
	}

	if user.Password == hashedPassword {
		return user.Id, true, nil
	}
	return "", false, nil

}

func (A Adapter) Register(name, username, password, email string) (string, error) {

	hashedPassword := A.core.Hash(password)

	var user = ports.User{Id: uuid.New().String(), Password: hashedPassword, Username: username, Name: name, Email: email}

	err := A.db.SaveUser(&user)

	if err != nil {
		log.Printf("Error saving user %v", err)
		return "", err
	}

	go func() {

		var user = ports.MsgBrokerUserInfo{Email: email, Name: name}
		err := A.msgBroker.PublishMessage(&user)

		if err != nil {
			log.Printf("Error publishing message %v", err)
		}
	}()

	return user.Id, nil

}
