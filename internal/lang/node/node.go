package node

type node struct{}

func New() *node {
	return &node{}
}

func (n *node) Pkg() []string {
	return []string{"# node", "nodejs"}
}

func (n *node) Tools() []string {
	return []string{"yarn"}
}
