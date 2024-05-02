package python

type python struct{}

func New() *python {
	return &python{}
}

func (p *python) Pkg() []string {
	return []string{"# python", "python3"}
}

func (p *python) Tools() []string {
	return []string{"nodePackages.pyright"}
}
