package agent_handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	models "github.com/maliven1/metrics/internal/model"
)

type Agent interface {
	GetMetrics() (map[string]float64, map[string]int64)
	CollectMetrics()
}

type SendClient struct {
	AddHandler Agent
}

func NewSendClient(s Agent) *SendClient {
	return &SendClient{AddHandler: s}
}

func (s SendClient) SendClientMetrics() {
	fmt.Println("Start")
	endpoint := "http://localhost:8080/update/"
	client := &http.Client{}
	go s.AddHandler.CollectMetrics()
	for {

		time.Sleep(time.Duration(models.ReportInterval) * time.Second)
		Gauge, Counter := s.AddHandler.GetMetrics()
		fmt.Println("Start2")
		for i, v := range Gauge {
			if i == "" {
				continue
			}
			// контейнер данных для запроса
			data, _ := url.JoinPath(models.Gauge, i, fmt.Sprint(v))

			// добавляем HTTP-клиент

			// пишем запрос
			// запрос методом POST должен, помимо заголовков, содержать тело
			// тело должно быть источником потокового чтения io.Reader
			request, err := http.NewRequest(http.MethodPost, endpoint+data, nil)
			if err != nil {
				panic(err)
			}
			// в заголовках запроса указываем кодировку
			request.Header.Add("Content-Type", "Content-Type: text/plain")
			// отправляем запрос и получаем ответ
			response, err := client.Do(request)
			if err != nil {
				panic(err)
			}
			fmt.Println(request)
			fmt.Println("Статус-код ", response.Status, i, v)
			// defer response.Body.Close()
		}
		for i, v := range Counter {
			if i == "" {
				continue
			}
			// контейнер данных для запроса
			data, _ := url.JoinPath(models.Counter, i, fmt.Sprint(v))

			// добавляем HTTP-клиент

			// пишем запрос
			// запрос методом POST должен, помимо заголовков, содержать тело
			// тело должно быть источником потокового чтения io.Reader
			request, err := http.NewRequest(http.MethodPost, endpoint+data, nil)
			if err != nil {
				panic(err)
			}
			// в заголовках запроса указываем кодировку
			request.Header.Add("Content-Type", "Content-Type: text/plain")
			// отправляем запрос и получаем ответ
			response, err := client.Do(request)
			if err != nil {
				panic(err)
			}
			fmt.Println(request)
			fmt.Println("Статус-код ", response.Status, i, v)
			// defer response.Body.Close()
		}

		// выводим код ответа

		// читаем поток из тела ответа
		// body, err := io.ReadAll(response.Body)
		// if err != nil {
		// 	panic(err)
		// }
		// и печатаем его
		// fmt.Println(string(body))
	}
}
