package deno

type deno struct{}

func New() *deno {
	return &deno{}
}

func (d *deno) Pkg() []string {
	return []string{"# deno", "deno"}
}

func (d *deno) Tools() []string {
	return []string{}
}
