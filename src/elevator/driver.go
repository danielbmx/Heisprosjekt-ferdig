package elevator

import( 
		"fmt"
		"time"
)
// Direction
type Direction int
const (
    NONE Direction = iota
    UP
    DOWN
)
var prevDir Direction = NONE

// Light Value (ON/OFF)
type LightVal int
const (
    ON LightVal = iota
    OFF
)


// Button type
type Button struct {
    Floor int
    Dir Direction
}

// Speed
const MAX_SPEED = 3000
const MIN_SPEED = 2048

// Initialize system and drive car down to closest floor
func Init (buttonEventChan chan Button, floorEventChan chan int, initFloorChan chan int){

    // Check if hardware can be initialized:
    val := IoInit()
    if !val {
        fmt.Printf("Driver initiated\n")
    } else {
        fmt.Printf("Driver not initiated\n")
    }
    
    // Clear all buttons lights
    ClearButtons()
    
    SetMotorDir(NONE)

   	Poller(buttonEventChan, floorEventChan, initFloorChan)
    
    for{
        time.Sleep(25*time.Millisecond)
        select{
            case floor := <- initFloorChan:
                ElevatorStop(DOWN)
				SetFloorLight(floor)
                return
            default:
                SetMotorDir(DOWN)
        }
    }
}


// Polling channels every 50 ms
func Poller(buttonEventChan         chan Button,
            floorEventChan          chan int,
            initFloorChan           chan int) {
	
	
    var floorMap = map[int] int {
        SENSOR1 : 0,
        SENSOR2 : 1,
        SENSOR3 : 2,
        SENSOR4 : 3,
    }

    var buttonMap = map[int] Button {
        FLOOR_COMMAND1 : { 0, NONE },
        FLOOR_COMMAND2 : { 1, NONE },
        FLOOR_COMMAND3 : { 2, NONE },
        FLOOR_COMMAND4 : { 3, NONE },
        FLOOR_UP1      : { 0,   UP },
        FLOOR_UP2      : { 1,   UP },
        FLOOR_UP3      : { 2,   UP },
        FLOOR_DOWN2    : { 1, DOWN },
        FLOOR_DOWN3    : { 2, DOWN },
        FLOOR_DOWN4    : { 3, DOWN },
    }


	
    buttonList := make(map[int]bool)
    for key, _ := range buttonMap {
        buttonList[key] = Read_bit(key)
    }

    floorList := make(map[int]bool)

    go func(){
        initCounter := 0
        for {
            time.Sleep(25*time.Millisecond)
            
            for key, floor := range floorMap {
                newValue := Read_bit(key)
                if newValue != floorList[key]  &&  newValue == true {
                    newFloor := floor
                    
                    for initCounter < 1 {
                        initFloorChan <- newFloor
                        initCounter++
                    }
                    
                    
                    floorEventChan <- newFloor
                
                }
                floorList[key] = newValue
            }
        }
     }()

    go func(){
        for {
            time.Sleep(25*time.Millisecond)
            for key, btn := range buttonMap {
                newValue := Read_bit(key)
                if newValue && !buttonList[key] {
                    newButton := btn
                    buttonEventChan <- newButton
                }
                buttonList[key] = newValue
            }
        }
        
    }()


}

// Sets motor direction
func SetMotorDir(newDir Direction) {
    if (newDir == NONE) && (prevDir == UP) {
        Set_bit(MOTORDIR)
        Write_analog(MOTOR, MIN_SPEED)
    } else if (newDir == NONE) && (prevDir == DOWN) {
        Clear_bit(MOTORDIR)
        Write_analog(MOTOR, MIN_SPEED)
    } else if (newDir == UP) {
        Clear_bit(MOTORDIR)
        Write_analog(MOTOR, MAX_SPEED)
    } else if (newDir == DOWN) {
        Set_bit(MOTORDIR)
        Write_analog(MOTOR, MAX_SPEED)
    } else {
        Write_analog(MOTOR, MIN_SPEED)
    }
    prevDir = newDir
}

// Setting Button lights
func SetButtonLight(floor int, dir Direction, onoff LightVal) {
    switch onoff {
    case ON:
        switch {
        case    floor == 0 && dir == NONE:
                Set_bit(LIGHT_COMMAND1)
        case    floor == 1 && dir == NONE:
                Set_bit(LIGHT_COMMAND2)
        case    floor == 2 && dir == NONE:
                Set_bit(LIGHT_COMMAND3)
        case    floor == 3 && dir == NONE:
                Set_bit(LIGHT_COMMAND4)
        case    floor == 0 && dir == UP:
                Set_bit(LIGHT_UP1)
        case    floor == 1 && dir == UP:
                Set_bit(LIGHT_UP2)

        case    floor == 2 && dir == UP:
                Set_bit(LIGHT_UP3)
        case    floor == 1 && dir == DOWN:
                Set_bit(LIGHT_DOWN2)
        case    floor == 2 && dir == DOWN:
                Set_bit(LIGHT_DOWN3)
        case    floor == 3 && dir == DOWN:
                Set_bit(LIGHT_DOWN4)
        }
    case OFF:
        switch {
        case    floor == 0 && dir == NONE:
                Clear_bit(LIGHT_COMMAND1)
        case    floor == 1 && dir == NONE:
                Clear_bit(LIGHT_COMMAND2)
        case    floor == 2 && dir == NONE:
                Clear_bit(LIGHT_COMMAND3)
        case    floor == 3 && dir == NONE:
                Clear_bit(LIGHT_COMMAND4)
        case    floor == 0 && dir == UP:
                Clear_bit(LIGHT_UP1)
        case    floor == 1 && dir == UP:
                Clear_bit(LIGHT_UP2)
        case    floor == 2 && dir == UP:
                Clear_bit(LIGHT_UP3)
        case    floor == 1 && dir == DOWN:
                Clear_bit(LIGHT_DOWN2)
        case    floor == 2 && dir == DOWN:
                Clear_bit(LIGHT_DOWN3)
        case    floor == 3 && dir == DOWN:
                Clear_bit(LIGHT_DOWN4)
        }
    }
}

// Setting floor light
func SetFloorLight(floor int) {
    switch floor {
    case 0:
        Clear_bit (FLOOR_IND1)
        Clear_bit (FLOOR_IND2)
    case 1:
        Clear_bit (FLOOR_IND1)
        Set_bit   (FLOOR_IND2)
    case 2:
        Set_bit   (FLOOR_IND1)
        Clear_bit (FLOOR_IND2)
    case 3:
        Set_bit   (FLOOR_IND1)
        Set_bit   (FLOOR_IND2)
    }
}

// Setting Stop Button light
func SetStopButtonLight(onoff LightVal) {
    switch {
    case onoff == ON:
        Set_bit(LIGHT_STOP)
    case onoff == OFF:
        Clear_bit(LIGHT_STOP)
    }
}

// Setting Door Open lamp
func SetDoorOpenLight(onoff LightVal) {
    switch {
    case onoff == ON:
        Set_bit(DOOR_OPEN)
    case onoff == OFF:
        Clear_bit(DOOR_OPEN)
    }
}

// Stoping elevator
func ElevatorStop(direction Direction) {
	if direction == UP {
		SetMotorDir(DOWN)
		time.Sleep(8*time.Millisecond)
		SetMotorDir(NONE)
	}
	if direction == DOWN {
		SetMotorDir(UP)
		time.Sleep(8*time.Millisecond)
		SetMotorDir(NONE)
	}else{
		SetMotorDir(NONE)
	}
}

// clear buttons:
func ClearButtons() {
	SetDoorOpenLight(OFF)
	for i := 0; i < 4; i++ {
		SetButtonLight(i, NONE, OFF)
		SetButtonLight(i, UP, OFF)
		SetButtonLight(i, DOWN, OFF)
	}
}





























