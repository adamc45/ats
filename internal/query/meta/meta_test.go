package meta

import (
	"testing"
)

type MetaTest struct {
}

func (m *MetaTest) GetColumnNames() []string {
	return []string{"abc", "def"}
}

func Test_Invalid_Limit(t *testing.T) {
	q := QueryMeta{Limit: 0}
	m := &MetaTest{}
	c := q.GetSqlClause(m)
	if len(c.Params) != 0 {
		t.Errorf(
			`
there were unfulfilled expectations.
Expected Params to be a length of 0 but got %d.
`,
			len(c.Params))
	}
}

func Test_Valid_Limit_And_Default_Offset(t *testing.T) {
	q := QueryMeta{Limit: 5}
	m := &MetaTest{}
	c := q.GetSqlClause(m)
	if len(c.Params) != 2 || c.Params[1] != "5" || c.Params[0] != "0" {
		t.Errorf(
			`
there were unfulfilled expectations.
Expected Params to be a length of 2 but got %d.
Expected Params[0] to be 0 but got %s and Params[1] to be 5 but got %s
`,
			len(c.Params), c.Params[1], c.Params[0],
		)
	}
}

func Test_Invalid_Column_In_Order_By(t *testing.T) {
	q := QueryMeta{
		OrderBy: []OrderBy{
			{
				Column: "99",
				Order:  "asc",
			},
		},
	}
	m := &MetaTest{}
	c := q.GetSqlClause(m)
	if len(c.Sql) != 0 {
		t.Errorf(
			`
there were unfulfilled expectations.
Expected Sql to be empty but got %d. Value: %s
`,
			len(c.Sql),
			c.Sql,
		)
	}
}

func Test_Invalid_Order_In_Order_By(t *testing.T) {
	q := QueryMeta{
		OrderBy: []OrderBy{
			{
				Column: "abc",
				Order:  "ascd",
			},
		},
	}
	m := &MetaTest{}
	c := q.GetSqlClause(m)
	if len(c.Sql) != 0 {
		t.Errorf(
			`
there were unfulfilled expectations.
Expected Sql to be empty but got a length of %d. Value: %s
`,
			len(c.Sql),
			c.Sql,
		)
	}
}

func Test_Invalid_And_Valid_Order_In_Order_By(t *testing.T) {
	q := QueryMeta{
		OrderBy: []OrderBy{
			{
				Column: "abc",
				Order:  "asc",
			},
			{
				Column: "123",
				Order:  "desc",
			},
		},
	}
	m := &MetaTest{}
	c := q.GetSqlClause(m)
	if c.Sql != "order by abc asc" {
		t.Errorf(
			`
there were unfulfilled expectations.
Expected Sql to be value: order by abc asc.
Got value: %s
`,
			c.Sql,
		)
	}
}

func Test_Fully_Valid_Clause(t *testing.T) {
	q := QueryMeta{
		Limit:  20,
		Offset: 10,
		OrderBy: []OrderBy{
			{
				Column: "abc",
				Order:  "asc",
			},
			{
				Column: "def",
				Order:  "desc",
			},
		},
	}
	m := &MetaTest{}
	c := q.GetSqlClause(m)
	isValid := c.Params[0] == "10" && c.Params[1] == "20"
	if c.Sql != "order by abc asc, def desc limit ?, ?" || !isValid {
		t.Errorf(
			`
there were unfulfilled expectations.
Expected Sql to be value: order by abc asc, def desc limit ?, ?.
Got value: %s.
Expected Params to be [10, 20] but got %s
`,
			c.Sql,
			c.Params,
		)
	}
}
