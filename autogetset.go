package autosetget

import (
	"reflect"
	"strings"
	"fmt"
)


import (
"fmt"
"reflect"
"strings"
)

const (
	FIELD_NONE = iota
	FIELD_PRIVATE
	FIELD_PUBLIC
)

type IData interface {
	get(key string) interface{}
	set(key string, data ...interface{})
}

type DataImpl struct {
	fieldPrivilege map[string]bool
	privateSetFn   map[string]reflect.Value
	privateGetFn   map[string]reflect.Value
}

func NewDataImpl() *DataImpl {
	return &DataImpl{
		fieldPrivilege: make(map[string]bool),
		privateSetFn:   make(map[string]reflect.Value),
		privateGetFn:   make(map[string]reflect.Value),
	}
}

func (p *DataImpl) setFieldPrivilege(key string, bool bool) {
	if fieldPrivilege, ok := p.fieldPrivilege[key]; !ok {
		fieldPrivilege = bool
	}
}

func (p *DataImpl) GetPrivateAttr(key string) int {

	if fieldPrivilege, ok := p.fieldPrivilege[key]; ok {

		if fieldPrivilege {
			return FIELD_PUBLIC
		} else {
			return FIELD_PRIVATE
		}
	}

	return FIELD_NONE
}

func (p *DataImpl) setPrivateSetFn(key string, fn reflect.Value) {
	privilegeAttr := p.GetPrivateAttr(key)

	if privilegeAttr == FIELD_PRIVATE {
		if privilegeField, ok := p.privateSetFn[key]; !ok {
			privilegeField = fn
		}

	}
}

func (p *DataImpl) setPrivateGetFn(key string, fn reflect.Value) {
	privilegeAttr := p.GetPrivateAttr(key)

	if privilegeAttr == FIELD_PRIVATE {
		if privilegeField, ok := p.privateGetFn[key]; !ok {
			privilegeField = fn
		}

	}
}

func (p *DataImpl) get(key string) interface{} {

	privilegeAttr := p.GetPrivateAttr(key)

	if privilegeAttr == FIELD_PUBLIC {

	} else if privilegeAttr == FIELD_PRIVATE {

		if getFn, ok := p.privateGetFn[key]; ok {

			in := []reflect.Value{}

			reValue := getFn.Call(in)

		}

	}

	return nil
}

func (p *DataImpl) set(key string, data ...interface{}) {

}

func (p *DataImpl) SetField(obj interface{}) {

	rt := reflect.TypeOf(obj)

	rv := reflect.ValueOf(obj)

	for i := 0; i < rt.NumField(); i++ {
		rfd := rt.Field(i)

		name := rfd.Name
		tag := rfd.Tag.Get("privilege")

		if tag == "" {
			continue
		}

		privilege := Privilege(tag)
		p.setFieldPrivilege(name, privilege)

		if !privilege {

			// register  set
			setFnName := "set" + strings.ToUpper(name[0]) + name[1:]
			setFn := rv.MethodByName(setFnName)
			if setFn.Pointer() != 0 {
				p.setPrivateSetFn(name, setFn)
			}

			// register get

			getFnName := "get" + strings.ToUpper(name[0]) + name[1:]

			getFn := rv.MethodByName(getFnName)
			if getFn.Pointer() != 0 {
				p.setPrivateGetFn(name, getFn)
			}

		}

	}

}

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

func Privilege(tag string) bool {
	switch tag {
	case "public":
		return true
	case "private":
		return false
	default:
		return false
	}
}

func (p *Data) set(key string, data ...interface{}) {

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

//func main() {
//
//	// HelloWord()
//	fmt.Print(HelloWord)
//}
