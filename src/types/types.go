package types

var ControllerStates = []string{
	"Running",
}

var RegistryStates = []string{
	"Running",
	"FAIL",
}

var NodeStates = []string{
	"Offline",
	"Inactive (Invalid License)",
	"Active",
	"Updating",
	"Drain Mode",
}

var Architectures = []string{
	"amd64",
	"arm64",
}

var InstanceStates = []string{
	"Scheduling",
	"Pulling",
	"Started",
	"Stopping",
	"Stopped",
	"Terminating",
	"Terminated",
	"Error",
	"Pushing",
}

type Status struct {
	Status          string `json:"status"`
	Version         string `json:"version"`
	RegistryAddress string `json:"registry_address"`
	RegistryStatus  string `json:"registry_status"`
	License         string `json:"license"`
}

type Node struct {
	NodeID         string      `json:"node_id"`
	NodeName       string      `json:"node_name"`
	CPU            uint        `json:"cpu_count"`
	RAM            uint        `json:"ram"`
	VMCount        uint        `json:"vm_count"`
	UsedVCPUCount  uint        `json:"vcpu_count"`
	UsedVRAM       uint        `json:"vram"`
	CPUUtilization float64     `json:"cpu_util"`
	RAMUtilization float64     `json:"ram_util"`
	FreeDiskSpace  uint        `json:"free_disk_space"`
	AnkaDiskUsage  uint        `json:"anka_disk_usage"`
	DiskSize       uint        `json:"disk_size"`
	State          string      `json:"state"`
	Capacity       uint        `json:"capacity"`
	HostArch       string      `json:"host_arch"`
	Groups         []NodeGroup `json:"groups"`
}

type NodeGroup struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	FallBackGroupId string `json:"fallback_group_id"`
}

type RegistryDisk struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

type Instance struct {
	InstanceID string `json:"instance_id"`
	Vm         VmData `json:"vm"`
}

type VmData struct {
	State          string `json:"instance_state"`
	TemplateUUID   string `json:"vmid"`
	TemplateName   string
	GroupUUID      string `json:"group_id"`
	NodeUUID       string `json:"node_id"`
	Arch           string `json:"arch"`
	CreationTime   string `json:"cr_time"`
	LastUpdateTime string `json:"ts"`
}

type Response interface {
	GetStatus() string
	GetMessage() string
	GetBody() interface{}
}

type DefaultResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (dr *DefaultResponse) GetStatus() string {
	return dr.Status
}

func (dr *DefaultResponse) GetMessage() string {
	return dr.Message
}

type StatusResponse struct {
	DefaultResponse
	Body Status `json:"body"`
}

func (sr *StatusResponse) GetBody() interface{} {
	return sr.Body
}

type NodesResponse struct {
	DefaultResponse
	Body []Node `json:"body"`
}

func (nr *NodesResponse) GetBody() interface{} {
	return nr.Body
}

type RegistryDiskResponse struct {
	DefaultResponse
	Body RegistryDisk `json:"body"`
}

func (rdr *RegistryDiskResponse) GetBody() interface{} {
	return rdr.Body
}

type Template struct {
	UUID string `json:"id"`
	Name string `json:"name"`
	Size uint   `json:"size"`
	Tags []TemplateTag
}
type RegistryTemplateResponse struct {
	DefaultResponse
	Body []Template `json:"body"`
}

func (rtr *RegistryTemplateResponse) GetBody() interface{} {
	return rtr.Body
}

type TemplateTag struct {
	Name string `json:"tag"`
	Size uint   `json:"size"`
}

type RegistryTemplateTags struct {
	Versions []TemplateTag `json:"versions"`
}
type RegistryTemplateTagsResponse struct {
	DefaultResponse
	Body RegistryTemplateTags `json:"body"`
}

func (rttr *RegistryTemplateTagsResponse) GetBody() interface{} {
	return rttr.Body
}

type InstancesResponse struct {
	DefaultResponse
	Body []Instance `json:"body"`
}

func (ir *InstancesResponse) GetBody() interface{} {
	return ir.Body
}
