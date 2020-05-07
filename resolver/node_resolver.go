package resolver

type NodeResolver struct {
	Id      string
	Url     string
	Ok      bool
	Message string
}

func (n NodeResolver) ID() string {
	return n.Id
}

func (n NodeResolver) URL() string {
	return n.Url
}

func (n NodeResolver) OK() bool {
	return n.Ok
}

func (n NodeResolver) MESSAGE() string {
	return n.Message
}
