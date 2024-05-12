package preview

import (
	"net/url"

	"backend/type/payload"
	"backend/type/response"
	"backend/type/table"
	"backend/util/value"
)

func CreateInfo(component *table.SystemComponent, group *table.SystemGroup) (*payload.ComponentInfo, error) {
	path, err := url.JoinPath(*group.Address, *component.Name)
	if err != nil {
		return nil, response.Error(true, "Unable to join path", err)
	}
	var queryInfo *payload.QueryInfo
	if component.QueryKey != nil {
		queryInfo = &payload.QueryInfo{
			Key:        component.QueryKey,
			ValueStart: component.QueryValStart,
			ValueEnd:   component.QueryValEnd,
		}
	}
	componentInfo := &payload.ComponentInfo{
		Url:       value.Ptr(path),
		QueryInfo: queryInfo,
	}

	return componentInfo, nil
}
