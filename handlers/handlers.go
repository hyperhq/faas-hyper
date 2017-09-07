package handlers

import "github.com/alexellis/faas-hyper/hyper"

type Handler struct {
	*hyper.Hyper
}

func New() (*Handler, error) {
	_hyper, err := hyper.New()
	if err != nil {
		return nil, err
	}
	return &Handler{_hyper}, nil
}
