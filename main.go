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

func count_strinbody(string_html, q_str string) int {
	// подсчет количества вхождений строки в полученый запрос
	return strings.Count(strings.ToLower(string_html), q_str)
}

func get_url(url string) []byte {
	// get запрос на URL
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
	return contents
}

func scaner_urls() []string {
	// полкчение массива ссылок передаваемых при запуске программы
	a := []string{}
	buf := bufio.NewScanner(os.Stdin)
	for {
		url := strings.TrimSpace(buf.Text())
		a = append(a, url)  //добавление линка в массив
		//fmt.Println(url)
		if !buf.Scan() {
			break
		}
	}
	return a[1:]
}

func parse(urls []string, q_str string) int {
	//  парс
	total := 0
	for _, url := range urls {
		fmt.Println(url, q_str)
		total += 1
	}

	return total
}


func main() {
	flag.Parse() //И запускаем разбор аргументов
	//b := count_strinbody(string(contents), q_str)
	//fmt.Printf("Найдено %d совпадений\n", b)

	fmt.Println("========================")
	urls := scaner_urls()
	//s := scaner_urls()
	fmt.Printf("получили %d url\n", len(urls))
	fmt.Printf("array urls: %s\n", urls)
	fmt.Println("========================")

	result := parse(urls, q_str)
	fmt.Println("result = ", result)
	fmt.Println("========================")
	//input := "foo\nbar\nbaz"
    //scanner := bufio.NewScanner(os.Stdin)
    //scanner.Split(bufio.ScanLines)
    //for scanner.Scan() {
    //   fmt.Println(scanner.Text())
    //}







	fmt.Println("========================")
}






