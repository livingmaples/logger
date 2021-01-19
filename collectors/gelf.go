package collectors

import (
	"fmt"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"sync"
	"time"
)

type protocol string

const (
	TCP protocol = "tcp"
	UDP protocol = "udp"
)

type GelfCollector struct {
	Timeout  time.Duration
	Addr     string
	Protocol protocol

	udpWriter *gelf.UDPWriter
	tcpWriter *gelf.TCPWriter
	mu        sync.Mutex
}

func (c *GelfCollector) Ping() (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Set default protocol as UDP if not presented
	if c.Protocol == "" {
		c.Protocol = UDP
	}

	if c.Protocol == TCP && c.tcpWriter == nil {
		tcpWriter, err := gelf.NewTCPWriter(c.Addr)
		if err != nil {
			return false, fmt.Errorf("cannot connect to the gelf server: (%v)", err)
		}

		c.tcpWriter = tcpWriter
		return true, nil
	} else if c.Protocol == UDP && c.udpWriter == nil {
		udpWriter, err := gelf.NewUDPWriter(c.Addr)
		if err != nil {
			return false, fmt.Errorf("cannot connect to the gelf server: (%v)", err)
		}

		c.udpWriter = udpWriter
		return true, nil
	}

	return false, fmt.Errorf("given protocol not supported by gelf collector")
}

func (c *GelfCollector) Write(data *[]byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Protocol == TCP && c.tcpWriter != nil {
		c.tcpWriter.Write(*data)
	} else if c.Protocol == UDP && c.udpWriter != nil {
		c.udpWriter.Write(*data)
	}
}
