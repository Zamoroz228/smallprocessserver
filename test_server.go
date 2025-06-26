package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Process struct{
	start_time time.Time
	end_time time.Time
	process_running bool
}

type Main_data struct{
	massive map[int] *Process
}

func ImitationOfActivity(data *Process){
	time.Sleep(time.Duration(rand.Intn(3)+3)*time.Minute)
	data.process_running = false
	data.end_time = time.Now()
}

func PathToId(w http.ResponseWriter,r *http.Request) int{
	list_with_path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(list_with_path[len(list_with_path)-1])
	if err != nil {
		fmt.Fprintln(w, "Id процесса должно быть в числовом виде")
		return -1
	}
	return id

}

func (server *Main_data) AddProcess(w http.ResponseWriter, r *http.Request) {
	id := PathToId(w, r)
	if id==-1{
		return
	}
	if _, exist := server.massive[id]; exist{
		fmt.Fprintln(w, "Такой процесс уже существует")
		return
	}
	server.massive[id] = &Process{
		start_time: time.Now(),
		process_running: true,
	}
	go ImitationOfActivity(server.massive[id])
	fmt.Fprintf(w,"Процесс %d успешно запущен", id)
}

func (server *Main_data) RemoveProcess(w http.ResponseWriter, r *http.Request) {
	id := PathToId(w, r)
	if id==-1{
		return
	}
	if _, exist := server.massive[id]; !exist{
		fmt.Fprintln(w, "Такого процесса нет")
		return
	}
	delete(server.massive,id)
	fmt.Fprintf(w,"Процесс %d успешно удален", id)
}

func (server *Main_data) Info(w http.ResponseWriter, r *http.Request) {
	id := PathToId(w, r)
	if id==-1{
		return
	}
	if process, exist := server.massive[id]; !exist{
		fmt.Fprintln(w, "Такого процесса нет")
		return
	} else{
		if !process.process_running{
			fmt.Fprintln(w, "Процесс закончил работу")
			fmt.Fprintln(w, process.start_time.Format("Процесс был запущен в 15:04:05"))
			fmt.Fprintln(w, process.end_time.Format("Процесс закончил работу в 15:04:05"))
			fmt.Fprintf(w,"Процесс работал %d секунд", process.end_time.Sub(process.start_time) / time.Second)
		} else{
			fmt.Fprintln(w, "Процесс в работе")
			fmt.Fprintln(w, process.start_time.Format("Процесс был запущен в 15:04:05"))
			fmt.Fprintf(w,"Процесс работает уже %d секунд", time.Now().Sub(process.start_time) / time.Second)
		}
	}
}

func main() {
	server := &Main_data{massive: make(map[int]*Process)}
    http.HandleFunc("/add/", server.AddProcess)
	http.HandleFunc("/remove/", server.RemoveProcess)
	http.HandleFunc("/info/",server.Info)
    log.Fatal(http.ListenAndServe(":8080", nil))
}