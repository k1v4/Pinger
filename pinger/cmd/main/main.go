package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/docker/docker/client"
)

type dockerContainer struct {
	Id     string
	Image  string
	Status string
	Ip     string
}

type PingResult struct {
	IP             string    `json:"ip"`            // ip-адрес контейнера
	PingTime       int       `json:"ping_time"`     // продолжительность пинга в миллисекундах
	Success        bool      `json:"is_successful"` // успешен ли pingFunc
	LastSuccessful time.Time `json:"last_successful"`
}

//func pingFunc(ip string) (int, bool) {
//	pinger, err := ping.NewPinger(ip)
//	if err != nil {
//		log.Printf("Ошибка при создании Pinger: %s\n", err)
//		return 0, false
//	}
//
//	pinger.Count = 5
//	pinger.Timeout = time.Second * 5
//	err = pinger.Run()
//	if err != nil {
//		log.Printf("Ошибка при выполнении пинга: %s\n", err)
//		return 0, false
//	}
//
//	stats := pinger.Statistics()
//	return int(stats.AvgRtt.Milliseconds()), stats.PacketsRecv > 0
//}

func pingFunc(ip string) (int, bool) {
	start := time.Now()                                    // Засекаем время начала
	_, err := exec.Command("ping", "-c", "4", ip).Output() // Выполняем ping
	if err != nil {
		fmt.Println(ip, err)
		return 0, false // Если ошибка, возвращаем неудачный статус
	}

	return int(time.Since(start).Milliseconds()), true
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

		if _, ok := containerDetails.NetworkSettings.Networks["ping_network"]; !ok {
			continue
		}
		
		ipAddress := containerDetails.NetworkSettings.Networks["ping_network"].IPAddress

		if ipAddress == "" {
			for _, netw := range containerDetails.NetworkSettings.Networks {
				ipAddress = netw.IPAddress
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
	_, err := http.Post(
		fmt.Sprintf("http://backend:8080/v1/containers/%s", result.IP),
		"application/json",
		responseBody,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
}

func uniteConteiners() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(
		client.WithHost("tcp://host.docker.internal:2375"), // Используем TCP
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatalf("Ошибка при создании клиента: %s", err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		log.Fatalf("Ошибка при получении списка контейнеров: %s", err)
	}

	networkName := "ping_network"
	net, err := cli.NetworkInspect(ctx, networkName, network.InspectOptions{})
	if err != nil {
		log.Fatalln("Ошибка сети")
	}

	for _, c := range containers {
		isConnected := false

		for k := range net.Containers {
			fmt.Println(k)
			if k == c.ID {
				isConnected = true
				break
			}
		}

		if !isConnected {
			log.Println(c.Image, c.NetworkSettings.Networks)
			err = cli.NetworkConnect(ctx, networkName, c.ID, nil)
			if err != nil {
				log.Printf("Ошибка подключения контейнера %s к сети %s: %v", c.ID, networkName, err)
			} else {
				log.Printf("Контейнер %s подключен к сети %s\n", c.Image, networkName)
				// Перезапускаем контейнер после подключения к сети
				err = cli.ContainerRestart(ctx, c.ID, container.StopOptions{})
				if err != nil {
					log.Printf("Ошибка перезапуска контейнера %s: %v", c.Image, err)
				} else {
					log.Printf("Контейнер %s перезапущен\n", c.Image)
				}
			}
		} else {
			log.Printf("Контейнер %s уже подключен к сети %s\n", c.Image, networkName)
		}
	}

	log.Println("Все контейнеры были обработаны")
}

func main() {
	for {
		uniteConteiners()
		dockerContainers := takeDockerContainers()
		fmt.Println(dockerContainers)

		for _, ip := range dockerContainers {
			pingTime, success := pingFunc(ip.Ip) // Пингуем IP-адрес

			if success {
				sendPingResult(PingResult{
					IP:             ip.Ip,
					PingTime:       pingTime,
					Success:        success,
					LastSuccessful: time.Now().UTC(),
				})
			}

			log.Printf("IP: %s, PingTime: %d Success: %t Image: %s\n", ip.Ip, pingTime, success, ip.Image)
		}

		time.Sleep(10 * time.Second)
		log.Printf("sleeping %d second\n", 10)
	}

}
