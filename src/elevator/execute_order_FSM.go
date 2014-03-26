// Execute order functions for FSM

package elevator

import(
    //"fmt"
	"net"
)



func InitOrderMatrix() [4][3]int {
	
	var ordermatrix[4][3]int
	
	for i := 0; i<4; i++ {
		for j := 0; j<3; j++{
			ordermatrix[i][j] = 0;
		}
	}
	return ordermatrix
}



func SaveOrder(readorder Button, ordermatrix [4][3]int) [4][3]int{
	
	//time.Sleep(50*time.Millisecond)
	ordermatrix[readorder.Floor][readorder.Dir] = 1
    return ordermatrix
}


func DeleteOrder(floor int, dir Direction, ordermatrix [4][3]int, orderTakenChan chan<- OrderSetLight) [4][3]int{
   var deleteLight OrderSetLight
   deleteLight.Floor = floor
   deleteLight.Light = OFF
   
   ordermatrix[floor][dir] = 0
   deleteLight.Dir = dir 
   orderTakenChan <- deleteLight
   
   ordermatrix[floor][NONE] = 0
   deleteLight.Dir = NONE
   orderTakenChan <- deleteLight
   
   switch dir{
       case UP:
            if !OrderAbove(floor, ordermatrix) {
                ordermatrix[floor][DOWN] = 0
                deleteLight.Dir = DOWN 
                orderTakenChan <- deleteLight
            }
       case DOWN:
            if !OrderBelow(floor, ordermatrix) {
                ordermatrix[floor][UP] = 0
                deleteLight.Dir = UP
                orderTakenChan <- deleteLight
            }    
   }
   return ordermatrix
}



func StopAtFloor(dir Direction, floor int, ordermatrix [4][3]int) bool {
	// Return TRUE if order in floor in the right direction or command button
	if ordermatrix[floor][dir] != 0 || ordermatrix[floor][NONE] != 0 {
		return true
	}
	// Return TRUE if elevator is on top or bottom
	if floor == 0 || floor == 3 {
	    return true	
	}
	// Return TRUE if there is no order above or belowe
	if !(OrderAbove(floor, ordermatrix) || OrderBelow(floor, ordermatrix)) {
	    return true
	}else{
		return false
	}
}

// Returns true if there exists an order above current floor
func OrderAbove(floor int, ordermatrix [4][3]int) bool {
	for floor+=1; floor < 4; floor++{
		if ordermatrix[floor][NONE] != 0 || ordermatrix[floor][UP] != 0 || ordermatrix[floor][DOWN] != 0{
			return true
			}
	}
	return false

}

// Returns true if there exists an order below current floor
func OrderBelow(floor int, ordermatrix [4][3]int) bool {
	for floor-=1; floor >= 0; floor--{
		if ordermatrix[floor][NONE]!=0 || ordermatrix[floor][UP]!=0 || ordermatrix[floor][DOWN]!=0{
			return true
			}
	}
	return false

}

// This now returns a direction
func GetNextDirection(dir Direction, prevDir Direction, floor int, ordermatrix [4][3]int) Direction{
//    fmt.Println("Getting dir, floor:", floor)
	switch dir {
		case NONE:
			switch prevDir {
				case UP:
					if OrderAbove(floor, ordermatrix) {
						return UP
					}
				case DOWN: 
					if OrderBelow(floor, ordermatrix) {
						return DOWN
					}
			}
			if OrderBelow(floor, ordermatrix) {
				return DOWN
			}
			if OrderAbove(floor, ordermatrix) {
				return UP
			}
		case UP:
			if OrderAbove(floor, ordermatrix) {
				return UP
			}
			if OrderBelow(floor, ordermatrix) {
				return DOWN
			}
		case DOWN:
			if OrderBelow(floor, ordermatrix) {
				return DOWN
			}
			if OrderAbove(floor, ordermatrix) {
				return UP
			}
		}
    return NONE	
}

// Resets orders if one elevator is lost
func ResetOrder(elevator int, ordermatrix [4][3]int) [4][3]int {
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			if ordermatrix[i][j] == elevator{
				ordermatrix[i][j] = 1
			}
		}
	}
	return ordermatrix
}


// Function for returning local IP
func LocalIP() (net.IP, error) { 
	tt, err := net.Interfaces() 
	if err != nil { 
		return nil, err 
	} 
	for _, t := range tt { 
		aa, err := t.Addrs() 
		if err != nil { 
			return nil, err 
		} 
		for _, a := range aa { 
			ipnet, ok := a.(*net.IPNet) 
			if !ok { 
				continue 
			} 
			v4 := ipnet.IP.To4() 
			if v4 == nil || v4[0] == 127 { 
				// loopback address 
				continue
			} 
			return v4, nil 
		} 
	} 
	return nil, nil 
}




