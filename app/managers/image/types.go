package managers

type Manager interface {
	IsURLAllowed(imageURL string) bool
	ProcessImage(imageURL string, width, height int, format string) (string, error)
}
