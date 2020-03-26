package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
    k int = 5 //кол-во "потоков"
    q string = "go" //строка поиска
)

func init() {
	//Задаем правила разбора:
	flag.IntVar(&k, "k", k, "максимальное количество потоков")
	flag.StringVar(&q, "q", q, "что ищем")
}

type ProducerObject struct {
	// структура парсинга
	Counter *counter
	Q string
	StackChannel chan int
}

type counter struct {
	// структура счетчика
	count int
	lock *sync.Mutex
}

func new_counter() *counter  {
	//создание экз счетчика
	return &counter{lock:&sync.Mutex{}}  
}

func (c *counter) count_add(count int) {
	c.lock.Lock()  // блокироква
	defer c.lock.Unlock()
	c.count += count  // увеличенрие найденых вхождений
}

func (c *counter) count_get() int {
	c.lock.Lock()  // блокироква
	defer c.lock.Unlock()
	return c.count  // возврат общего количества найденых вхождений
}

func count_q_inbody(string_html, q string) int {
	// подсчет количества вхождений строки в полученый запрос
	return strings.Count(strings.ToLower(string_html), q)
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

func get_url(url string) (*string, error) {
	// get запрос на URL, возвращает содержимое htlm страницы
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer response.Body.Close()  // закрытие
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	result := string(contents)
	return &result, nil
}


func parse(urls []string, q string, k int,  wg *sync.WaitGroup) {
	// парсинг всех url через горутины передача через канал
	counter := new_counter()
	for _, url := range urls {
		wg.Add(1)  // инкриментирование стека горутин
		go parse_one_url(url, q, wg)
	}
	wg.Wait()   // ожидание закрытия всех потоков
	fmt.Printf("Total: %d\n", counter.count_get())
}

func parse_one_url(url string, q string, wg *sync.WaitGroup)  {
	// парсинг одного url
	defer wg.Done()  //вычитаем стек потоков
	contents, err := get_url(url)
	if err != nil {
		return
	}
	count_i := count_q_inbody(*contents, q)
	counter.count_add(count_i)
	fmt.Printf("Count '%s' for %s : %d\n", q, url, count_i)
}


func main() {
	flag.Parse() //И запускаем разбор аргументов
	fmt.Println("========================")
	urls := scaner_urls()
	fmt.Printf("Получено %d url\n", len(urls))
	wg := new(sync.WaitGroup)  // индикатор завершения потоков
	parse(urls, q, k, wg)
	fmt.Println("========================")
}