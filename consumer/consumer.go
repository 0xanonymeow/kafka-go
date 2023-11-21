package consumer

type Topic struct {
	Name string
}

type Consumer struct {
	Type   interface{}
	Topics []Topic
}
