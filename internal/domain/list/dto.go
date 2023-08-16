package list

import (
	"errors"
	"net/http"
	"time"
)

type Request struct {
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
	Status   bool   `json:"-"`
}

func (s *Request) Bind(r *http.Request) error {
	if s.Title == "" {
		return errors.New("title: cannot be blank")
	}
	if len(s.Title) >= 200 {
		return errors.New("title: не более 200 символов")
	}

	if s.ActiveAt == "" {
		return errors.New("activeAt: cannot be blank")

	}

	// Проверка валидности даты
	layout := "2006-01-02"
	_, err := time.Parse(layout, s.ActiveAt)
	if err != nil {
		return errors.New("activeAt: неверный формат даты")
	}

	return nil
}

type Response struct {
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		Title:    *data.Title,
		ActiveAt: *data.ActiveAt,
	}
	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
