package types

var NodeStates = []string{
	"Offline",
	"Inactive (Invalid License)",
	"Active",
	"Updating",
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

type Node struct {
	NodeID         string      `json:"node_id"`
	NodeName       string      `json:"node_name"`
	CPU            uint        `json:"cpu_count,omitempty"`
	RAM            uint        `json:"ram,omitempty"`
	VMCount        uint        `json:"vm_count,omitempty"`
	VCPUCount      uint        `json:"vcpu_count,omitempty"`
	VRAM           uint        `json:"vram,omitempty"`
	CPUUtilization float32     `json:"cpu_util,omitempty"`
	RAMUtilization float32     `json:"ram_util,omitempty"`
	FreeDiskSpace  uint        `json:"free_disk_space,omitempty"`
	AnkaDiskUsage  uint        `json:"anka_disk_usage,omitempty"`
	DiskSize       uint        `json:"disk_size,omitempty"`
	State          string      `json:"state"`
	Capacity       uint        `json:"capacity"`
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
	State        string `json:"instance_state"`
	TemplateUUID string `json:"vmid"`
	TemplateNAME string `json:"name"`
	GroupUUID    string `json:"group_id"`
	NodeUUID     string `json:"node_id"`
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

func (this *DefaultResponse) GetStatus() string {
	return this.Status
}

func (this *DefaultResponse) GetMessage() string {
	return this.Message
}

type NodesResponse struct {
	DefaultResponse
	Body []Node `json:"body,omitempty"`
}

func (this *NodesResponse) GetBody() interface{} {
	return this.Body
}

type RegistryDiskResponse struct {
	DefaultResponse
	Body RegistryDisk `json:"body,omitempty"`
}

func (this *RegistryDiskResponse) GetBody() interface{} {
	return this.Body
}

type Template struct {
	UUID string `json:"id"`
	Name string `json:"name"`
	Size uint   `json:"size"`
	Tags []TemplateTag
}
type RegistryTemplateResponse struct {
	DefaultResponse
	Body []Template `json:"body,omitempty"`
}

func (this *RegistryTemplateResponse) GetBody() interface{} {
	return this.Body
}

type TemplateTag struct {
	Name string `json:"tag"`
	Size uint   `json:"size"`
}

type RegistryTemplateTags struct {
	Versions []TemplateTag `json:"versions,omitempty"`
}
type RegistryTemplateTagsResponse struct {
	DefaultResponse
	Body RegistryTemplateTags `json:"body,omitempty"`
}

func (this *RegistryTemplateTagsResponse) GetBody() interface{} {
	return this.Body
}

type InstancesResponse struct {
	DefaultResponse
	Body []Instance `json:"body,omitempty"`
}

func (this *InstancesResponse) GetBody() interface{} {
	return this.Body
}
