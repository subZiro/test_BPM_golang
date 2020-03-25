# test_BPM_golang

==========
ввод:
	echo -e 'https://golang.org\nhttps://golang.org' | go run main.go
вывод:
	Count for https://golang.org: 58
	Count for https://golang.org: 58
	Total: 116

app_with_doctor 
==========
	скрипт автоматизации записи:
	main - основная программа
	config - данные настроек (путь к chromedriver, адрес специалиста, личные данные пользователя)
	
	*для MacOS изменение прав доступа для драйвера "cd /путь/до/драйвера/ sudo chmod +x chromedriver"