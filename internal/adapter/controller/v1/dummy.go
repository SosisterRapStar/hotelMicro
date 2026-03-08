package v1

import "net/http"

type DummyController interface {
	GetDummy(w http.ResponseWriter, r *http.Request)
}

type dummyController struct{}

func NewDummyController() *dummyController {
	return &dummyController{}
}

// GetDummy возвращает 200 OK (заглушка для проверки доступности API).
//
// @Summary  Dummy endpoint
// @Description Возвращает 200 OK. Используется для проверки доступности сервиса.
// @Tags     v1
// @Produce  json
// @Success  200
// @Router   /dummy [get]
func (d *dummyController) GetDummy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
