package echo

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestResult_New(t *testing.T) {
	//type fields struct {
	//	Code  int
	//	Msg   string
	//	Data  interface{}
	//	Extra map[string]interface{}
	//}
	//

	ret := NewRetResult()
	m := make(map[string]interface{})
	m["yo~~"] = "hei~~~"
	m["hu~~"] = "pu~~~~"
	ret.Extra = m
	ret.Code = 0
	ret.Data = "oliver loves annabelle~"
	bytes, _ := json.Marshal(ret)
	fmt.Println(string(bytes))

}
