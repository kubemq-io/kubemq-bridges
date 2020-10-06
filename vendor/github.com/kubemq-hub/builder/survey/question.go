package survey

type Question interface {
	Complete() error
	Render(target interface{}) error
}
