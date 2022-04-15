package go_moonraker

type PrinterObjects struct {
	Webhooks      *Webhooks      `json:"webhooks,omitempty"`
	GcodeMove     *GcodeMove     `json:"gcode_move,omitempty"`
	Toolhead      *Toolhead      `json:"toolhead,omitempty"`
	ConfigFile    *ConfigFile    `json:"config_file,omitempty"`
	Extruder      *Extruder      `json:"extruder,omitempty"`
	HeaterBed     *HeaterBed     `json:"heater_bed,omitempty"`
	Fan           *Fan           `json:"fan,omitempty"`
	IdleTimeout   *IdleTimeout   `json:"idle_timeout,omitempty"`
	VirtualSdcard *VirtualSdcard `json:"virtual_sdcard,omitempty"`
	PrintStats    *PrintStats    `json:"print_stats,omitempty"`
	DisplayStatus *DisplayStatus `json:"display_status,omitempty"`
	BedMesh       *BedMesh       `json:"bed_mesh,omitempty"`
}

type Webhooks struct {
	State        string `json:"state,omitempty"`
	StateMessage string `json:"state_message,omitempty"`
}

type GcodeMove struct {
	SpeedFactor         float32   `json:"speed_factor,omitempty"`
	Speed               float32   `json:"speed,omitempty"`
	ExtrudeFactor       float32   `json:"extrude_factor,omitempty"`
	AbsoluteCoordinates bool      `json:"absolute_coordinates,omitempty"`
	AbsoluteExtrude     bool      `json:"absolute_extrude,omitempty"`
	HomingOrigin        []float32 `json:"homing_origin,omitempty"`
	Position            []float32 `json:"position,omitempty"`
	GcodePosition       []float32 `json:"gcode_position,omitempty"`
}

type Toolhead struct {
	HomedAxes            string    `json:"homed_axes,omitempty"`
	PrintTime            float32   `json:"print_time,omitempty"`
	EstimatedPrintTime   float32   `json:"estimated_print_time,omitempty"`
	Extruder             string    `json:"extruder,omitempty"`
	Position             []float32 `json:"position,omitempty"`
	MaxVelocity          float32   `json:"max_velocity,omitempty"`
	MaxAccel             float32   `json:"max_accel,omitempty"`
	MaxAccelToDecel      float32   `json:"max_accel_to_decel,omitempty"`
	SquareCornerVelocity float32   `json:"square_corner_velocity,omitempty"`
}

type ConfigFile struct {
	Config            map[string]string `json:"config,omitempty"`
	Settings          map[string]string `json:"settings,omitempty"`
	SafeConfigPending string            `json:"safeConfigPending,omitempty"`
}

type Extruder struct {
	Temperature     float32 `json:"temperature,omitempty"`
	Target          float32 `json:"target,omitempty"`
	Power           float32 `json:"power,omitempty"`
	PressureAdvance float32 `json:"pressureAdvance,omitempty"`
	SmoothTime      float32 `json:"smoothTime,omitempty"`
}

type HeaterBed struct {
	Temperature float32 `json:"temperature,omitempty"`
	Target      float32 `json:"target,omitempty"`
	Power       float32 `json:"power,omitempty"`
}

type Fan struct {
	Speed float32 `json:"speed"`
	Rpm   float32 `json:"rpm"`
}

type IdleTimeout struct {
	State        string  `json:"state"`
	PrintingTime float32 `json:"printingTime"`
}

type VirtualSdcard struct {
	Progress     float32 `json:"progress"`
	IsActive     bool    `json:"is_active"`
	FilePosition int     `json:"file_position"`
}

type PrintStats struct {
	Filename      string  `json:"filename"`
	TotalDuration float32 `json:"total_duration"`
	PrintDuration float32 `json:"print_duration"`
	FilamentUsed  float32 `json:"filament_used"`
	State         string  `json:"state"`
	Message       string  `json:"message"`
}

type DisplayStatus struct {
	Message  string  `json:"message"`
	Progress float32 `json:"progress"`
}

type BedMesh struct {
	ProfileName  string      `json:"profile_name"`
	MeshMin      []float32   `json:"mesh_min"`
	MeshMax      []float32   `json:"mesh_max"`
	ProbedMatrix [][]float32 `json:"probed_matrix"`
	MeshMatrix   [][]float32 `json:"mesh_matrix"`
}
