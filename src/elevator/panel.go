// Panel

package elevator

import(
		
		"fmt"
		)

// Handels the panel and sets lights
func PanelHandler(buttonEventChan <-chan Button, setLightChan <-chan OrderSetLight, btnPanelToOrder chan<- Button) {
	fmt.Println("PanelHandler running")
	var passOn Button
	var setLight OrderSetLight
	
	for{
		select{
			case passOn = <- buttonEventChan: 
				// Pass on buttons to ordersystem on the btnPanelToOrder channel
				btnPanelToOrder <- passOn

			case setLight = <- setLightChan:
				// Set lights on panel
				SetButtonLight(setLight.Floor, setLight.Dir, setLight.Light)
		}

	}

}
