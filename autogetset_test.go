package autosetget_test

import(
 "testing"
	. "github.com/xw340721/autoGetSet"
)

type Data struct {
	*DataImpl
	Id   int    `privilege:"public"`
	Name string `privilege:"private"`
}

func (p *Data) setName(name string) {
	p.Name = name
}

func (p *Data) getName() string {
	return p.Name
}


func NewData() *Data {
	
	data := &Data{
		DataImpl: NewDataImpl(),
		Id:       100,
		Name:     "haha",
	}
	
	data.SetField(data)
	
	return data
}

func TestRun(t *testing.T)  {


	var id int

	data := NewData()
	
	
	id =data.Get("Id").(int)
	
	if id!=100{
		t.Fail()
	}
	
	newId := 101
	
	data.Set("Id",newId)
	
	id = data.Get("Id").(int)
	
	if id!=newId{
		t.Fail()
	}
	

}