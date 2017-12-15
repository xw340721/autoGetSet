package autosetget

import (
	"reflect"
	"strings"
)

const (
	FIELD_NONE = iota
	FIELD_PRIVATE
	FIELD_PUBLIC
)

type IData interface {
	Get(key string) interface{}
	Set(key string, data interface{})
}

type DataImpl struct {
	fieldPrivilege map[string]bool
	privateSetFn   map[string]reflect.Value
	privateGetFn   map[string]reflect.Value
	obj            interface{}
}

func NewDataImpl() *DataImpl {
	return &DataImpl{
		fieldPrivilege: make(map[string]bool),
		privateSetFn:   make(map[string]reflect.Value),
		privateGetFn:   make(map[string]reflect.Value),
	}
}

func (p *DataImpl) setFieldPrivilege(key string, bool bool) {
	if _, ok := p.fieldPrivilege[key]; !ok {
		p.fieldPrivilege[key] = bool
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
		if _, ok := p.privateSetFn[key]; !ok {
			p.privateSetFn[key] = fn
		}

	}
}

func (p *DataImpl) setPrivateGetFn(key string, fn reflect.Value) {
	privilegeAttr := p.GetPrivateAttr(key)

	if privilegeAttr == FIELD_PRIVATE {
		if _, ok := p.privateGetFn[key]; !ok {
			p.privateGetFn[key] = fn
		}

	}
}

func (p *DataImpl) Get(key string) interface{} {

	privilegeAttr := p.GetPrivateAttr(key)

	if privilegeAttr == FIELD_PUBLIC {

		rv := reflect.ValueOf(p.obj)

		value := rv.FieldByName(key)

		return value.Interface()

	} else if privilegeAttr == FIELD_PRIVATE {

		if getFn, ok := p.privateGetFn[key]; ok {

			in := []reflect.Value{}

			reValue := getFn.Call(in)

			return reValue[0].Interface()
		}
	}

	return nil
}

func (p *DataImpl) Set(key string, data interface{}) {

	privilegeAttr := p.GetPrivateAttr(key)

	if privilegeAttr == FIELD_PUBLIC {

		rv := reflect.ValueOf(p.obj)

		value := rv.FieldByName(key)

		vdata := reflect.ValueOf(data)

		value.Set(vdata)
	} else if privilegeAttr == FIELD_PRIVATE {
		if setFn, ok := p.privateSetFn[key]; ok {

			rv := reflect.ValueOf(data)

			in := []reflect.Value{
				rv,
			}
			setFn.Call(in)
		}
	}
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

func (p *DataImpl) SetField(obj interface{}) {

	p.obj = obj

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
			setFnName := "set" + strings.ToUpper(string(name[0])) + name[1:]
			setFn := rv.MethodByName(setFnName)
			if setFn.Pointer() != 0 {
				p.setPrivateSetFn(name, setFn)
			}

			// register get

			getFnName := "get" + strings.ToUpper(string(name[0])) + name[1:]

			getFn := rv.MethodByName(getFnName)
			if getFn.Pointer() != 0 {
				p.setPrivateGetFn(name, getFn)
			}

		}

	}

}

//func main() {
//
//	// HelloWord()
//	fmt.Print(HelloWord)
//}
