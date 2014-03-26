// State Machine 

package elevator

import (
		"time"
		"fmt"
		"net"		
)


// States
type State int
const (
    MOVING State = iota
    STANDSTILL
    //EMG_STOPPED
    N_STATES
)

// Events
type Event int
const(
	MOVE = iota
	HALT
	//EMGSTOP
	N_EVENTS
)

type LocalClient struct {
    CurrentState State
    CurrentDir Direction
	Floor int
	IpAddr net.IP

}

// Private variables:

var timeStart time.Time
var event Event
var localClient LocalClient
var OrderMatrix = InitOrderMatrix()
var doorClose chan bool


// Statemachine:
func UpdateState(floorEventChan <-chan int, OrderToFSMChan <-chan Button, OrderTakenChan chan<- OrderSetLight, LocalClientChan chan<- LocalClient) {  
	// Variabels to use in function
	var event Event
	doorClose = make(chan bool, 2)
    var doorIsOpen bool
	var prevDir Direction

	// initialize localClient
	localClient.CurrentState = STANDSTILL
	localClient.CurrentDir = NONE
	localClient.Floor = -1
	localClient.IpAddr,_ = LocalIP()
	
	for{
		time.Sleep(25*time.Millisecond)
		
		
		//Read order(s) from ordersystem:
		select {
			case LocalClientChan <- localClient:
				// Send local client to network for broadcast
				break
		    case readOrder := <- OrderToFSMChan:
		    	// Recived order from ordersystem
        		// Save order in ordermatrix
		        OrderMatrix = SaveOrder(readOrder, OrderMatrix)
	            fmt.Println("newOrder:", OrderMatrix)
		        // Check for other orders
                if OrderAbove(localClient.Floor, OrderMatrix) || OrderBelow(localClient.Floor, OrderMatrix){
                    if !doorIsOpen { 
                        event = MOVE
                    }
                }
				
		    case newFloor := <- floorEventChan:
				// If floor is updated, check if stop is needed
		        localClient.Floor = newFloor
				SetFloorLight(newFloor)
		        if StopAtFloor(localClient.CurrentDir, localClient.Floor, OrderMatrix) {
		            
		            event = HALT
		            // Delete order from ordermatrix and tell panel to turn off lights
		            OrderMatrix = DeleteOrder(localClient.Floor, localClient.CurrentDir, OrderMatrix, OrderTakenChan)
		            fmt.Println(OrderMatrix)
		        }
		 
	        case <- doorClose:
				//If timer is out -> close door and get next direction
	            SetDoorOpenLight(OFF)
	            doorIsOpen = false
	            newDir := GetNextDirection(localClient.CurrentDir, prevDir, localClient.Floor, OrderMatrix)
	            if newDir != NONE {
	                event = MOVE
	            }
		    }


		switch localClient.CurrentState {
		
		    case MOVING:
		        switch event {
		            case MOVE:
		            	// No change is needed
		                break 
		            case HALT:
		                //Stop elevator and open door for 3 sec
		               	ElevatorStop(localClient.CurrentDir)
						SetDoorOpenLight(ON)
						doorIsOpen = true
                		go TimeAfter(doorClose, 3*time.Second)
		          		//Update client
		                localClient.CurrentState = STANDSTILL
						prevDir = localClient.CurrentDir
		                localClient.CurrentDir = NONE
                        break
                }
                break      
                       
		    case STANDSTILL:
		        switch event {		            

		            case MOVE:
		                //Figure out new direction and move
		                newDir := GetNextDirection(localClient.CurrentDir, prevDir, localClient.Floor, OrderMatrix)
		                SetMotorDir(newDir)
						//Update client
		                localClient.CurrentState = MOVING
		                localClient.CurrentDir = newDir

		                break
		            
		            case HALT:
						// No change is needed
                		break
		       }
    		   break
		   }	    
	}	
}


func TimeAfter(ch chan bool, t time.Duration){
    time.Sleep(t)
    ch <- true
}
























