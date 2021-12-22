package types

type SlaveProfile struct {
	SlaveUUId string `json:"slave_uuid"`

	Platform        string `json:"platform"`
	PlatformFamily  string `json:"platform_family"`
	PlatformVersion string `json:"platform_version"`

	MemoryTotalSize uint64 `json:"memory_total_size"`

	CpuModelName    string `json:"cpu_model_name"`
	LogicalCoreCnt  int    `json:"logical_core_cnt"`
	PhysicalCoreCnt int    `json:"physical_core_cnt"`
}

type APISlaveUUIDListResponse struct {
	SlavesUUID []string `json:"slaves"`
}

type APISlaveProfileListRequest struct {
	SlavesUUID []string `json:"slaves"`
}

type APISlaveProfileListResponse struct {
	Slaves []SlaveProfile `json:"slaves"`
}

type APISlaveStatusRequest struct {
	SlaveUUID []string `json:"uuids"`
}

type APISlaveStatusItem struct {
	UUID             string `json:"uuid"`
	Stats            string `json:"stats"`
	CoreAvailable    int    `json:"core_available"`
	MemAvailable     uint64 `json:"mem_available"`
	StorageAvailable uint64 `json:"storage_available"`
}

type APISlaveStatusResponse struct {
	Status []APISlaveStatusItem `json:"status"`
}

type APISlaveAddToken struct {
	Token string `json:"token"`
}

type APISlaveDeleteRequest struct {
	APICallBackRequestBase
	SlaveUUID string `json:"slave_uuid"`
}

type APISlaveDeleteResponse struct {
	APICallBackResponseBase
	SlaveUUID string `json:"slave_uuid"`
}

type APISlaveImageListRequest struct {
	SlaveUUID string `json:"slave_uuid"`
}

type ImageListItem struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	CreatedTime string `json:"created_time"`
}

type APISlaveImageListResponse struct {
	Images []ImageListItem `json:"images"`
}
