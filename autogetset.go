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
	rv             reflect.Value
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

		value := p.rv.Elem().FieldByName(key)

		return value.Interface()

	} else if privilegeAttr == FIELD_PRIVATE {

		if getFn, ok := p.privateGetFn[key]; ok {

			in := []reflect.Value{}

			//fmt.Println(getFn.String())

			reValue := getFn.Call(in)

			return reValue[0].Interface()
		}
	}

	return nil
}

func (p *DataImpl) Set(key string, data interface{}) {

	privilegeAttr := p.GetPrivateAttr(key)

	if privilegeAttr == FIELD_PUBLIC {

		value := p.rv.Elem().FieldByName(key)

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

func (p *DataImpl) SetField(rt reflect.Type, rv reflect.Value) {
	p.rv = rv

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
			setFnName := "Set" + strings.ToUpper(string(name[0])) + name[1:]
			setFn := rv.MethodByName(setFnName)

			//fmt.Printf("register set func name %s\r\n func to string %s \r\n", setFnName, setFn.String())

			if setFn.IsValid() {
				p.setPrivateSetFn(name, setFn)
			}

			// register get
			getFnName := "Get" + strings.ToUpper(string(name[0])) + name[1:]

			getFn := rv.MethodByName(getFnName)

			//fmt.Printf("register get func name %s\r\n func to string %s \r\n", getFnName, setFn.String())

			if getFn.IsValid() {
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
