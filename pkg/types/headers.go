package types

type HEADER_KEY string

func (hk HEADER_KEY) String() string { return string(hk) }

const (
	AuthorizationHeaderKey HEADER_KEY = "Authorization"
)
