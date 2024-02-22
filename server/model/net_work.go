package model

type NetWork struct {
	Name    string ` json:"name"`
	IP      string ` json:"ip"`
	NetMask string ` json:"netMask"`
	Gateway string ` json:"gateway"`
}
