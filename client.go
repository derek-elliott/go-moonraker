package go_moonraker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/wschannel"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type MoonClient struct {
	Conn *jrpc2.Client
	Host string
}

func logger(text string) {
	log.Info(text)
}

func NewClient(host, path string, notifyHandler func(*jrpc2.Request)) (*MoonClient, error) {
	var client MoonClient
	opts := &jrpc2.ClientOptions{
		Logger:   logger,
		OnNotify: notifyHandler,
	}
	u := url.URL{Scheme: "ws", Host: host, Path: path}
	channel, err := wschannel.Dial(u.String(), nil)
	if err != nil {
		return &MoonClient{}, err
	}
	client.Conn = jrpc2.NewClient(channel, opts)
	client.Host = host
	return &client, nil
}

func (c *MoonClient) Close() (err error) {
	if err := c.Conn.Close(); err != nil {
		return err
	}
	return
}

type IdentifyParams struct {
	ClientName string `json:"client_name"`
	Version    string `json:"version"`
	Type       string `json:"type"`
	Url        string `json:"url"`
}

type IdentifyResp struct {
	ConnectionId int `json:"connection_id"`
}

func (c *MoonClient) Identify(params *IdentifyParams) (int, error) {
	ctx := context.Background()
	var resp *IdentifyResp
	if err := c.Conn.CallResult(ctx, "server.connection.identify", params, &resp); err != nil {
		log.WithError(err).Error("call error")
		return 0, err
	}
	return resp.ConnectionId, nil
}

type PrinterInfo struct {
	State           string `json:"state"`
	StateMessage    string `json:"state_message"`
	Hostname        string `json:"hostname"`
	SoftwareVersion string `json:"software_version"`
	CpuInfo         string `json:"cpu_info"`
	KlipperPath     string `json:"klipper_path"`
	PythonPath      string `json:"python_path"`
	LogFile         string `json:"log_file"`
	ConfigFile      string `json:"config_file"`
}

func (c *MoonClient) Info() (*PrinterInfo, error) {
	ctx := context.Background()
	var resp *PrinterInfo
	if err := c.Conn.CallResult(ctx, "printer.info", nil, &resp); err != nil {
		return &PrinterInfo{}, err
	}
	return resp, nil
}

func (c *MoonClient) EmergencyStop() error {
	ctx := context.Background()
	_, err := c.Conn.Call(ctx, "printer.emergency_stop", nil)
	if err != nil {
		return err
	} else if _, ok := err.(*jrpc2.Error); ok {
		return err
	}
	return nil
}

func (c *MoonClient) FirmwareRestart() error {
	ctx := context.Background()
	resp, err := c.Conn.Call(ctx, "printer.firmware_restart", nil)
	if err != nil {
		return err
	}
	if respError := resp.Error(); respError != nil {
		errCode := respError.Code.String()
		return fmt.Errorf("RPC Error. Code: %s, Error: %s", errCode, respError.Message)
	}
	return nil
}

func (c *MoonClient) ListObjects() (*[]string, error) {
	ctx := context.Background()
	var objects struct {
		Objects []string
	}
	if err := c.Conn.CallResult(ctx, "printer.objects.list", nil, &objects); err != nil {
		return nil, err
	}
	return &objects.Objects, nil
}

type QueryObjectParams struct {
	Objects map[string]interface{} `json:"objects"`
}

func (c *MoonClient) QueryObject(params QueryObjectParams, results interface{}) error {
	ctx := context.Background()
	if err := c.Conn.CallResult(ctx, "printer.objects.query", params, results); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) Subscribe(params QueryObjectParams, results interface{}) error {
	ctx := context.Background()
	if err := c.Conn.CallResult(ctx, "printer.objects.subscribe", params, results); err != nil {
		return err
	}
	return nil
}

type Endstops struct {
	X string `json:"x"`
	Y string `json:"y"`
	Z string `json:"z"`
}

func (c *MoonClient) QueryEndstops() (*Endstops, error) {
	ctx := context.Background()
	var resp Endstops
	if err := c.Conn.CallResult(ctx, "printer.query_endstops.status", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

type ServerInfo struct {
	KlippyConnected       bool     `json:"klippy_connected"`
	KlippyState           string   `json:"klippy_state"`
	Components            []string `json:"components"`
	FailedComponents      []string `json:"failed_components"`
	RegisteredDirectories []string `json:"registered_directories"`
	Warnings              []string `json:"warnings"`
	WebsocketCount        int      `json:"websocket_count"`
	MoonrakerVersion      string   `json:"moonraker_version"`
	APIVersion            []int    `json:"api_version"`
	APIVersionString      string   `json:"api_version_string"`
}

func (c *MoonClient) QueryServerInfo() (*ServerInfo, error) {
	ctx := context.Background()
	var resp ServerInfo
	if err := c.Conn.CallResult(ctx, "server.info", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) TemperatureStore(results interface{}) error {
	ctx := context.Background()
	if err := c.Conn.CallResult(ctx, "server.temperature_store", nil, results); err != nil {
		return err
	}
	return nil
}

type GcodeStore struct {
	GcodeStore []GcodeStoreEntry `json:"gcode_store"`
}

type GcodeStoreEntry struct {
	Message string  `json:"message"`
	Time    float64 `json:"time"`
	Type    string  `json:"type"`
}

func (c *MoonClient) GcodeStore(count int) (*GcodeStore, error) {
	ctx := context.Background()
	var resp GcodeStore
	if err := c.Conn.CallResult(ctx, "server.gcode_store", struct{ count int }{count: count}, &resp); err != nil {
		return &GcodeStore{}, err
	}
	return &resp, nil
}

func (c *MoonClient) Restart() error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "server.restart", nil); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) RunGcode(code string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "printer.gcode.script", struct{ script string }{script: code}); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) GcodeHelp() (*map[string]string, error) {
	ctx := context.Background()
	var resp map[string]string
	if err := c.Conn.CallResult(ctx, "printer.gcode.help", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) Print(file string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "printer.print.start", struct{ filename string }{filename: file}); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) PausePrint() error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "printer.print.pause", nil); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) ResumePrint() error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "printer.print.resume", nil); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) CancelPrint() error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "printer.print.cancel", nil); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) MachineInfo() (*MachineInfo, error) {
	ctx := context.Background()
	var resp MachineInfo
	if err := c.Conn.CallResult(ctx, "machine.system_info", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) ShutdownOS() error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "machine.shutdown", nil); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) RebootOS() error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "machine.reboot", nil); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) RestartService(service string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "machine.services.restart", struct{ service string }{service: service}); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) StopService(service string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "machine.services.stop", struct{ service string }{service: service}); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) StartService(service string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "machine.services.start", struct{ service string }{service: service}); err != nil {
		return err
	}
	return nil
}

type ProcStats struct {
	MoonrakerStats       []MoonrakerStats `json:"moonraker_stats"`
	ThrottledState       ThrottledState   `json:"throttled_state"`
	CpuTemp              float32          `json:"cpu_temp"`
	Network              interface{}      `json:"network"`
	SystemCpuUsage       interface{}      `json:"system_cpu_usage"`
	WebsocketConnections int              `json:"websocket_connections"`
}
type MoonrakerStats struct {
	Time     float64 `json:"time"`
	CPUUsage float32 `json:"CPUUsage"`
	Memory   int     `json:"memory"`
	MemUnits string  `json:"mem_units"`
}

type ThrottledState struct {
	Bits  int      `json:"bits"`
	Flags []string `json:"flags"`
}

func (c *MoonClient) ProcStats() (*ProcStats, error) {
	ctx := context.Background()
	var resp ProcStats
	if err := c.Conn.CallResult(ctx, "machine.proc_stats", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

type MoonrakerFile struct {
	Path        string  `json:"path"`
	Modified    float64 `json:"modified"`
	Size        int     `json:"size"`
	Permissions string  `json:"permissions"`
}

func (c *MoonClient) ListFiles(root string) (*[]*MoonrakerFile, error) {
	ctx := context.Background()
	var resp []*MoonrakerFile
	if err := c.Conn.CallResult(ctx, "server.files.list", struct{ root string }{root: root}, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

type GcodeMetadata struct {
	PrintStartTime     int         `json:"print_start_time"`
	JobId              int         `json:"job_id"`
	Size               int         `json:"size"`
	Modified           float64     `json:"modified"`
	Slicer             string      `json:"slicer"`
	SlicerVersion      string      `json:"slicer_version"`
	LayerHeight        float32     `json:"layer_height"`
	FirstLayerHeight   float32     `json:"first_layer_height"`
	ObjectHeight       float32     `json:"object_height"`
	FilamentTotal      float32     `json:"filament_total"`
	EstimatedTime      int         `json:"estimated_time"`
	Thumbnails         []Thumbnail `json:"thumbnails"`
	FirstLayerBedTemp  int         `json:"first_layer_bed_temp"`
	FirstLayerExtrTemp int         `json:"first_layer_extr_temp"`
	GcodeStartByte     int         `json:"gcode_start_byte"`
	GcodeEndByte       int         `json:"gcode_end_byte"`
	Filename           string      `json:"filename"`
}

type Thumbnail struct {
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Size         int    `json:"size"`
	RelativePath string `json:"relative_path"`
}

func (c *MoonClient) GcodeMetadata(file string) (*GcodeMetadata, error) {
	ctx := context.Background()
	var resp GcodeMetadata
	if err := c.Conn.CallResult(ctx, "server.files.metadata", struct{ filename string }{filename: file}, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

type DirInfo struct {
	Dirs      []Dir           `json:"dirs"`
	Files     []MoonrakerFile `json:"files"`
	DiskUsage Usage           `json:"disk_usage"`
	RootInfo  RootInfo        `json:"root_info"`
}

type Dir struct {
	Modified    float64 `json:"modified"`
	Size        int     `json:"size"`
	Permissions string  `json:"permissions"`
	DirName     string  `json:"dirname"`
}

type Usage struct {
	Total int `json:"total"`
	Used  int `json:"used"`
	Free  int `json:"free"`
}

type RootInfo struct {
	Name        string `json:"name"`
	Permissions string `json:"permissions"`
}

func (c *MoonClient) DirectoryInfo(path string, extended bool) (*[]*DirInfo, error) {
	ctx := context.Background()
	var resp []*DirInfo
	if err := c.Conn.CallResult(ctx, "server.files.get_directory", struct {
		path     string
		extended bool
	}{path: path, extended: extended}, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) CreateDirectory(path string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "server.files.post_directory", struct{ path string }{path: path}); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) DeleteDirectory(path string, force bool) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "server.files.delete_directory", struct {
		path  string
		force bool
	}{path: path, force: force}); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) MoveFile(source string, dest string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "server.files.move", struct {
		source string
		dest   string
	}{source: source, dest: dest}); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) CopyFile(source string, dest string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "server.files.copy", struct {
		source string
		dest   string
	}{source: source, dest: dest}); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) DownloadFile(filename string, dest io.Writer) error {
	u := url.URL{Scheme: "http", Host: c.Host, Path: fmt.Sprintf("/server/files/%s", filename)}
	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if _, err := io.Copy(dest, resp.Body); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) UploadFile(filename string, data io.Reader, startPrint string) error {
	u := url.URL{Scheme: "http", Host: c.Host, Path: "/server/files/upload"}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filename)
	if _, err := io.Copy(part, data); err != nil {
		return err
	}
	if err := writer.WriteField("print", startPrint); err != nil {
		return err
	}
	writer.Close()

	r, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	if _, err := client.Do(r); err != nil {
		return err
	}
	return nil
}

func (c *MoonClient) DeleteFile(filename string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "server.files.delete_file", struct{ path string }{path: filename}); err != nil {
		return err
	}
	return nil
}

type JobQueueItems struct {
	QueuedJobs []JobQueueItem `json:"queued_jobs"`
	QueueState string         `json:"queue_state"`
}

type JobQueueItem struct {
	Filename    string  `json:"filename"`
	JobId       string  `json:"job_id"`
	TimeAdded   float64 `json:"time_added"`
	TimeInQueue float64 `json:"time_in_queue"`
}

func (c *MoonClient) ListJobQueue() (*JobQueueItems, error) {
	ctx := context.Background()
	var resp JobQueueItems
	if err := c.Conn.CallResult(ctx, "server.job_queue.status", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) QueueJobs(jobs []string) (*JobQueueItems, error) {
	ctx := context.Background()
	var resp JobQueueItems
	if err := c.Conn.CallResult(ctx, "server.job_queue.post_job", struct{ filenames []string }{filenames: jobs}, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) DeleteQueueJobs(jobIds []string) (*JobQueueItems, error) {
	ctx := context.Background()
	var resp JobQueueItems
	if err := c.Conn.CallResult(ctx, "server.job_queue.delete_job", struct{ job_ids []string }{job_ids: jobIds}, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) PauseJobQueue() (*JobQueueItems, error) {
	ctx := context.Background()
	var resp JobQueueItems
	if err := c.Conn.CallResult(ctx, "server.job_queue.pause", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) StartJobQueue() (*JobQueueItems, error) {
	ctx := context.Background()
	var resp JobQueueItems
	if err := c.Conn.CallResult(ctx, "server.job_queue.start", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

type JobHistory struct {
	Count int   `json:"count"`
	Jobs  []Job `json:"jobs"`
}

type Job struct {
	JobId         string        `json:"job_id"`
	Exists        bool          `json:"exists"`
	EndTime       float64       `json:"end_time"`
	FilamentUsed  float32       `json:"filament_used"`
	Filename      string        `json:"filename"`
	Metadata      GcodeMetadata `json:"metadata"`
	PrintDuration float64       `json:"print_duration"`
	Status        string        `json:"status"`
	StartTime     float64       `json:"start_time"`
	TotalDuration float64       `json:"total_duration"`
}

func (c *MoonClient) JobHistoryList(limit int, start int, since float64, before float64, order string) (*JobHistory, error) {
	ctx := context.Background()
	var resp JobHistory
	if err := c.Conn.CallResult(ctx, "server.history.list", struct {
		limit  int
		start  int
		since  float64
		before float64
		order  string
	}{limit, start, since, before, order}, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

type JobHistoryTotals struct {
	JobTotals HistoryTotals `json:"job_totals"`
}

type HistoryTotals struct {
	TotalJobs         int     `json:"total_jobs"`
	TotalTime         float64 `json:"total_time"`
	TotalPrintTime    float64 `json:"total_print_time"`
	TotalFilamentUsed float64 `json:"total_filament_used"`
	LongestJob        float64 `json:"longest_job"`
	LongestPrint      float64 `json:"longest_print"`
}

func (c *MoonClient) JobHistoryTotals() (*JobHistoryTotals, error) {
	ctx := context.Background()
	var resp JobHistoryTotals
	if err := c.Conn.CallResult(ctx, "server.history.totals", nil, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) ResetJobHistoryTotals() error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "server.history.reset_totals", nil); err != nil {
		return err
	}
	return nil
}

type JobHistorySingle struct {
	Job Job `json:"job"`
}

func (c *MoonClient) JobHistoryGetJob(uid string) (*JobHistorySingle, error) {
	ctx := context.Background()
	var resp JobHistorySingle
	if err := c.Conn.CallResult(ctx, "server.history.get_job", struct{ uid string }{uid}, &resp); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (c *MoonClient) JobHistoryDeleteJob(uid string) error {
	ctx := context.Background()
	if _, err := c.Conn.Call(ctx, "server.history.delete_job", struct{ uid string }{uid}); err != nil {
		return err
	}
	return nil
}
