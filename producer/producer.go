package producer

type Topic struct {
	Key  string
	Name string
}

type Producer struct {
	Topics []Topic
}
