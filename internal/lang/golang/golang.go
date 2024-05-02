package golang

type golang struct{}

func New() *golang {
	return &golang{}
}

func (g *golang) Pkg() []string {
	return []string{"# golang", "go"}
}

func (g *golang) Tools() []string {
	return []string{"gopls"}
}
