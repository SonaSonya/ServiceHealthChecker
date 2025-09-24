package colors

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

func Wrap(text, color string) string {
	return color + text + Reset
}
