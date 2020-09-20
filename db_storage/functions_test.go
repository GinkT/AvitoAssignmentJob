package db_storage

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	sqlStatement :=  `SELECT balance FROM public."users" WHERE`

	expectedIntValue := 1

	rows := sqlmock.NewRows([]string{"One"}).AddRow(expectedIntValue)

	mock.ExpectQuery(sqlStatement).WillReturnRows(rows).WithArgs("1")
	test, _ := GetBalance(db, "1", "")
	fmt.Println(test)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

// Интерфейс чтобы мок не ругался на сгенерированное внутри функции время
type Any struct{}

func (a Any) Match(v driver.Value) bool {
	return true
}

func TestAddTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	sqlStatement :=  `INSERT INTO public."transactions"`

	mock.ExpectExec(sqlStatement).WithArgs(
		driver.Value("transfer"),
		driver.Value("1"),
		driver.Value("3"),
		driver.Value("350"),
		Any{})
	AddTransaction(db, "transfer", "1", "3", "350")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetHistoryForId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	sqlStatement := `
			SELECT type, sender, receiver, amount, time 
			FROM public."transactions"
		`
	rows := sqlmock.NewRows([]string{"One", "Two", "Three", "Four", "Five"}).AddRow("1", "2", "3", "4", "5")
	mock.ExpectQuery(sqlStatement).WillReturnRows(rows).WithArgs(driver.Value(driver.Value("1")))
	GetHistoryForId(db, "1", "", "")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestBalanceChange(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	sqlStatement :=  `INSERT INTO public."users"`

	mock.ExpectExec(sqlStatement).WithArgs(driver.Value("2"), driver.Value("350"))
	BalanceChange(db, "2", "350")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestNullIntCheck(t *testing.T) {
	nullIntCheckData := []string{"", "0", "1", ""}
	testResults := []interface{}{sql.NullInt64{}, "0", "1", sql.NullInt64{}}

	for idx, _ := range nullIntCheckData {
		if NullIntCheck(nullIntCheckData[idx]) != testResults[idx] {
			t.Errorf("Bad assert! Expected: %s, got: %s", testResults[idx], nullIntCheckData[idx])
		}
	}
}