package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type client chan <- string
var (
	entering = make(chan client)
	leaving = make(chan client)
	messages = make(chan string)
)

func broadcaster()  {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <- messages:
			for cli := range clients{
				cli <- msg
			}
		case cli := <- entering:
			clients[cli] = true

		case cli := <- leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn)  {
	ch := make(chan string)
	go clientWriter(conn, ch)


	who := conn.RemoteAddr().String()
	ch <- "You are"+ who
	messages <- who + "has arrived"
	entering <-ch

	input := bufio.NewScanner(conn)
	for input.Scan(){
		messages <- who + ": " +input.Text()
	}

	leaving <- ch
	messages <- who + "has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <- chan string)  {
	for msg := range ch{
		fmt.Println(conn,msg)
	}
}



func main()  {
	port := "localhost:5000"
	listener,err := net.Listen("tcp", port)
	if err != nil{
		log.Fatal(err)
	}
	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil{
			fmt.Fprintf(os.Stderr, "Detected %s\n", err.Error())
			continue
		}
		fmt.Fprintf(conn,"Welcome on board decadevs")
		go match(conn)
		go handleConn(conn)

	}
}



var partner = make(chan io.ReadWriteCloser)

func chat(a,b io.ReadWriteCloser)  {
	fmt.Println(a, "Hello")
	fmt.Println(b, "hi")
	go io.Copy(a,b)
	io.Copy(b,a)
}

func match(c io.ReadWriteCloser)  {
	//fmt.Fprintf(c, "waiting for partner...")

	select {
	case partner <-c:

	case p := <-partner:
		chat(p,c)
	}
}