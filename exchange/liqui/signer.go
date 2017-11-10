package liqui

type Signer interface {
	GetLiquiKey() string
	LiquiSign(msg string) string
}
