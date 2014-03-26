package main

import (
		"./networkmodule"
        "./elevator"
        "fmt"
        "time"
        
)

func main(){
   
   	fmt.Println("main")


	// Driver channels:
	ButtonEventChan 	:= make(chan elevator.Button, 1)
	FloorEventChan 		:= make(chan int, 1)
	InitFloorChan       := make(chan int, 1)
   	
	// Channels between panel and ordersystem:
	SetLightChan 		:= make(chan elevator.OrderSetLight, 1)
	BtnPanelToOrderChan := make(chan elevator.Button, 1)

	// Channels between statemachine and ordersystem:
	OrderTakenChan      := make(chan elevator.OrderSetLight, 1)
	OrderToFSMChan		:= make(chan elevator.Button, 1)
	LocalClientFSMToOrderChan	:= make(chan elevator.LocalClient, 1)


	// Channels between net and ordersystem:
	ClientOrderToNetChan := make(chan elevator.LocalClient, 1)
	ClientNetToOrderChan 	:= make(chan elevator.LocalClient, 1)
	BtnOrderToNetChan      := make(chan elevator.Button, 1)
	BtnNetToOrderChan	:= make(chan elevator.Button, 1)


	// Channels between Net and UDP: 
	BtnFromUDPChan		:= make(chan elevator.Button, 1)
	ClientFromUDPChan	:= make(chan elevator.LocalClient, 1)

	elevator.Init(ButtonEventChan, FloorEventChan, InitFloorChan)

	go elevator.PanelHandler(ButtonEventChan, SetLightChan, BtnPanelToOrderChan)
	
	go elevator.OrderHandler(SetLightChan, BtnPanelToOrderChan, OrderTakenChan, OrderToFSMChan, LocalClientFSMToOrderChan, ClientOrderToNetChan , ClientNetToOrderChan, BtnOrderToNetChan, BtnNetToOrderChan)
	
	go networkmodule.NetworkHandler(BtnOrderToNetChan, BtnFromUDPChan, BtnNetToOrderChan, ClientOrderToNetChan, ClientFromUDPChan, ClientNetToOrderChan)
	

    go elevator.UpdateState(FloorEventChan, OrderToFSMChan, OrderTakenChan, LocalClientFSMToOrderChan)
    
	
	for {
		time.Sleep(10000*time.Hour)
	}
}




















