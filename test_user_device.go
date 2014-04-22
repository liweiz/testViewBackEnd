package main

import (
	"fmt"
)

type userX struct {
	UserNo   int
	Email    string
	Password string
	UserId   string
	Devices  []deviceX
}

type deviceX struct {
	DeviceNo       int
	DeviceUuid     string
	DeviceUuidCode string
	AccessToken    string
	RefreshToken   string
	ReqId          []string
}

func (u *userX) SetUserTokens(p *publicDataSet) {
	fmt.Println("To user: ", u.UserNo)
	if u.Devices[0].DeviceUuid == p.Uuid {
		u.Devices[0].AccessToken = p.AccessToken
		u.Devices[0].RefreshToken = p.RefreshToken
		fmt.Println("To device: ", u.Devices[0].DeviceNo)
	} else {
		u.Devices[1].AccessToken = p.AccessToken
		u.Devices[1].RefreshToken = p.RefreshToken
		fmt.Println("To device: ", u.Devices[1].DeviceNo)
	}
}
