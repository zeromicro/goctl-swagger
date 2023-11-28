package v2

import "github.com/zeromicro/go-zero/tools/goctl/api/spec"

func asDefineStruct(obj spec.Type) (spec.DefineStruct, bool) {
	stru, ok := obj.(spec.DefineStruct)
	if ok {
		return stru, ok
	}
	if ptr, ok := obj.(spec.PointerType); ok {
		return asDefineStruct(ptr.Type)
	}
	return spec.DefineStruct{}, false
}

func asArrayType(obj spec.Type) (spec.ArrayType, bool) {
	array, ok := obj.(spec.ArrayType)
	if ok {
		return array, ok
	}
	if ptr, ok := obj.(spec.PointerType); ok {
		return asArrayType(ptr.Type)
	}
	return spec.ArrayType{}, false
}
