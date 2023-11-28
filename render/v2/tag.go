package v2

import (
	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func lookupGozeroTag(tags []*spec.Tag) *spec.Tag {
	for _, tag := range tags {
		switch tag.Key {
		case types.PathTagKey, types.FormTagKey, types.HeaderTagKey, types.JsonTagKey:
			return tag
		}
	}
	return nil
}

func isOptionalTag(tag *spec.Tag) bool {
	if len(tag.Options) == 0 {
		return false
	}
	return stringx.Contains(tag.Options, "optional") || stringx.Contains(tag.Options, "omitempty")
}
