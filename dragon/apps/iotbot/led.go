package main

// GreenHouse represents a green house
type GreenHouse struct {
	Id   string     `json:"id"`
	Name string     `json:"name"`
	Leds []LedModel `json:"leds"`
}

//LedModel represents a led device with network and display template.
type LedModel struct {
	IP       string        `json:"ip"`
	Port     int           `json:"port"`
	Template []DisplayItem `json:"template"`
}

// DisplayItem defines a led display item with content,posision,binding device specified point
type DisplayItem struct {
	Text     string `json:"text"`
	PosX     int    `json:"posx"`
	PosY     int    `json:"posy"`
	DeviceId string `json:"device_id"`
	PointId  string `json:"point_id"`
	Order    int    `json:"order"`
}

var greenHouses []GreenHouse

//LoadLedConf load green house config file
func LoadLedConf(cfgPath string) error {
	return loadJsonConfig(cfgPath, &greenHouses)
}

func syncLedDisplay(house GreenHouse) {
	// api interface: /getCollectorRuntimeInfoByGreenHouses

}
