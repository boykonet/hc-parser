package printer

type textColors struct {
	Reset, Red, Green, Yellow, Blue, Cyan string
}

func NewColors() *textColors {
	return &textColors{
		Reset: "\033[0m",
		Red: "\033[31m",
		Green: "\033[32m",
		Yellow: "\033[33m",
		Blue: "\033[34m",
		Cyan: "\033[36m",
	}
}