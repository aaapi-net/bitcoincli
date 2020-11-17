package bitcoincli

const (
	RPC_REQUEST_TIMEOUT = 30
)

type BitcoinCliConfig struct {
	IsTest bool
	Host string
	Port int
	User string
	Password string
	UseSsl bool
	Timeout int
	WalletNotifyConfig *WalletNotifyConfig
}

type WalletNotifyConfig struct {
	Host string
	Port int
}

func NewDefaultBitcoinCliConfig() *BitcoinCliConfig {
	return &BitcoinCliConfig{IsTest: true, Host: "127.0.0.1", Port: 17332, UseSsl: false, Timeout: RPC_REQUEST_TIMEOUT}
}

func (c *BitcoinCliConfig) WithUser(user string) *BitcoinCliConfig {
	c.User = user
	return c
}

func (c *BitcoinCliConfig) WithPassword(password string) *BitcoinCliConfig {
	c.Password = password
	return c
}

func (c *BitcoinCliConfig) WithTimeout(timeout int) *BitcoinCliConfig {
	c.Timeout = timeout
	return c
}

func (c *BitcoinCliConfig) WithWalletNotify(config *WalletNotifyConfig) *BitcoinCliConfig {
	c.WalletNotifyConfig = config
	return c
}