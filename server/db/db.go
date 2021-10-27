package db

import (
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/smallwood/chatapp-server/msg"
)

var instance *DB
var once sync.Once

func Handler() *DB {
	once.Do(func() {
		instance = newHandler()
	})
	return instance
}

type DB struct {
	lite *sqlx.DB
}

func newHandler() *DB {
	db, err := sqlx.Connect("sqlite3", "msg.db")
	if err != nil {
		log.Fatal("failed to start sqlite db:", err.Error())
	}

	handler := &DB{lite: db}

	handler.CreateTable()

	return handler
}

func (db *DB) CreateTable() error {
	create := `CREATE TABLE IF NOT EXISTS messages (
          sender TEXT NOT NULL,
          msg TEXT
        );`

	_, err := db.lite.Exec(create)
	return err
}

func (db *DB) InsertMessage(message *msg.Msg) error {
	insert := `INSERT INTO messages (
          sender, msg
        ) VALUES (
          $1, $2
        );`

	_, err := db.lite.Exec(insert, message.GetSender().GetName(), message.GetMsg())
	return err
}

func (db *DB) SelectAllMessages() (*msg.Msgs, error) {
	query := `SELECT * FROM messages;`

	var msgs []struct {
		Sender string `db:"sender,omitempty"`
		Msg    string `db:"msg,omitempty"`
	}
	err := db.lite.Select(&msgs, query)

	messages := &msg.Msgs{}
	for _, m := range msgs {
		message := &msg.Msg{
			Sender: &msg.User{Name: m.Sender},
			Msg:    m.Msg,
		}
		messages.Msgs = append(messages.Msgs, message)
	}

	return messages, err
}
