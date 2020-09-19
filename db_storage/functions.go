package db_storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"time"

	"log"
)

// Пополнение баланса. В случае если юзера не было в таблице - создаётся новая запись.
// Если юзер был - выполняет UPDATE баланса и прибавляет сумму пополнения.
// Соответственно, если в начале amount стоит минус - убавляет.
func BalanceChange(db *sql.DB, id, amount string) (sql.Result, error) {
	sqlStatement := `
			INSERT INTO public."users"
			VALUES($1, $2)
			ON CONFLICT ("id")
			DO
			UPDATE SET "balance" = "users"."balance" + $2
		`

	res, err := db.Exec(sqlStatement, id, amount)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Получение баланса. Выполняет запрос к БД и получает значение баланса пользователя.
// При указанном параметре currency - выполняет конвертацию.
func GetBalance(db *sql.DB, id, currency string) (float64, error){
	var conversionRate float64 = 1
	if currency != "" {
		var err error
		conversionRate, err = GetConversionRate(currency)
		if err != nil {
			return 0, err
		}
	}

	var balance float64

	row := db.QueryRow(`SELECT balance FROM public."users" WHERE id = $1`, id)
	if err := row.Scan(&balance); err != nil {
		return 0, err
	}

	return balance / conversionRate, nil
}

// Использует API валютного сайта для того чтобы получить коэффициент обмена РУБЛЯ на указанную валюту
func GetConversionRate(currency string) (float64, error) {
	resp, err := http.Get("https://api.exchangeratesapi.io/latest?symbols=RUB&base=" + currency) // Проверить с невалидной валютой
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("Status code was not OK!")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if !json.Valid(body) {
		return 0, errors.New("Invalid JSON!")
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	rate := result["rates"].(map[string]interface{})
	log.Printf("[Rate API] Unmarshalled current rate for %s: %f\n", currency, rate["RUB"])

	return rate["RUB"].(float64), nil
}

// Проверка на Null, возвращает NullInt или пустой интерфейс
func nullIntCheck(str string) interface{} {
	if len(str) == 0 {
		return sql.NullInt64{}
	}
	return str
}

// Добавляет транзакцию в Таблицу транзакций. Имеет проверку значений sender, receiver на NULL
func AddTransaction(db *sql.DB, transactionType, sender, receiver, amount string) error {
	sqlStatement := `
			INSERT INTO public."transactions" (type, sender, receiver, amount, time)
			VALUES($1, $2, $3, $4, $5)
		`

	uTime := time.Now().Unix()

	_, err := db.Exec(sqlStatement, transactionType, nullIntCheck(sender), nullIntCheck(receiver), amount, uTime)
	if err != nil {
		return err
	}
	return nil
}

type Transaction struct {
	Desc 		string				`json:"description"`
	Sender 		sql.NullInt64		`json:"sender"`
	Receiver 	sql.NullInt64		`json:"receiver"`
	Amount 		string				`json:"amount"`
	Time 		string				`json:"time"`
}

// Получение журнала транзакций по ID. Опционально сортировка sortBy по сумме(amount), дате транзакции(time)
// и изменение порядка сортировки orderBy по возрастанию(asc), убыванию(desc).
func GetHistoryForId(db *sql.DB, id, sortBy, orderBy string) []*Transaction {
	sqlStatement := `
			SELECT type, sender, receiver, amount, time 
			FROM public."transactions" 
			WHERE sender = $1 OR receiver = $1 
		`
	if sortBy == "amount" || sortBy == "time" {
		sqlStatement += ` ORDER BY ` + sortBy
		if orderBy == "asc" || orderBy == "desc" {
			sqlStatement += ` ` + orderBy
		}
	}

	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	transactions := make([]*Transaction, 0)
	for rows.Next() {
		trans := new(Transaction)
		err := rows.Scan(&trans.Desc, &trans.Sender, &trans.Receiver, &trans.Amount, &trans.Time)
		if err != nil {
			log.Println("Error doing scan in Get History", err)
			return nil
		}
		transactions = append(transactions, trans)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}
	return transactions
}

