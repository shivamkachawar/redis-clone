package server

import (
	"fmt"
	"go-redis/aof"
	"go-redis/cache"
	"syscall"
)

func Run(cache *cache.Cache, aofFile *aof.AOF) error {
	serverFD, err := syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_STREAM,
		0,
	)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)

	err = syscall.SetsockoptInt(
		serverFD,
		syscall.SOL_SOCKET,
		syscall.SO_REUSEADDR,
		1,
	)
	if err != nil {
		return err
	}

	err = syscall.SetNonblock(serverFD, true)
	if err != nil {
		return err
	}

	err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: 6379,
	})

	if err != nil {
		return fmt.Errorf("Bind failed: %w", err)
	}
	err = syscall.Listen(serverFD, 128)
	if err != nil {
		return fmt.Errorf("Listen Failed: %w", err)
	}

	fmt.Println("Server listening on Port 6379")

	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}

	defer syscall.Close(epollFD)

	event := &syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}

	err = syscall.EpollCtl(
		epollFD,
		syscall.EPOLL_CTL_ADD,
		serverFD,
		event,
	)
	if err != nil {
		return err
	}
	events := make([]syscall.EpollEvent, 10)
	clients := make(map[int]*Client)
	for {
		n, err := syscall.EpollWait(epollFD, events, -1)

		if err == syscall.EINTR {
			continue
		}

		if err != nil {
			return err
		}

		for i := 0; i < n; i++ {

			if int(events[i].Fd) == serverFD {
				for {
					clientFD, _, err := syscall.Accept(serverFD)

					if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
						break
					}

					if err != nil {
						fmt.Println("Accept Error:", err)
						continue
					}

					err = syscall.SetNonblock(clientFD, true)
					if err != nil {
						fmt.Println("SetNonblock Error:", err)
						syscall.Close(clientFD)
						continue
					}

					clientEvent := &syscall.EpollEvent{
						Events: syscall.EPOLLIN,
						Fd:     int32(clientFD),
					}

					err = syscall.EpollCtl(
						epollFD,
						syscall.EPOLL_CTL_ADD,
						clientFD,
						clientEvent,
					)
					if err != nil {
						fmt.Println("EpollCtl Error:", err)
						syscall.Close(clientFD)
						continue
					}

					fmt.Println("Accepted Client:", clientFD)
					clients[clientFD] = &Client{
						FD:          clientFD,
						InputBuffer: make([]byte, 0),
					}
				}
			} else {

				client := clients[int(events[i].Fd)]

				alive := HandleReadableClient(client, cache, aofFile)

				if !alive {

					err := syscall.EpollCtl(
						epollFD,
						syscall.EPOLL_CTL_DEL,
						client.FD,
						nil,
					)
					if err != nil {
						fmt.Println("EpollCtl DEL Error:", err)
					}

					delete(clients, client.FD)

					err = syscall.Close(client.FD)
					if err != nil {
						fmt.Println("Close Error:", err)
					}

					fmt.Println("Cleaned up client:", client.FD)
				}

			}
		}
	}

}
