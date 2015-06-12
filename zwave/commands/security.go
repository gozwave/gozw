package commands

func NewSecurityNonceGet() []byte {
	return []byte{
		CommandClassSecurity,
		SecurityNonceGet,
	}
}
