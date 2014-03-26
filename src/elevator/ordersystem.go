// Ordersystem
package elevator

import (
	//"fmt"
	"time"

) 

// Dummy variable for running only one elevator:
var takeOrder bool = true

type OrderSetLight struct {
    Floor int
    Dir Direction
	Light LightVal
}


func OrderHandler(SetLightChan chan<- OrderSetLight, BtnPanelToOrderChan <-chan Button, OrderTakenChan <-chan OrderSetLight, OrderToFSMChan chan<- Button, LocalClientFSMToOrderChan <-chan LocalClient, ClientOrderToNetChan chan<- LocalClient, ClientNetToOrderChan <-chan LocalClient, BtnOrderToNetChan chan<- Button, BtnNetToOrderChan <-chan Button) {
   	/*
   	Provide neccesary order-handling based on information from driver via channels.	
	Send buttons to be turned on/off to "panel" via setLightChan
   	*/
	
	var btnFromPanel Button
	var setBtnToPanel OrderSetLight
	var btnFromNetwork Button 
	var orderToFSM Button
	var localClient LocalClient
	
	for{
		time.Sleep(25*time.Millisecond)
		//fmt.Println("OrderHandler: her")
		select{
			// Case 1: Button pressed on own panel
			case btnFromPanel = <- BtnPanelToOrderChan:
				
				// Send order to network module if not Command Button
				if btnFromPanel.Dir != NONE {
					BtnOrderToNetChan <- btnFromPanel
				}else{
					// Tell panel to set light on
					setBtnToPanel.Floor = btnFromPanel.Floor
					setBtnToPanel.Dir = btnFromPanel.Dir
					setBtnToPanel.Light = ON
					SetLightChan <- setBtnToPanel
					
					// tell runElevator to go to that floor.
					orderToFSM.Floor = btnFromPanel.Floor
					orderToFSM.Dir = btnFromPanel.Dir
					OrderToFSMChan <- orderToFSM
				
				}
				// Calculate cost related to taking order:
                
				// Compare cost with other elevators:
				
				// If should take order, tell statemachine and other participants
				// Else, someone else takes order.
				
				/*if takeOrder == true {
					// store order in ordermatrix
					//fmt.Println("HandleOrder: true")
					
					// Tell panel to set light on
					setBtnToPanel.Floor = btnFromPanel.Floor
					setBtnToPanel.Dir = btnFromPanel.Dir
					setBtnToPanel.Light = ON
					SetLightChan <- setBtnToPanel
					
					// tell runElevator to go to that floor.
					orderToFSM.Floor = btnFromPanel.Floor
					orderToFSM.Dir = btnFromPanel.Dir
					OrderToFSMChan <- orderToFSM

				}*/
			
			// Case 2: Order is taken and can be deleted	
			case DeleteOrder := <- OrderTakenChan:
			    setBtnToPanel.Floor = DeleteOrder.Floor
				setBtnToPanel.Dir = DeleteOrder.Dir
				setBtnToPanel.Light = DeleteOrder.Light
				SetLightChan <- setBtnToPanel
				    

			// Case 3: Button recieved from Networkmodule, and needs to be taken
			case btnFromNetwork = <- BtnNetToOrderChan:
							
				// Tell panel to set light on
				setBtnToPanel.Floor = btnFromNetwork.Floor
				setBtnToPanel.Dir = btnFromNetwork.Dir
				setBtnToPanel.Light = ON
				SetLightChan <- setBtnToPanel
				
				// tell statemachine to go to that floor
				orderToFSM.Floor = btnFromNetwork.Floor
				orderToFSM.Dir = btnFromNetwork.Dir
				OrderToFSMChan <- orderToFSM
 
			case localClient = <- LocalClientFSMToOrderChan: 
				select{				
					case ClientOrderToNetChan <- localClient:
					default:
						continue
				}
			
			case clientFromNet := <- ClientNetToOrderChan:
				_ = clientFromNet
			default:
				continue
		
	
		}
	}

}




































