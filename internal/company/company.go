package company

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type Company struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func EditCompany(db *sql.DB, company *Company) (*Company, error) {
	if company == nil {
		return company, errors.New("expected to receive company struct but got nil instead")
	}
	result, err := db.Exec(
		`Update company set name = ? where id = ?`,
		company.Name,
		company.Id,
	)
	if err != nil {
		return company, err
	}
	count, rowsAffectedErr := result.RowsAffected()
	if rowsAffectedErr != nil {
		return company, rowsAffectedErr
	}
	if count == 0 {
		return company, errors.New("no rows were affected by update query")
	}
	return company, nil
}

func GetAllCompanies(db *sql.DB) ([]Company, error) {
	rows, err := db.Query(`select name, id from company order by name desc`)
	if err != nil {
		return make([]Company, 0), err
	}
	return scanRows(rows)
}

func GetCompaniesByName(db *sql.DB, name *string) ([]Company, error) {
	rows, err := db.Query(
		`select name, id from company where name like concat('%', ?, '%') order by name asc`,
		name,
	)
	if err != nil {
		return make([]Company, 0), err
	}
	return scanRows(rows)
}

func GetCompanyById(db *sql.DB, id *int) (Company, error) {
	row := db.QueryRow(
		`select name, id from company where id = ?`,
		id,
	)
	company := Company{}
	err := row.Scan(&company.Name, &company.Id)
	// sql.ErrorNoRows. Handle in some other way perhaps?
	return company, err
}

func GetCompanyByName(db *sql.DB, name *string) (Company, error) {
	row := db.QueryRow(
		`select name, id from company where name = ?`,
		name,
	)
	company := Company{}
	err := row.Scan(&company.Name, &company.Id)
	// sql.ErrorNoRows. Handle in some other way perhaps?
	return company, err
}

func InsertCompany(db *sql.DB, name *string) (Company, error) {
	result, err := db.Exec(
		`insert into company (name) values (?)`,
		name,
	)
	company := Company{}
	if err != nil {
		return company, err
	}
	id, lastInsertErr := result.LastInsertId()
	if lastInsertErr != nil {
		return company, lastInsertErr
	}
	insertedRow := db.QueryRow(
		`select name, id from company where id = ?`,
		id,
	)
	insertErr := insertedRow.Scan(&company.Name, &company.Id)
	if insertErr != nil {
		return company, insertErr
	}
	return company, nil
}

func scanRows(rows *sql.Rows) ([]Company, error) {
	var err error
	companies := make([]Company, 0)
	defer rows.Close()
	for rows.Next() {
		var (
			id   int
			name string
		)
		if err = rows.Scan(&name, &id); err != nil {
		} else {
			companies = append(companies, Company{id, name})
		}
	}
	if err != nil {
		return make([]Company, 0), err
	}
	return companies, nil
}
