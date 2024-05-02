package clang

type clang struct{}

func New() *clang {
	return &clang{}
}

func (c *clang) Pkg() []string {
	return []string{}
}

func (c *clang) Tools() []string {
	return []string{"# clang", "cmake", "ninja", "clang-tools"}
}
