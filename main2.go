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
	return &counter{lock:&sync.Mutex{}}  //создание экз счетчика
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

func scaner_urls() chan string{
	// получение массива ссылок передаваемых при запуске программы
	//a := []string{}
	chan_url := make(chan string)  // создание канала для передачи сканированных url
	buf := bufio.NewScanner(os.Stdin)
	go func() {
		var wg_url sync.WaitGroup
		for {
			url := strings.TrimSpace(buf.Text())
			if url != ""{
				wg_url.Add(1)
				chan_url <- url  // передача url в канал
				//a = append(a, url)  //добавление линка в массив
			}
			if !buf.Scan() {
				break
			}
		}
		wg_url.Wait()
		close(chan_url)
	}()
	return chan_url
	//return a[1:]
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

func parse(urls chan string, q string, k int, wg *sync.WaitGroup) {
	// парсинг всех url через горутины передача через канал
	counter := new_counter()
	stackChan := make(chan int, k)
	for url := range urls {
		wg.Add(1)  // инкриментирование стека горутин
		go parse_one_url(url, ProducerObject{StackChannel:stackChan, Counter:counter, Q:q}, wg)
	}
	close(stackChan)
	wg.Wait()   // ожидание закрытия всех потоков
	fmt.Printf("Total: %d\n", counter.count_get())
}

func parse_one_url(url string, pobj ProducerObject, wg *sync.WaitGroup) {
	// парсинг одного url
	defer wg.Done()  //вычитаем стек потоков
	contents, err := get_url(url)
	if err != nil {
		return
	}
	count_i := count_q_inbody(*contents, pobj.Q)  // количество вхождений
	pobj.Counter.count_add(count_i)  // добавление количества вхождений в счетчик
	fmt.Printf("Count '%s' for %s : %d\n", q, url, count_i)
}

func main() {
	fmt.Println("========================")
	flag.Parse()  //запускаем разбор аргументов
	//urls := scaner_urls()
	wg := new(sync.WaitGroup)  // индикатор завершения потоков
	parse(scaner_urls(), q, k, wg)
	fmt.Println("========================")
}