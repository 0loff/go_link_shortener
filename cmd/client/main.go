package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	endpoint := "http://127.0.0.1:8080"
	// Контейнер данных для запроса
	data := url.Values{}
	// Приглашение в консоли
	fmt.Println("Введите длинный URL")
	// Открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// Читаем строку из консоли
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\n")
	// Заполняем контейнер данными
	data.Set("url", long)
	// Добавляем HTTP клиент
	client := &http.Client{}
	// Пишем запрос
	// Запрос методом POST должен, помимо заголовков, содержать тело
	// Тело должно быть источником потокового чтения io.Reader
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	// В заголовках запроса указываем кодировку
	request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// Выводим код ответа
	fmt.Println("Статус код ", response.Status)
	defer response.Body.Close()

	// Читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// И печатаем его
	fmt.Println(string(body))
}
