package udpchan

import (
	"fmt"
	"net"
)

type UdpChan struct {
	lconn   *net.UDPConn
	sconn   *net.UDPConn
	receive *chan []byte
	send    chan []byte
}

func NewUdpChan(port int, receive *chan []byte) (udpChan *UdpChan, err error) {
	srvAddr := fmt.Sprintf("224.0.0.1:%d", port)

	addr, err := net.ResolveUDPAddr("udp4", srvAddr)
	if err != nil {
		return nil, err
	}

	lconn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	sconn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		lconn.Close()
		return nil, err
	}

	udpChan = &UdpChan{
		lconn:   lconn,
		sconn:   sconn,
		receive: receive,
		send:    make(chan []byte, 10),
	}

	go udpChan.receiveLoop()
	go udpChan.sendLoop()

	return udpChan, err

}

func (u *UdpChan) receiveLoop() {
	for {
		buffer := make([]byte, 1024)
		n, _, err := u.lconn.ReadFromUDP(buffer)
		if err != nil {
			return
		}
		*u.receive <- buffer[:n]
	}
}

func (u *UdpChan) sendLoop() {
	for {
		if toSend, ok := <-u.send; ok {
			u.sconn.Write(toSend)
		}
	}
}

func (u *UdpChan) Close() (err error) {

	close(u.send)

	if u.lconn != nil {
		err = u.lconn.Close()
	}
	if u.sconn != nil {
		err = u.sconn.Close()
	}
	return err
}

func (u *UdpChan) Send(msg []byte) {
	u.send <- msg
}
