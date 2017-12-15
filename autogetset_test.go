package autosetget_test

import (
	"fmt"
	"testing"

	"reflect"

	. "github.com/xw340721/autoGetSet"
)

type Data struct {
	*DataImpl
	Id   int    `privilege:"public"`
	Name string `privilege:"private"`
}

func (p *Data) SetName(name string) {
	p.Name = name
}

func (p *Data) GetName() string {
	return p.Name
}

func NewData() *Data {

	data := Data{
		DataImpl: NewDataImpl(),
		Id:       100,
		Name:     "haha",
	}

	rt := reflect.TypeOf(data)

	rv := reflect.ValueOf(&data)

	data.SetField(rt, rv)

	return &data
}

func TestRun(t *testing.T) {

	var id int

	data := NewData()

	id = data.Get("Id").(int)

	if id != 100 {
		t.Fail()
	}

	fmt.Println(id)

	name := data.Get("Name").(string)

	if name != "haha" {
		t.Fail()
	}

	fmt.Println(name)

	newId := 101

	data.Set("Id", newId)

	id = data.Get("Id").(int)

	if id != newId {
		t.Fail()
	}
	fmt.Println(data.Get("Id"))

	name2 := "haha2"

	data.Set("Name", name2)

	if data.Get("Name").(string) != name2 {
		t.Fail()
	}

	fmt.Println(data.Get("Name"))

}
