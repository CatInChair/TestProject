package db_controller

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	randomkey "github.com/swloopit/golang-random-key"
)

var Db *sql.DB

/*
 * Transaction data structure
 */
type Movement struct {
	from      string
	to        string
	amount    float64
	timestamp int64
}

/*
 * Serialize( Movement ) public
 * Makes map[string]string from Movement structure
 */
func Serialize(t Movement) map[string]string {
	return map[string]string{"from": t.from, "to": t.to, "amount": strconv.FormatFloat(t.amount, 'f', 4, 64), "timestamp": strconv.FormatInt(t.timestamp, 10)}
}

/*
 *	fillWallets() private
 *	Creates 10 wallets in db
 */
func fillWallets() {
	for i := 0; i < 10; i++ {
		key := randomkey.CreateRandomKey(24)
		_, err := Db.Exec("INSERT INTO wallets VALUES (?, 100)", key)

		if err != nil {
			i--
			//Маловероятно, что хотя бы 2 кода из 10 при длине слов в 24 символа совпадут
			// Незамысловатое решение, бд не даст создать в таблице wallets две записи с идентичными адресами,
			// о чем с радостью сообщит нам возникающая при выполнении запроса ошибка; это дает право предположить, что можно сразу генерить новый код
		} else {
			fmt.Println("New wallet", key)
		}

	}
}

/*
 * GetBalance( string ) ( float64, error ) public
 * Return wallet balance to the corresponding provided address
 */
func GetBalance(address string) (float64, error) {
	res, err := Db.Query("SELECT balance FROM wallets WHERE address = ?", address)

	if err != nil {
		return -1, err
	}
	defer res.Close()

	var balance float64
	state := res.Next()

	if state {
		res.Scan(&balance)
		return balance, nil
	}

	return -1, errors.New("wallet does not exist")
}

/*
 * GetLast( int64 ) ( []map[string]string, error ) public
 * Return last time transactions
 */
func GetLast(limit int64) ([]map[string]string, error) {
	res, err := Db.Query("SELECT * FROM transactions ORDER BY timestamp DESC LIMIT ?", limit)

	if err != nil {
		return []map[string]string{}, err
	}
	defer res.Close()

	movements := []map[string]string{}

	for res.Next() {
		var movement Movement

		res.Scan(&movement.from, &movement.to, &movement.amount, &movement.timestamp)
		movements = append(movements, Serialize(movement))
	}

	return movements, nil
}

/*
 * Send( string, string, float64) ( error ) public
 * Create new transaction
 */
func Send(from string, to string, amount float64) error {
	balSender, errSender := GetBalance(from)
	_, errRecipient := GetBalance(to)

	if errSender != nil {
		return errors.New("\"from\" wallet does not exist")
	}

	if errRecipient != nil {
		return errors.New("\"to\" wallet does not exist")
	}

	if balSender < amount {
		return errors.New("insufficient funds on sender wallet")
	}

	if _, err := Db.Exec(`
		UPDATE wallets SET balance = balance - ? WHERE address = ?;
		UPDATE wallets SET balance = balance + ? WHERE address = ?;
	`, amount, from, amount, to); err != nil {
		return err
	}

	if _, err := Db.Exec("INSERT INTO transactions VALUES (?, ?, ?, ?)", from, to, amount, time.Now().Unix()); err != nil {
		return err
	}

	return nil
}

// Prepare database controller for work
func Init() {
	var isDataExist = true

	if _, err := os.Stat("./data.db"); errors.Is(err, os.ErrNotExist) {
		isDataExist = false
	}

	db, err := sql.Open("sqlite3", "./data.db")

	if err != nil {
		log.Fatal(err)
	}

	db.Exec(`
	CREATE TABLE IF NOT EXISTS wallets(address TEXT PRIMARY KEY, balance REAL); 
	CREATE TABLE IF NOT EXISTS transactions(sender TEXT, recepient TEXT, amount REAL, timestamp INTEGER);
	`)

	Db = db

	if !isDataExist {
		fillWallets() // Fill wallets table if database is empty
	}
}
