package managers

type ImageManager struct {
}

func NewManager() *ImageManager {
	return &ImageManager{}
}

func (m *ImageManager) IsURLAllowed(imageURL string) bool {
	return true
}

func (m *ImageManager) ProcessImage(imageURL string, width, height int, format string) (string, error) {
	return "", nil
}
