package meta

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/adamc45/ats/internal/query/clause"
)

type OrderBy struct {
	Column string
	Order  string
}

type QueryMeta struct {
	Limit   int
	Offset  int
	OrderBy []OrderBy
}

type ColumnNames interface {
	GetColumnNames() []string
}

type SqlMeta interface {
	GetSqlClause(c ColumnNames) clause.SqlClause
}

func (q *QueryMeta) getLimit() string {
	if q.Limit == 0 {
		return ""
	}
	return strconv.Itoa(q.Limit)
}

func (q *QueryMeta) getLimitClause() clause.SqlClause {
	limit := q.getLimit()
	// Doesn't matter what offset is at this point. An offset with a limit of 0 is meaningless
	if len(limit) == 0 {
		return clause.SqlClause{}
	}
	return clause.SqlClause{
		// order by must be string literal so we can't assign as a param. Documentation says the same for limit but that isn't the case.
		Sql: "limit ?, ?",
		Params: []string{
			q.getOffset(),
			limit,
		},
	}
}

func (q *QueryMeta) getOffset() string {
	return strconv.Itoa(q.Offset)
}

func (q *QueryMeta) getOrderby(c ColumnNames) string {
	if len(q.OrderBy) == 0 {
		return ""
	}
	validColumnNames := c.GetColumnNames()
	orderBy := []string{}
	for _, o := range q.OrderBy {
		// don't have a valid combination of order and column so don't include it
		if len(o.Column) == 0 || len(o.Order) == 0 {
			continue
		}
		order := strings.ToLower(o.Order)
		// don't have a valid order so don't include it
		if order != "asc" && order != "desc" {
			continue
		}
		// vscode go extension has decided slices package isn't currently available. This is equivalent of slices.contains(...)
		sliceList := make([]string, 0)
		for _, c := range validColumnNames {
			if c != o.Column {
				continue
			}
			sliceList = append(sliceList, o.Column)
			break
		}
		// since the column will be coming from user input, we can't simply consume the value otherwise we open ourselves up to sql injection
		// checking that the column is in a list that the caller knows beforehand should be ample protection against this
		isValidColumn := len(sliceList) > 0
		if !isValidColumn {
			continue
		}
		orderBy = append(orderBy,
			fmt.Sprintf("%s %s", o.Column, o.Order),
		)
	}
	// If we end up with no order by clauses, treat it like no order by
	if len(orderBy) == 0 {
		return ""
	}
	return fmt.Sprintf(
		"%s %s",
		"order by",
		strings.Join(orderBy, ", "),
	)
}

func (q *QueryMeta) GetSqlClause(c ColumnNames) clause.SqlClause {
	r := regexp.MustCompile("[ ]+")
	limitClause := q.getLimitClause()
	return clause.SqlClause{
		Sql: strings.TrimSpace(
			r.ReplaceAllString(
				fmt.Sprintf("%s %s", q.getOrderby(c), limitClause.Sql),
				" ",
			),
		),
		Params: limitClause.Params,
	}
}
