package payload

type StateResponse struct {
	SideA map[string]*ComponentInfo `json:"sideA"`
	SideB map[string]*ComponentInfo `json:"sideB"`
}

type ComponentInfo struct {
	Url       *string    `json:"url"`
	QueryInfo *QueryInfo `json:"queryInfo"`
}

type QueryInfo struct {
	Key        *string `json:"key"`
	ValueStart *int64  `json:"valueStart"`
	ValueEnd   *int64  `json:"valueEnd"`
}
