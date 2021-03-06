// Network module
package networkmodule

import (
	"fmt"
	"net"
	"time"
	"encoding/json"
	elevator "../elevator" 
) 


// Handelig all network communication
func NetworkHandler(BtnOrderToNetChan <-chan elevator.Button, BtnFromUDPChan chan elevator.Button, BtnNetToOrderChan chan<- elevator.Button, ClientOrderToNetChan <-chan elevator.LocalClient, ClientFromUDPChan chan elevator.LocalClient, ClientNetToOrderChan chan<- elevator.LocalClient  ) { // 6

	// Storing Local IP	
	var client elevator.LocalClient
	client.IpAddr, _ = elevator.LocalIP()
	

	// map with all clients joined:
	AllClients := make(map[string]elevator.LocalClient)

    // Make connections
    BtnPort := "129.241.187.255:20005"
	ClientPort := "129.241.187.255:20007"
    
    BtnConnection := UdpConnect(BtnPort)
	ClientConnection := UdpConnect(ClientPort)
	
	var noClient = make(chan bool, 1)
	
    
    // Recive structs from UDP and send on channel
    go UdpButtonReciver(BtnFromUDPChan)
    go UdpClientReciver(ClientFromUDPChan)
    
    for{
        select{
        	// Case 1: Button is recived from ordersystem
            case btnFromOrder := <- BtnOrderToNetChan:
                // Send Button via UDP to other elevators using the same port
                UdpButtonSender(btnFromOrder, BtnConnection)

			// Case 2: Button is recived via UDP
            case btnFromNet := <- BtnFromUDPChan:
				// Calculate cost related to taking the order:
				bestClient := OrderDistribute(AllClients, btnFromNet)				
				
				// Take order if you have the lowest cost			
				if bestClient.IpAddr.String() == client.IpAddr.String(){
					// Send to Ordersystem
					BtnNetToOrderChan <- btnFromNet
				}
				
			// Case 3: Client is recived from ordersystem
			case clientFromOrder := <- ClientOrderToNetChan:
                // Send LocalClient via UDP to other elevators using the same port
				UdpClientSender(clientFromOrder, ClientConnection)
			
			// Case 4: Client is recived  via UDP
			case clientFromNet := <- ClientFromUDPChan:
				// Put Client in map
				
				//elevator.TimeAfter(noClient, 3*time.Second)
				
				AllClients[clientFromNet.IpAddr.String()] = clientFromNet
				ClientNetToOrderChan <- clientFromNet
				
			case <-noClient: 
				
        }
    }
}


// Create UDP connection
func UdpConnect(address string) *net.UDPConn{
	serverAddr_udp, err := net.ResolveUDPAddr("udp", address)
	PrintError(err)

    con_udp, err := net.DialUDP("udp", nil, serverAddr_udp)
    PrintError(err)
    
    return con_udp
}

// Broadcast Button struct via UDP using Json
func UdpButtonSender(parameter elevator.Button, con_udp *net.UDPConn) {

    message, err := json.Marshal(parameter) 
    PrintError(err)
	
	_, err2 := con_udp.Write(message)
	PrintError(err2)
}

// Recieve Button struct via UDP
func UdpButtonReciver(message_channel chan elevator.Button) {
    
    serverAddr_udp, err := net.ResolveUDPAddr("udp", ":20005")
	PrintError(err)

    con_udp, err := net.ListenUDP("udp", serverAddr_udp)
    PrintError(err)
    save := elevator.Button{} 
    buffer := make([]byte,1024)
	
	for {
        n, _,err := con_udp.ReadFromUDP(buffer)
        PrintError(err)
        
        err1 := json.Unmarshal(buffer[0:n],&save)
        PrintError(err1)
        message_channel <- save
    }
    
}

// Broadcast LocalClient struct via UDP using Json
func UdpClientSender(parameter elevator.LocalClient, con_udp *net.UDPConn) {
    message, err := json.Marshal(parameter) 
    PrintError(err)
	
	for i := 0; i<2; i++ {
		time.Sleep(10 * time.Millisecond)
		_, err2 := con_udp.Write(message)
		PrintError(err2)
	}
}

// Recieve LocalClient struct via UDP
func UdpClientReciver(message_channel chan elevator.LocalClient) {
    
    serverAddr_udp, err := net.ResolveUDPAddr("udp", ":20007")
	PrintError(err)

    con_udp, err := net.ListenUDP("udp", serverAddr_udp)
    PrintError(err)
    save := elevator.LocalClient{} 
    buffer := make([]byte,1024)
    
	for {
        n, _,err := con_udp.ReadFromUDP(buffer)
        PrintError(err)
        
        err1 := json.Unmarshal(buffer[0:n],&save)
        PrintError(err1)
        message_channel <- save
    }
    
}

// Calculate and return a cost for the elevator to take an order
func GetCost(client elevator.LocalClient, btn elevator.Button) int {
	cost := 0
	// Because of bug in elevator cant take order in own floor
	if client.Floor == btn.Floor {
		cost += 10
		return cost
	}
	
	if client.CurrentDir == elevator.NONE{
		cost = client.Floor-btn.Floor
		if cost < 0 {
			cost = cost * -1
		}	
		cost += 1
		return cost
	}
	// If elevator is on its way to ordered floor
	if client.CurrentDir == elevator.UP && btn.Dir == elevator.UP && btn.Floor >= client.Floor {
		cost = btn.Floor - client.Floor
		return cost
	}
	if client.CurrentDir == elevator.DOWN && btn.Dir == elevator.DOWN && btn.Floor <= client.Floor {
		cost = client.Floor - btn.Floor
		return cost
	}
	// If elevator is driving in oposite direction
	if client.CurrentDir != btn.Dir {
		cost = client.Floor-btn.Floor
		if cost < 0 {
			cost = cost * -1
		}	
		cost += 3
		return cost
	}else{
		cost = 2
		return cost
	}
	
}


// Finds the IP with the lowest cost and gives the order
func OrderDistribute(Clients map[string]elevator.LocalClient, btn elevator.Button) elevator.LocalClient {
	least_cost := 20
	var bestClient elevator.LocalClient
	
	// Setting first IP to best IP
	for k,_ := range Clients {
		bestClient = Clients[k]
		break

	}
	
	// Go through IP-map to find best elevator to take order
	for key,_ := range Clients {
		temp_cost := GetCost(Clients[key],btn)
		if temp_cost < least_cost {
			least_cost = temp_cost
			bestClient = Clients[key]		

		}
	}
	return bestClient
}



func PrintError(err error) {
	if err != nil{
        fmt.Println(err)
	}
}

















































