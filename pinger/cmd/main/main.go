package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"io/ioutil"
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
	IP       string `json:"ip"`        // IP-адрес контейнера
	PingTime int    `json:"ping_time"` // Время пинга в миллисекундах
	Success  bool   `json:"success"`   // Успешен ли ping
}

func ping(ip string) (int, bool) {
	start := time.Now()                                    // Засекаем время начала
	_, err := exec.Command("ping", "-c", "1", ip).Output() // Выполняем ping
	if err != nil {
		return 0, false // Если ошибка, возвращаем неудачный статус
	}
	return int(time.Since(start).Milliseconds()), true // Возвращаем время и успешный статус
}

func takeDockerContainers() []dockerContainer {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Ошибка при создании Docker-клиента: %s", err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		log.Fatalf("Ошибка при получении списка контейнеров: %s", err)
	}

	var dockerContainers []dockerContainer
	for _, container := range containers {
		containerDetails, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Printf("Ошибка при получении деталей контейнера %s: %s", container.ID[:10], err)
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
			Id:     container.ID,
			Image:  container.Image,
			Status: container.Status,
			Ip:     ipAddress,
		})

		fmt.Printf(
			"ID: %s, Image: %s, Status: %s, IP: %s\n",
			container.ID,
			container.Image,
			container.Status,
			ipAddress,
		)
	}

	return dockerContainers
}

func sendPingResult(result PingResult) {
	postBody, _ := json.Marshal(result)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		fmt.Sprintf("localhost:8080/v1/containers/%s", result.IP),
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

func main() {
	dockerContainers := takeDockerContainers()

	for _, ip := range dockerContainers {
		pingTime, success := ping(ip.Ip) // Пингуем IP-адрес
		if success {
			// Если ping успешен, отправляем данные в Backend
			sendPingResult(PingResult{
				IP:       ip.Ip,
				PingTime: pingTime,
				Success:  success,
			})
		}
	}

	time.Sleep(10 * time.Second)

}
