package clause

import (
	"testing"
)

func Test_Get_Params_Count(t *testing.T) {
	clause := SqlClause{Params: []string{"123", "abc"}}
	params := clause.GetParams()
	if len(params) != 2 {
		t.Errorf("there were unfulfilled expectations. Expected length of params to be but got %d.", len(params))
	}
}

func Test_Get_Params_Order(t *testing.T) {
	clause := SqlClause{Params: []string{"123", "abc"}}
	params := clause.GetParams()
	if params[0] != "123" || params[1] != "abc" {
		t.Errorf("there were unfulfilled expectations. Expected params to be in the same passed order but got %s, %s.", params[0], params[1])
	}
}
