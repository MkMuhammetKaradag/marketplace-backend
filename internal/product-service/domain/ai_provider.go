package domain

type AiProvider interface {
	GetVector(text string) ([]float32, error)
}
