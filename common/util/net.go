package util

import (
	"fmt"
	"net"
)

func demo() {
	interfaceName := "en0" // 指定要查询的网口名称

	// 获取网口信息
	interfaceInfo, err := getInterfaceInfo(interfaceName)
	if err != nil {
		fmt.Println("无法获取网口信息:", err)
		return
	}

	fmt.Println("网口名称:", interfaceInfo.Name)
	fmt.Println("MAC 地址:", interfaceInfo.HardwareAddr)
	fmt.Println("IP 地址:", interfaceInfo.IPAddresses)
}

// InterfaceInfo 网口信息结构体
type InterfaceInfo struct {
	Name         string           // 网口名称
	HardwareAddr net.HardwareAddr // MAC 地址
	IPAddresses  []net.IP         // IP 地址列表
}

// 获取特定网口的信息
func getInterfaceInfo(interfaceName string) (*InterfaceInfo, error) {
	// 获取所有网口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 查找指定名称的网口
	var targetInterface net.Interface
	found := false
	for _, iface := range interfaces {
		if iface.Name == interfaceName {
			targetInterface = iface
			found = true
			break
		}
	}

	// 未找到指定名称的网口
	if !found {
		return nil, fmt.Errorf("未找到网口 %s", interfaceName)
	}

	// 获取网口 MAC 地址
	macAddr := targetInterface.HardwareAddr

	// 获取网口 IP 地址列表
	var ipAddrs []net.IP
	addrs, err := targetInterface.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}
		ipAddrs = append(ipAddrs, ip)
	}

	// 构建网口信息结构体
	interfaceInfo := &InterfaceInfo{
		Name:         targetInterface.Name,
		HardwareAddr: macAddr,
		IPAddresses:  ipAddrs,
	}

	return interfaceInfo, nil
}
