package transform

import (
	"errors"
	"reflect"
	"unsafe"
)

func Copy(s interface{}, t interface{}) interface{} {
	sValue := reflect.ValueOf(s)
	tValue := reflect.ValueOf(t)
	if s, e := lsRequire(sValue, tValue); s == false {
		panic(e)
	}
	if sValue.Kind() == reflect.Ptr {
		if s, e := lsRequire(sValue.Elem(), tValue.Elem()); s == false {
			panic(e)
		}
		return transform(reflect.ValueOf(s).Elem(), reflect.ValueOf(t).Elem(), reflect.TypeOf(s).Elem(), reflect.TypeOf(t).Elem())
	}
	return transform(reflect.ValueOf(&s).Elem(), reflect.ValueOf(&t).Elem(), reflect.TypeOf(&s).Elem(), reflect.TypeOf(&t).Elem())
}

func transform(sValue reflect.Value, tValue reflect.Value, sType reflect.Type, tType reflect.Type) interface{} {
	if s, e := lsRequire(sValue, tValue); s == false {
		panic(e)
	}
	// 结构体
	if sValue.Kind() == reflect.Struct {
		targetMap := make(map[string]reflect.Type)
		// key：字段名；value：字段value
		for i := 0; i < tType.NumField(); i++ {
			targetMap[tType.Field(i).Name] = tType.Field(i).Type
		}
		for tFieldName, tFieldType := range targetMap {
			sField, find := sType.FieldByName(tFieldName)
			// 源中没有该字段则不做处理
			if find == false {
				continue
			}
			if lsSet(sField.Type, tFieldType) {
				setValue(sValue.FieldByName(tFieldName), tValue.FieldByName(tFieldName))
			} else {
				// 递归
				tField, _ := tType.FieldByName(tFieldName)
				transform(sValue.FieldByName(tFieldName), tValue.FieldByName(tFieldName), sField.Type, tField.Type)
			}
		}
	} else {
		// array or slice
	}
}

func lsSet(sType reflect.Type, tType reflect.Type) bool {
	if sType.Name() == tType.Name() {
		return true
	}
	return false
}

func lsRequire(sValue reflect.Value, tValue reflect.Value) (bool, error) {
	if sValue.Kind() != tValue.Kind() {
		return false, errors.New("Both kind is different ")
	}
	if sValue.Kind() == reflect.Map {
		return false, errors.New("kind is map, gtrans not support map")
	}
	return true, nil
}

func setValue(sValue reflect.Value, tValue reflect.Value) {
	reflect.NewAt(tValue.Type(), unsafe.Pointer(tValue.UnsafeAddr())).Elem().Set(sValue)
}
