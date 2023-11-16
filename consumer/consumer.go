package consumer

type Topic struct {
	Topic    string
	Name    string
	Handler interface{}
}