package autosetget_test

import (
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

func reflectUint(b *testing.B) {
	var id int

	data := NewData()

	id = data.Get("Id").(int)

	if id != 100 {
		b.Fail()
	}

	name := data.Get("Name").(string)

	if name != "haha" {
		b.Fail()
	}

	newId := 101

	data.Set("Id", newId)

	id = data.Get("Id").(int)

	if id != newId {
		b.Fail()
	}

	name2 := "haha2"

	data.Set("Name", name2)

	if data.Get("Name").(string) != name2 {
		b.Fail()
	}
}

func originUint(b *testing.B) {

	var id int

	data := &Data{
		Id:   100,
		Name: "haha",
	}

	id = data.Id

	if id != 100 {
		b.Fail()
	}

	name := data.Name

	if name != "haha" {
		b.Fail()
	}

	newId := 101

	data.Id = newId

	id = data.Id

	if id != newId {
		b.Fail()
	}

	name2 := "haha2"

	data.Name = name2

	if data.Name != name2 {
		b.Fail()
	}

}

func BenchmarkReflectRun(b *testing.B) {

	for i := 0; i < b.N; i++ {
		reflectUint(b)
	}

}

func BenchmarkOriginRun(b *testing.B) {

	for i := 0; i < b.N; i++ {
		originUint(b)
	}

}

func TestRun(t *testing.T) {

	var id int

	data := NewData()

	id = data.Get("Id").(int)

	if id != 100 {
		t.Fail()
	}

	//fmt.Println(id)

	name := data.Get("Name").(string)

	if name != "haha" {
		t.Fail()
	}

	//fmt.Println(name)

	newId := 101

	data.Set("Id", newId)

	id = data.Get("Id").(int)

	if id != newId {
		t.Fail()
	}
	//fmt.Println(data.Get("Id"))

	name2 := "haha2"

	data.Set("Name", name2)

	if data.Get("Name").(string) != name2 {
		t.Fail()
	}

	//fmt.Println(data.Get("Name"))
}
