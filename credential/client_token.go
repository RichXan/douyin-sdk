package credential

// ClientTokenHandle client token获取
type ClientTokenHandle interface {
	// GetClientToken 获取client token
	GetClientToken() (client_token string, err error)
}
