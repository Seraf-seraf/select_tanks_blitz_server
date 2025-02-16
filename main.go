package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/eiannone/keyboard"
)

const hostsFilePath = "C:\\Windows\\System32\\drivers\\etc\\hosts"

var servers = map[int]string{
	0: "login0.tanksblitz.ru",
	1: "login1.tanksblitz.ru",
	2: "login2.tanksblitz.ru",
	3: "login3.tanksblitz.ru",
	4: "login4.tanksblitz.ru",
	5: "Сбросить настройки",
}

func main() {
	defer waitForKeyPress()
	printMenu()

	var choice int
	fmt.Println("Введите номер настройки: ")
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 0 || choice > 5 {
		fmt.Println("Некорректный выбор")
		return
	}

	hosts, err := os.Open(hostsFilePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer hosts.Close()

	var lines []string
	scanner := bufio.NewScanner(hosts)
	for scanner.Scan() {
		line := scanner.Text()
		shouldRemove := false

		for _, server := range servers {
			if strings.Contains(line, server) {
				shouldRemove = true
				break
			}
		}

		if !shouldRemove {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Ошибка чтения файла:", err)
		return
	}

	if choice == 5 {
		writeToFile(lines)
		fmt.Println("Настройки сброшены")
		return
	}

	for key, server := range servers {
		if key == choice {
			lines = append(lines, "# "+server)
			continue
		}

		if key != 5 {
			lines = append(lines, "127.0.0.1 "+server)
		}
	}

	writeToFile(lines)
	fmt.Printf("Выбран сервер %v\n", choice)
	fmt.Println("Подключение к выбранному серверу может занять время")
}

func printMenu() {
	fmt.Println("Выберите из меню:")

	keys := make([]int, 0, len(servers))
	for key := range servers {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, key := range keys {
		fmt.Printf("Сервер %d: %s\n", key, servers[key])
	}
}

func writeToFile(lines []string) {
	hosts, err := os.OpenFile(hostsFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer hosts.Close()

	writer := bufio.NewWriter(hosts)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}
	writer.Flush()
}

func waitForKeyPress() {
	fmt.Println("Нажмите любую клавишу для завершения...")

	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	_, _, err = keyboard.GetKey()
	if err != nil {
		log.Fatal(err)
	}
}
