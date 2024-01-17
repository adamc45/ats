package clause

type SqlClause struct {
	Sql    string
	Params []string
}

// db.Query expects a generic list of params to bind to the query. This converts a specific type of Params to a generic slice.
func (s *SqlClause) GetParams() []any {
	p := make([]any, len(s.Params))
	for i, param := range s.Params {
		p[i] = param
	}
	return p
}
