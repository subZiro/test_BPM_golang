package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
    k int = 5 //кол-во "потоков"
    q_str string = "go" //строка поиска
)

func init() {
	//Задаем правила разбора:
	flag.IntVar(&k, "k", k, "максимальное количество потоков")
	flag.StringVar(&q_str, "q", q_str, "что ищем")
}

func count_q_inbody(string_html, q string) int {
	// подсчет количества вхождений строки в полученый запрос
	return strings.Count(strings.ToLower(string_html), q)
}

func get_url(url string) (*string, error) {
	// get запрос на URL, возвращает содержимое htlm страницы
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer response.Body.Close()  // закрытие
	//s, err := io.Copy(os.Stdout, resp.Body)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	result := string(contents)
	return &result, nil
}

func scaner_urls() []string {
	// полкчение массива ссылок передаваемых при запуске программы
	a := []string{}
	buf := bufio.NewScanner(os.Stdin)
	for {
		url := strings.TrimSpace(buf.Text())
		a = append(a, url)  //добавление линка в массив
		if !buf.Scan() {
			break
		}
	}
	return a[1:]
}

func parse(urls []string, q string) int {
	//  парс
	total := 0
	for _, url := range urls {
		contents, err := get_url(url)
		if err != nil {
			return 0
		}
		count_i := count_q_inbody(*contents, q)
		fmt.Printf("Count '%s' for %s : %d\n", q, url, count_i)
		
		total += count_i
	}
	return total
}


func main() {
	flag.Parse() //И запускаем разбор аргументов
	fmt.Println("========================")
	urls := scaner_urls()
	fmt.Printf("получили %d url\n", len(urls))
	fmt.Printf("array urls: %s\n", urls)
	fmt.Println("========================")

	result := parse(urls, q_str)
	fmt.Println("result = ", result)
	fmt.Println("========================")
}