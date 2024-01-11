package company

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// In the following tests where the error is also checked from the tested function, ExpectationsWereMet doesn't seem to include Row.Scan errors so we need to check explicitly

func Test_Handles_Edit_Company_Successfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	newName := "abc"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	sqlmock.NewRows([]string{"name", "id"}).
		AddRow("abc", 1)
	// runs the mock query
	mock.ExpectExec(regexp.QuoteMeta(`Update company set name = ? where id = ?`)).
		WithArgs(newName, 1).
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	// runs the real query
	EditCompany(db, &Company{Id: 1, Name: newName})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Error_For_Edit_Company(t *testing.T) {
	db, mock, err := sqlmock.New()
	newName := "abc"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	sqlmock.NewRows([]string{"name", "id"}).
		AddRow("abc", 1)
	// runs the mock query
	mock.ExpectExec(regexp.QuoteMeta(`Update companys set name = ? where id = ?`)).
		WithArgs(newName, 1).
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	// runs the mock query (unsuccessfully since the table name is wrong)
	EditCompany(db, &Company{Id: 1, Name: newName})
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Unchanged_Edit_Company(t *testing.T) {
	db, mock, err := sqlmock.New()
	name := "abc"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	sqlmock.NewRows([]string{"name", "id"}).
		AddRow(name, 1)
	// runs the mock query
	mock.ExpectExec(regexp.QuoteMeta(`Update company set name = ? where id = ?`)).
		WithArgs(name, 1).
		WillReturnResult(
			sqlmock.NewResult(1, 0),
		)
	// runs the real query
	EditCompany(db, &Company{Id: 1, Name: name})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Nil_Company_For_Edit_Company(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	var company *Company = nil
	_, err = EditCompany(db, company)
	if err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Returns_Unfiltered_Companies_Successfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow("abc123", 1).
		AddRow("def456", 2)
	// runs the mock query
	mock.ExpectQuery("select name, id from company").
		WillReturnRows(mockRows)
	// runs the real query
	GetAllCompanies(db)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Error_For_Unfiltered_Companies(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow("abc123", 1).
		AddRow("def456", 2)
	// runs the mock query
	mock.ExpectQuery("select name, id from companys").
		WillReturnRows(mockRows)
	// runs the mock query (unsuccessfully since the table name is wrong)
	GetAllCompanies(db)
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Error_For_Unfiltered_Companies_With_Bad_Entry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow(nil, nil).
		AddRow("def456", 2).
		RowError(1, fmt.Errorf("row error"))
	// runs the mock query
	mock.ExpectQuery("select name, id from company").
		WillReturnRows(mockRows)
	// runs the mock query (unsuccessfully since the table name is wrong)
	_, e := GetAllCompanies(db)
	if err := mock.ExpectationsWereMet(); err != nil || e == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Returns_Filtered_Companies_Successfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	searchTerm := "abc"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow("abc123", 1)
	// runs the mock query
	mock.ExpectQuery(regexp.QuoteMeta(`select name, id from company where name like concat('%', ?, '%')`)).
		WithArgs(searchTerm).
		WillReturnRows(mockRows)
	// runs the real query
	GetCompaniesByName(db, &searchTerm)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Error_For_Filtered_Companies(t *testing.T) {
	db, mock, err := sqlmock.New()
	searchTerm := "abc"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow("abc123", 1)
	// runs the mock query
	mock.ExpectQuery(regexp.QuoteMeta(`select name, id from companys where name like concat('%', ?, '%')`)).
		WithArgs(searchTerm).
		WillReturnRows(mockRows)
	// runs the mock query (unsuccessfully since the table name is wrong)
	GetCompaniesByName(db, &searchTerm)
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Error_For_Filtered_Companies_With_Bad_Entry(t *testing.T) {
	db, mock, err := sqlmock.New()
	searchTerm := "abc"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow(nil, nil).
		AddRow("abc123", 1).
		RowError(1, fmt.Errorf("row error"))
	// runs the mock query
	mock.ExpectQuery(regexp.QuoteMeta(`select name, id from company where name like concat('%', ?, '%')`)).
		WithArgs(searchTerm).
		WillReturnRows(mockRows)
	// runs the real query
	_, e := GetCompaniesByName(db, &searchTerm)
	if err := mock.ExpectationsWereMet(); err != nil || e == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Returns_Filtered_Company_Successfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	searchTerm := "abc123"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow(searchTerm, 1)
	// runs the mock query
	mock.ExpectQuery(regexp.QuoteMeta(`select name, id from company where name = ?`)).
		WithArgs(searchTerm).
		WillReturnRows(mockRows)
	// runs the real query
	GetCompanyByName(db, &searchTerm)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Error_For_Filtered_Company(t *testing.T) {
	db, mock, err := sqlmock.New()
	searchTerm := "abc123"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow(searchTerm, 1)
	// runs the mock query
	mock.ExpectQuery(regexp.QuoteMeta(`select name, id from companys where name = ?`)).
		WithArgs(searchTerm).
		WillReturnRows(mockRows)
	// runs the mock query (unsuccessfully since the table name is wrong)
	GetCompanyByName(db, &searchTerm)
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Error_For_Filtered_Company_With_Bad_Entry(t *testing.T) {
	db, mock, err := sqlmock.New()
	searchTerm := "abc123"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow(nil, nil).
		AddRow(searchTerm, 1).
		RowError(1, fmt.Errorf("row error"))
	// runs the mock query
	mock.ExpectQuery(regexp.QuoteMeta(`select name, id from company where name = ?`)).
		WithArgs(searchTerm).
		WillReturnRows(mockRows)
	// runs the real query
	_, e := GetCompanyByName(db, &searchTerm)
	if err := mock.ExpectationsWereMet(); err != nil || e == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Returns_Inserted_Company_Successfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	newName := "abc123"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectExec(regexp.QuoteMeta(`insert into company (name) values (?)`)).
		WithArgs(newName).
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow(newName, 1)
	// runs the mock query
	mock.ExpectQuery(regexp.QuoteMeta(`select name, id from company where id = ?`)).
		WithArgs(1).
		WillReturnRows(mockRows)
	InsertCompany(db, &newName)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Insert_Error_For_Company(t *testing.T) {
	db, mock, err := sqlmock.New()
	newName := "abc123"
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectExec(regexp.QuoteMeta(`insert into companys (name) values (?)`)).
		WithArgs(newName).
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow(newName, 1)
	// runs the mock query (unsuccessfully since the table name is wrong)
	mock.ExpectQuery(regexp.QuoteMeta(`select name, id from companys where id = ?`)).
		WithArgs(1).
		WillReturnRows(mockRows)
	InsertCompany(db, &newName)
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Handles_Insert_Error_For_Company_With_Bad_Entry(t *testing.T) {
	db, mock, err := sqlmock.New()
	var newName *string = nil
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectExec(regexp.QuoteMeta(`insert into company (name) values (?)`)).
		WithArgs(newName).
		WillReturnResult(
			sqlmock.NewResult(1, 1),
		)
	// adds rows
	mockRows := sqlmock.NewRows([]string{"name", "id"}).
		AddRow(nil, nil)
	// runs the mock query
	mock.ExpectQuery(regexp.QuoteMeta(`select name, id from company where id = ?`)).
		WithArgs(1).
		WillReturnRows(mockRows)
	_, e := InsertCompany(db, newName)
	if err := mock.ExpectationsWereMet(); err != nil || e == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
