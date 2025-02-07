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
	IP       string `json:"ip"`        // IP-адрес контейнера
	PingTime int    `json:"ping_time"` // Время пинга в миллисекундах
	Success  bool   `json:"success"`   // Успешен ли pingfunc
}

//func pingfunc(ip string) (int, bool) {
//	start := time.Now()                        // Засекаем время начала
//	cmd := exec.Command("pingfunc", "-n", "1", ip) // Используем "-n" для Windows
//
//	_, err := cmd.CombinedOutput()
//	if err != nil {
//		fmt.Printf("Ошибка: %s\n", err)
//	}
//
//	return int(time.Since(start).Milliseconds()), true // Возвращаем время и успешный статус
//}

func pingfunc(ip string) (int, bool) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		fmt.Printf("Ошибка при создании Pinger: %s\n", err)
		return 0, false
	}

	pinger.Count = 1
	pinger.Timeout = time.Second * 5
	err = pinger.Run()
	if err != nil {
		fmt.Printf("Ошибка при выполнении пинга: %s\n", err)
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
	for {
		dockerContainers := takeDockerContainers()

		for _, ip := range dockerContainers {
			pingTime, success := pingfunc(ip.Ip) // Пингуем IP-адрес
			//sendPingResult(PingResult{
			//	IP:       ip.Ip,
			//	PingTime: pingTime,
			//	Success:  success,
			//})
			fmt.Printf("IP: %s, PingTime: %d ", ip.Ip, pingTime)
			fmt.Printf("Success: %t\n", success)
		}

		time.Sleep(10 * time.Second)
	}

}

// написать конфиг и через него прокидывать адрес
//docker build -t pings .
// docker run -p 8081:8081 pings
