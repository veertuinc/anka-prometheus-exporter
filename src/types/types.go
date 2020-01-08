package types

type NodeState string
const ( 
	NodeInactive       = "Offline"
	NodeInvalidLicense = "Inactive (Invalid License)"
	NodeActive         = "Active"
	NodeUpdating       = "Updating"
)

type InstanceState string
const (
	StateScheduling  = "Scheduling"
	StateStarting    = "Pulling"
	StateStarted     = "Started"
	StateStopping    = "Stopping"
	StateStopped     = "Stopped"
	StateTerminating = "Terminating"
	StateTerminated  = "Terminated"
	StateError       = "Error"
	StatePushing     = "Pushing"
)

type NodeInfo struct {
	NodeID         string               `json:"node_id"`
	NodeName       string               `json:"node_name"`
	CPU            uint                 `json:"cpu_count"`
	RAM            uint                 `json:"ram"`
	VMCount        uint                 `json:"vm_count"`
	VCPUCount      uint                 `json:"vcpu_count"`
	VRAM           uint                 `json:"vram"`
	CPUUtilization float32              `json:"cpu_util"`
	RAMUtilization float32              `json:"ram_util"`
	State          NodeState 			`json:"state"`
	Capacity       uint                 `json:"capacity"`
}

type RegistryInfo struct {
	Total uint64 	`json:"total"`
	Free  uint64	`json:"free"`
}

type InstanceInfo struct {
	InstanceID		string		`json:"instance_id"`
	Vm			 	VmData		`json:"vm"`
}

type VmData struct {
	State		InstanceState		`json:"instance_state"`
}

type Response interface {
	GetStatus() string
	GetMessage() string
	GetBody()	interface{}
}

type DefaultResponse struct {
	Status  string 		`json:"status"`
	Message string 		`json:"message"`
}

func (this *DefaultResponse) GetStatus() string {
	return this.Status
}

func (this *DefaultResponse) GetMessage() string {
	return this.Message
}

type NodesResponse struct {
	DefaultResponse
	Body []NodeInfo		`json:"body,omtiempty"`
}

func (this *NodesResponse) GetBody() interface{} {
	return this.Body
}

type RegistryResponse struct {
	DefaultResponse
	Body RegistryInfo		`json:"body,omtiempty"`
}

func (this *RegistryResponse) GetBody() interface{} {
	return this.Body
}

type InstancesResponse struct {
	DefaultResponse
	Body	[]InstanceInfo	`json:"body,omtiempty"`
}

func (this *InstancesResponse) GetBody() interface{} {
	return this.Body
}
