package types

type CONTEXT_KEY string

func (ck CONTEXT_KEY) String() string { return string(ck) }

const (
	CredentialDataContextKey CONTEXT_KEY = "credential_data"
)
