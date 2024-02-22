package api

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net"
	"next-terminal/server/model"
	"next-terminal/server/service"
)

type NetWorkApi struct {
}

func (api NetWorkApi) GetNetWork(c echo.Context) error {
	networks, err := getNetworkInfo()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// 输出网络接口信息
	for _, network := range networks {
		fmt.Printf("Name: %s\n", network.Name)
		fmt.Printf("IP Address: %s\n", network.IP)
		fmt.Printf("NetMask: %s\n", network.NetMask)
		fmt.Printf("Gateway: %s\n", network.Gateway)
		fmt.Println()
	}

	return Success(c, networks)
}
func getNetworkInfo() ([]model.NetWork, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var networkInfo []model.NetWork
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error:", err)
				addrs = nil
			}
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}
				if ipNet.IP.To4() != nil {
					// 找到IPv4地址
					ip := ipNet.IP
					mask := ipNet.Mask
					gateway, err := getGateway(iface)
					if err != nil {
						fmt.Println("Error getting gateway:", err)
						gateway = ""
					}
					networkInfo = append(networkInfo, model.NetWork{
						Name:    iface.Name,
						IP:      ip.String(),
						NetMask: net.IP(mask).String(),
						Gateway: gateway,
					})
				}
			}
		}
	}
	return networkInfo, nil
}

// 获取网关地址
func getGateway(iface net.Interface) (string, error) {
	routes, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, route := range routes {
		if ipNet, ok := route.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.Contains(net.IPv4zero) {
				gatewayIP := ipNet.IP.Mask(ipNet.Mask)
				return gatewayIP.String(), nil
			}
		}
	}
	// 如果未找到网关，则返回空字符串
	return "", nil
}
func (api NetWorkApi) AddNetWork(c echo.Context) error {
	// 获取接口列表
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// 打印每个接口的信息
	for _, dev := range devices {
		fmt.Println("Name:", dev.Name)
		fmt.Println("Description:", dev.Description)
		fmt.Println("Addresses:")
		for _, addr := range dev.Addresses {
			fmt.Println("- IP:", addr.IP)
			fmt.Println("- Netmask:", addr.Netmask)
		}
		fmt.Println()
	}
	machineCode := ""
	machineCode = service.LicenseService.GetMachineId()
	return Success(c, machineCode)
}
