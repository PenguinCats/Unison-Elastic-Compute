package operation

import "sync"

type OprInfo struct {
	CallbackURL string
}

type OprInfoMap struct {
	mp map[string]OprInfo
	mu sync.RWMutex
}

func (o *OprInfoMap) GetOptInfo(OperationID string) OprInfo {
	o.mu.RLock()
	defer o.mu.RUnlock()
	info, _ := o.mp[OperationID]
	return info
}

func (o *OprInfoMap) SetOptInfo(OperationID string, info OprInfo) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.mp[OperationID] = info
}

func (o *OprInfoMap) DelOptInfo(operationID string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.mp, operationID)
}

var (
	OprInfoUtil = OprInfoMap{
		mp: make(map[string]OprInfo),
		mu: sync.RWMutex{},
	}
)
