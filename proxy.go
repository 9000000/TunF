package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
)

type ProxyService struct {
	listener       net.Listener
	ctx            context.Context
	cancel         context.CancelFunc
	running        bool
	mu             sync.Mutex
	listenAddr     string
	targetAddr     string
	manageFirewall bool
}

func NewProxyService() *ProxyService {
	return &ProxyService{}
}

func (p *ProxyService) Start(listenAddr, targetAddr string, manageFirewall bool) error {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return fmt.Errorf("proxy is already running")
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		p.mu.Unlock()
		return fmt.Errorf("failed to listen on %s: %v", listenAddr, err)
	}

	p.listener = listener
	p.listenAddr = listenAddr
	p.targetAddr = targetAddr
	p.manageFirewall = manageFirewall
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.running = true
	p.mu.Unlock()

	if manageFirewall {
		go p.addFirewallRule()
	}

	go p.acceptLoop()

	return nil
}

func (p *ProxyService) acceptLoop() {
	defer func() {
		p.mu.Lock()
		p.running = false
		p.mu.Unlock()
	}()

	for {
		conn, err := p.listener.Accept()
		if err != nil {
			select {
			case <-p.ctx.Done():
				return
			default:
				fmt.Printf("Accept error: %v\n", err)
				continue
			}
		}

		go p.handleConnection(conn)
	}
}

func (p *ProxyService) handleConnection(src net.Conn) {
	defer src.Close()

	dest, err := net.Dial("tcp", p.targetAddr)
	if err != nil {
		fmt.Printf("Dial error to %s: %v\n", p.targetAddr, err)
		return
	}
	defer dest.Close()

	done := make(chan struct{}, 2)

	go func() {
		io.Copy(dest, src)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(src, dest)
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-p.ctx.Done():
	}
}

func (p *ProxyService) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return
	}

	if p.manageFirewall {
		go p.removeFirewallRule()
	}

	if p.cancel != nil {
		p.cancel()
	}

	if p.listener != nil {
		p.listener.Close()
	}

	p.running = false
}

func (p *ProxyService) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}
