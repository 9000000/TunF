//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func (p *ProxyService) addFirewallRule() {
	port := strings.TrimPrefix(p.listenAddr, ":")
	ruleName := fmt.Sprintf("TunF_%s", port)

	// Delete any existing rule with same name first
	exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name="+ruleName).Run()

	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name="+ruleName,
		"dir=in",
		"action=allow",
		"protocol=TCP",
		"localport="+port)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to add firewall rule: %v\n", err)
	}
}

func (p *ProxyService) removeFirewallRule() {
	port := strings.TrimPrefix(p.listenAddr, ":")
	ruleName := fmt.Sprintf("TunF_%s", port)
	cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name="+ruleName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to remove firewall rule: %v\n", err)
	}
}
