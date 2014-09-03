package shared

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
)

var activeMap *gorp.DbMap

func InitDb() {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	finalPath, err := filepath.Abs(Config.DbFile)
	CheckErr(err, "Loading DBFile")
	log.Printf("Using db: %s\n", finalPath)
	db, err := sql.Open("sqlite3", finalPath)
	CheckErr(err, "sql.Open failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(Account{}, "accounts").SetKeys(false, "Uid")

	err = dbmap.CreateTablesIfNotExists()
	CheckErr(err, "Create tables failed")

	activeMap = dbmap
}

func FetchAccount(acc Account) *Account {
	existing, err := activeMap.Get(Account{}, acc.Uid)
	if err != nil {
		log.Print("Failed to lookup account")
		return nil
	}

	if existing == nil {
		return nil
	} else {
		return existing.(*Account)
	}
}

func AddUpdateAccount(acc Account) {
	existing := FetchAccount(acc)
	if existing != nil {
		existing.Token = acc.Token
		_, err := activeMap.Update(existing)
		if err != nil {
			log.Printf("Failed to update account: %v\n", err)
			return
		}
	} else {
		err := activeMap.Insert(&acc)
		if err != nil {
			log.Printf("Failed to insert account: %v\n", err)
			return
		}
	}
}

type Account struct {
	Uid   int64
	Token string
}
