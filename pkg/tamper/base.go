package tamper

type Tamper interface {
	Name() string
	Desc() string
	Exec(payload string) string
}
