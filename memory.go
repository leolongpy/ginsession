package ginsession

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
)

type MemSD struct {
	ID     string
	Data   map[string]interface{}
	rwLock sync.RWMutex
}

func NewMemorySessionData(id string) SessionData {
	return &MemSD{
		ID:   id,
		Data: make(map[string]interface{}, 8),
	}
}

func (m *MemSD) GetID() string {
	return m.ID
}

// Get 根据key获取值
func (m *MemSD) Get(key string) (value interface{}, err error) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	value, ok := m.Data[key]
	if !ok {
		err = fmt.Errorf("invalid key")
		return
	}
	return
}

// Set 创建数据
func (m *MemSD) Set(key string, value interface{}) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	m.Data[key] = value
}

// Del 删除key
func (m *MemSD) Del(key string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	delete(m.Data, key)
}
func (m *MemSD) Load(sessionID string) (err error) {
	return
}

//Save 保存session data
func (m *MemSD) Save() {
	return
}

func (m *MemSD) SetExpire(expired int) {
	return
}

type MemoryMgr struct {
	Session map[string]SessionData
	rwLock  sync.RWMutex
}

func NewMemoryMgr() Mgr {
	return &MemoryMgr{
		Session: make(map[string]SessionData, 1024),
	}
}

func (m *MemoryMgr) Init(addr string, options ...string) {
	return
}
func (m *MemoryMgr) GetSessionData(sessionID string) (sd SessionData, err error) {
	m.rwLock.RUnlock()
	defer m.rwLock.RUnlock()
	sd, ok := m.Session[sessionID]
	if !ok {
		err = fmt.Errorf("invalid session id")
		return
	}
	return
}

func (m *MemoryMgr) CreateSession() (sd SessionData) {
	uuidObj := uuid.NewV4()
	sd = NewMemorySessionData(uuidObj.String())
	m.Session[sd.GetID()] = sd
	return
}
