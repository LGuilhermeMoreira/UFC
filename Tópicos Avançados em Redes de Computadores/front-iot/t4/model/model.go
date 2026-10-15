package model

type SensorData struct {
	AX float32 `json:"X"`
	AY float32 `json:"Y"`
	AZ float32 `json:"Z"`
}

type LightSensor struct {
	Light int `json:"light"`
}
