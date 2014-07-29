package gameserver

import (
	"fmt"
	. "protocol"
)

func Handle(api int16, payload interface{}) {
	fmt.Println("Handle api:", api)
	switch api {
	case 1:
		P_login_req(payload.(PKT_login_info))
	default:
		fmt.Println("Unknown msg:", api)
	}
}
