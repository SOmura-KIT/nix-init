package latex

type latex struct{}

func New() *latex {
	return &latex{}
}

func (l *latex) Pkg() []string {
	return []string{"# latex", "texlive-full"}
}

func (l *latex) Tools() []string {
	return []string{"texlab"}
}
