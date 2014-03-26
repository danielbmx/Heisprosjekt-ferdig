package main
import "./elevdriver"
import (
    "fmt"
    //"time"
    //"os"
)

//const FLOORS = 4

func main() {


    buttonEventChan         := make(chan elevdriver.Button)
    floorEventChan          := make(chan int)
    //stopButtonEventChan     := make(chan bool)

    //var currButtonLights [FLOORS][3] bool

    elevdriver.Init(buttonEventChan,
                    floorEventChan,
                    nil,
                    nil)    // not listening to obstruction switch and stop button events

	fmt.Printf("Started!\n")
	
	elevdriver.SetDoorOpenLight(elevdriver.ON)
	//go elevdriver.Poller(buttonEventChan, floorEventChan, nil, nil)
	
      
/*
   fmt.Println("hit")

    for {
        select {
        case btnPress := <- buttonEventChan:
                fmt.Println("The", btnPress.Dir, "button on floor", btnPress.Floor, "has been pressed")
                currButtonLights[btnPress.Floor][btnPress.Dir] = !currButtonLights[btnPress.Floor][btnPress.Dir]
                elevdriver.SetButtonLight(btnPress.Floor, btnPress.Dir, elevdriver.ON)

        case newFloor := <- floorEventChan:
                fmt.Println("Arrived at floor", newFloor)
                elevdriver.SetFloorLight(Floor)
                switch newFloor {
                case 0:
                    elevdriver.SetMotorDir(elevdriver.UP)
                case 3:
                    elevdriver.SetMotorDir(elevdriver.DOWN)
                }

        case <- stopButtonEventChan:
                fmt.Println("The stop button has been pressed.")
                os.Exit(0)
        }
    }
    */
}











