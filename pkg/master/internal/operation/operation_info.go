package operation

import "sync"

type OprInfo struct {
	CallbackURL string
}

type OprInfoMap struct {
	mp map[int64]OprInfo
	mu sync.RWMutex
}

func (o *OprInfoMap) GetOptInfo(OperationID int64) (OprInfo, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	info, ok := o.mp[OperationID]
	return info, ok
}

func (o *OprInfoMap) SetOptInfo(OperationID int64, info OprInfo) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.mp[OperationID] = info
}

func (o *OprInfoMap) DelOptInfo(operationID int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.mp, operationID)
}

var (
	OprInfoUtil = OprInfoMap{
		mp: make(map[int64]OprInfo),
		mu: sync.RWMutex{},
	}
)
