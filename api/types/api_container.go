package types

type APIContainerCreateRequest struct {
	APICallBackRequestBase
	SlaveID        string `form:"slave_id" valid:"Required"`
	UECContainerID string `form:"uec_container_id" valid:"Required"`

	ImageName string `form:"image_name" valid:"Required"`

	ExposedTCPPorts []string `form:"exposed_tcp_ports"`
	ExposedUDPPorts []string `form:"exposed_udp_ports"`

	Mounts []string `form:"mounts"`

	CoreCnt int `form:"core_cnt" valid:"Required"`
	// max memory usage, in bytes
	MemorySize int64 `form:"memory_size" valid:"Required"`
	// max storage usage size, in bytes
	StorageSize int64 `form:"storage_size" valid:"Required"`
}

type APIContainerCreateResponse struct {
	APICallBackResponseBase `json:"api_response_base"`

	ExposedTCPPorts        []string `json:"exposed_tcp_ports"`
	ExposedTCPMappingPorts []string `json:"exposed_tcp_mapping_ports"`
	ExposedUDPPorts        []string `json:"exposed_udp_ports"`
	ExposedUDPMappingPorts []string `json:"exposed_udp_mapping_ports"`
}

type APIContainerStartRequest struct {
	APICallBackRequestBase
	SlaveID        string `form:"slave_id" valid:"Required"`
	UECContainerID string `form:"uec_container_id" valid:"Required"`
}

type APIContainerStartResponse struct {
	APICallBackResponseBase `json:"api_response_base"`
}

type APIContainerStopRequest struct {
	APICallBackRequestBase
	SlaveID        string `form:"slave_id" valid:"Required"`
	UECContainerID string `form:"uec_container_id" valid:"Required"`
}

type APIContainerStopResponse struct {
	APICallBackResponseBase `json:"api_response_base"`
}

type APIContainerRemoveRequest struct {
	APICallBackRequestBase
	SlaveID        string `form:"slave_id" valid:"Required"`
	UECContainerID string `form:"uec_container_id" valid:"Required"`
}

type APIContainerRemoveResponse struct {
	APICallBackResponseBase `json:"api_response_base"`
}
