//go:build !windows

package main

func (p *ProxyService) addFirewallRule() {
	// TODO: Implement using iptables or ufw if needed
}

func (p *ProxyService) removeFirewallRule() {
	// TODO: Implement using iptables or ufw if needed
}
