package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"github.com/go-ping/ping"
)

type dockerContainer struct {
	Id     string
	Image  string
	Status string
	Ip     string
}

type PingResult struct {
	IP       string `json:"ip"`        // ip-адрес контейнера
	PingTime int    `json:"ping_time"` // продолжительность пинга в миллисекундах
	Success  bool   `json:"success"`   // успешен ли pingfunc
}

func pingfunc(ip string) (int, bool) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Printf("Ошибка при создании Pinger: %s\n", err)
		return 0, false
	}

	pinger.Count = 7
	pinger.Timeout = time.Second * 5
	err = pinger.Run()
	if err != nil {
		log.Printf("Ошибка при выполнении пинга: %s\n", err)
		return 0, false
	}

	stats := pinger.Statistics()
	return int(stats.AvgRtt.Milliseconds()), stats.PacketsRecv > 0
}

func takeDockerContainers() []dockerContainer {
	cli, err := client.NewClientWithOpts(
		client.WithHost("tcp://host.docker.internal:2375"), // Используем TCP
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatalf("Ошибка при создании Docker-клиента: %s", err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		log.Fatalf("Ошибка при получении списка контейнеров: %s", err)
	}

	var dockerContainers []dockerContainer
	for _, c := range containers {
		containerDetails, err := cli.ContainerInspect(context.Background(), c.ID)
		if err != nil {
			log.Printf("Ошибка при получении деталей контейнера %s: %s", c.ID[:10], err)
			continue
		}

		ipAddress := containerDetails.NetworkSettings.IPAddress

		if ipAddress == "" {
			for _, network := range containerDetails.NetworkSettings.Networks {
				ipAddress = network.IPAddress
				break
			}
		}

		dockerContainers = append(dockerContainers, dockerContainer{
			Id:     c.ID,
			Image:  c.Image,
			Status: c.Status,
			Ip:     ipAddress,
		})
	}

	return dockerContainers
}

func sendPingResult(result PingResult) {
	postBody, _ := json.Marshal(result)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("http://backend:8080/v1/containers/%s", result.IP),
		"application/json",
		responseBody,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}

func uniteConteiners() {
	ctx := context.Background()

	// Создаем клиент Docker
	cli, err := client.NewClientWithOpts(
		client.WithHost("tcp://host.docker.internal:2375"), // Используем TCP
		client.WithAPIVersionNegotiation(),
	)

	// Получаем список всех контейнеров
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		log.Fatalf("Ошибка при получении списка контейнеров: %s", err)
	}

	// Создаем новую сеть
	networkName := "my_network"
	network, err := cli.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{})
	if err == nil {
		log.Printf("Сеть %s уже существует. Подключаем контейнеры к ней.\n", networkName)
	} else {
		// Если сети нет, создаем её
		_, err = cli.NetworkCreate(ctx, networkName, types.NetworkCreate{})
		if err != nil {
			log.Fatalf("Ошибка создания сети: %v", err)
		}
		log.Printf("Сеть %s была создана.\n", networkName)
	}

	// Подключаем все контейнеры к созданной сети
	for _, c := range containers {
		// Проверяем подключение
		isConnected := false
		for k := range network.Containers {
			if k == c.ID {
				isConnected = true
				break
			}
		}

		if !isConnected {
			// Подключаем контейнер к сети
			err = cli.NetworkConnect(ctx, networkName, c.ID, nil)
			if err != nil {
				log.Printf("Ошибка подключения контейнера %s к сети %s: %v", c.ID, networkName, err)
			} else {
				log.Printf("Контейнер %s подключен к сети %s\n", c.ID, networkName)
			}
		} else {
			log.Printf("Контейнер %s уже подключен к сети %s\n", c.ID, networkName)
		}
	}

	log.Println("Все контейнеры были обработаны.")
}

func main() {
	for {
		uniteConteiners()
		dockerContainers := takeDockerContainers()

		for _, ip := range dockerContainers {
			pingTime, success := pingfunc(ip.Ip) // Пингуем IP-адрес
			sendPingResult(PingResult{
				IP:       ip.Ip,
				PingTime: pingTime,
				Success:  success,
			})

			log.Printf("IP: %s, PingTime: %d Success: %t\n", ip.Ip, pingTime, success)
			fmt.Println()
		}

		time.Sleep(10 * time.Second)
	}

}
