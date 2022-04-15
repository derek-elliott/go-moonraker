package go_moonraker

type MachineInfo struct {
	SystemInfo     SystemInfo     `json:"system_info"`
	Virtualization Virtualization `json:"virtualization"`
	Python         Python         `json:"python"`
	Network        Network        `json:"network"`
}

type SystemInfo struct {
	CPUInfo           CPUInfo      `json:"cpu_info"`
	SDInfo            SDInfo       `json:"sd_info"`
	Distribution      DistInfo     `json:"distribution"`
	AvailableServices []string     `json:"available_services"`
	ServiceState      ServiceState `json:"service_state"`
}

type CPUInfo struct {
	CPUCount     int    `json:"cpu_count"`
	Bits         string `json:"bits"`
	Processor    string `json:"processor"`
	CPUDesc      string `json:"cpu_desc"`
	SerialNumber string `json:"serial_number"`
	HardwareDesc string `json:"hardware_desc"`
	Model        string `json:"model"`
	TotalMemory  int    `json:"total_memory"`
	MemoryUnits  string `json:"memory_units"`
}

type SDInfo struct {
	ManufacturerId   string `json:"manufacturer_id"`
	Manufacturer     string `json:"manufacturer"`
	OEMId            string `json:"oem_id"`
	ProductName      string `json:"product_name"`
	ProductRevision  string `json:"product_revision"`
	SerialNumber     string `json:"serial_number"`
	ManufacturerDate string `json:"manufacturer_date"`
	Capacity         string `json:"capacity"`
	TotalBytes       int    `json:"total_bytes"`
}

type DistInfo struct {
	Name         string       `json:"name"`
	Id           string       `json:"id"`
	Version      string       `json:"version"`
	VersionParts VersionParts `json:"version_parts"`
	Like         string       `json:"like"`
	Codename     string       `json:"codename"`
}

type VersionParts struct {
	Major       string `json:"major"`
	Minor       string `json:"minor"`
	BuildNumber string `json:"build_number"`
}

type ServiceState struct {
	Klipper    StateReport `json:"klipper"`
	KlipperMCU StateReport `json:"klipper_mcu"`
	Moonraker  StateReport `json:"moonraker"`
}

type StateReport struct {
	ActiveState string `json:"active_state"`
	SubState    string `json:"sub_state"`
}

type Virtualization struct {
	VirtType       string `json:"virt_type"`
	VirtIdentifier string `json:"virt_identifier"`
}

type Python struct {
	Version       []string `json:"version"`
	VersionString string   `json:"versionString"`
}

type Network struct {
	WLan0 NetDef `json:"wlan0"`
}

type NetDef struct {
	MACAddress  string      `json:"mac_address"`
	IPAddresses []IPAddress `json:"ip_addresses"`
}

type IPAddress struct {
	Family      string `json:"family"`
	Address     string `json:"address"`
	IsLinkLocal bool   `json:"is_link_local"`
}
