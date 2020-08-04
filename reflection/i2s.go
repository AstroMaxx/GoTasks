package main

import (
	"fmt"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	val, ok := out.(reflect.Value)
	if !ok {
		if reflect.ValueOf(out).Kind() == reflect.Ptr {
			val = reflect.ValueOf(out).Elem()
		} else {
			return fmt.Errorf("mistake type")
		}
	}
	dat, ok := data.(reflect.Value)
	if !ok {
		dat = reflect.ValueOf(data)
	}
	typeField_data := dat.Kind()
	//fmt.Println(dat, val, val.Kind())
	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			name := val.Type().Field(i).Name
			valueField := val.Field(i)
			typeField := val.Type().Field(i).Type.Kind()
			//fmt.Println(valueField, typeField)
		LOOP:
			switch typeField_data {
			case reflect.Map:
				for _, k := range dat.MapKeys() {
					//fmt.Println(k.String(), name)
					if k.String() != name {
						continue
					}
					//fmt.Println(dat.MapIndex(k).Elem(), dat.MapIndex(k).Elem().Kind())
					switch dat.MapIndex(k).Elem().Kind() {
					case reflect.Float64:
						if typeField == reflect.Int {
							valueField.Set(reflect.ValueOf(int(dat.MapIndex(k).Elem().Float())))
							//fmt.Println(valueField, out)
							break LOOP
						} else {
							return fmt.Errorf("mistake type")
						}
					case reflect.String:
						if typeField == reflect.String {
							valueField.SetString(dat.MapIndex(k).Elem().String())
							//fmt.Println(valueField, out)
							break LOOP
						} else {
							return fmt.Errorf("mistake type")
						}
					case reflect.Bool:
						if typeField == reflect.Bool {
							valueField.SetBool(dat.MapIndex(k).Elem().Bool())
							//fmt.Println(valueField, out)
							break LOOP
						} else {
							return fmt.Errorf("mistake type")
						}
					case reflect.Map:
						if typeField == reflect.Struct {
							err := i2s(dat.MapIndex(k).Elem(), valueField)
							if err != nil {
								return err
							}
							//fmt.Println(valueField, out)
							break LOOP
						} else {
							return fmt.Errorf("mistake type")
						}
					case reflect.Slice:
						if typeField == reflect.Slice {
							err := i2s(dat.MapIndex(k).Elem(), valueField)
							if err != nil {
								return err
							}
							//fmt.Println(valueField, out)
							break LOOP
						} else {
							return fmt.Errorf("mistake type")
						}
					}
				}
			default:
				return fmt.Errorf("mistake type")
			}
		}
	case reflect.Slice:
		sl := reflect.MakeSlice(val.Type(), dat.Len(), dat.Len())
		for i := 0; i < dat.Len(); i++ {
			err := i2s(dat.Index(i).Elem(), sl.Index(i))
			if err != nil {
				return err
			}
		}
		val.Set(sl)
	}
	return nil
}