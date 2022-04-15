package go_moonraker

import (
	"fmt"
	"github.com/creachadair/jrpc2"
	"github.com/stretchr/testify/assert"
	"testing"
)

var client *MoonClient

func init() {
	client, _ = NewClient("10.0.0.250:80", "websocket", func(*jrpc2.Request) { return })
}

func TestMoonClient_Info(t *testing.T) {
	assert := assert.New(t)

	info, err := client.Info()
	fmt.Printf("%#v\n", info)
	assert.NoError(err)
}

func TestMoonClient_Identify(t *testing.T) {
	assert := assert.New(t)
	params := &IdentifyParams{
		ClientName: "bot",
		Version:    "0.0.1",
		Type:       "bot",
		Url:        "https://example.com",
	}
	id, err := client.Identify(params)
	assert.NoError(err)
	assert.NotNil(id)
	fmt.Printf("Id is: %d\n", id)
}

func TestMoonClient_ListObjects(t *testing.T) {
	assert := assert.New(t)
	objs, err := client.ListObjects()
	assert.NoError(err)
	assert.NotNil(objs)
	fmt.Printf("%#v\n", objs)
}

func TestMoonClient_QueryObject(t *testing.T) {
	assert := assert.New(t)
	var results interface{}
	params := QueryObjectParams{
		map[string]interface{}{"gcode_move": nil},
	}
	err := client.QueryObject(params, &results)
	assert.NoError(err)
	assert.NotNil(results)
	fmt.Printf("%#v\n", results)
}

func TestMoonClient_QueryEndstops(t *testing.T) {
	assert := assert.New(t)
	endstops, err := client.QueryEndstops()
	assert.NoError(err)
	assert.NotNil(endstops)
	fmt.Printf("%#v\n", endstops)
}

func TestMoonClient_QueryServerInfo(t *testing.T) {
	assert := assert.New(t)
	info, err := client.QueryServerInfo()
	assert.NoError(err)
	assert.NotNil(info)
	fmt.Printf("%#v\n", info)
}

type TempStore struct {
	Extruder         PTempStore `json:"extruder"`
	HeaterBed        PTempStore `json:"heater_bed"`
	TempSensorMCU    STempStore `json:"temperature_sensor mcu"`
	TempSensorRaspPi STempStore `json:"temperature_sensor raspberry_pi"`
}

type PTempStore struct {
	Powers  []float32 `json:"powers"`
	Targets []float32 `json:"targets"`
	Temps   []float32 `json:"temperatures"`
}

type STempStore struct {
	Temps []float32 `json:"temperatures"`
}

func TestMoonClient_TemperatureStore(t *testing.T) {
	assert := assert.New(t)
	var results TempStore
	err := client.TemperatureStore(&results)
	assert.NoError(err)
	assert.NotNil(results)
	fmt.Printf("%#v\n", results)
}

func TestMoonClient_GcodeStore(t *testing.T) {
	assert := assert.New(t)
	store, err := client.GcodeStore(10)
	assert.NoError(err)
	assert.NotNil(store)
	fmt.Printf("%#v\n", store)
}
