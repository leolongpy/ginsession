package ginsession

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"sync"
	"time"
)

type RedisSD struct {
	ID      string
	Data    map[string]interface{}
	rwLock  sync.RWMutex
	expired int           // 过期时间
	client  *redis.Client //redis 连接池
}

func NewRedisSessionData(id string, client *redis.Client) SessionData {
	return &RedisSD{
		ID:     id,
		Data:   make(map[string]interface{}, 8),
		client: client,
	}
}
func (r *RedisSD) GetID() string {
	return r.ID
}

func (r *RedisSD) Get(key string) (value interface{}, err error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	value, ok := r.Data[key]
	if !ok {
		err = fmt.Errorf("invalid Key")
		return
	}
	return
}
func (r *RedisSD) Set(key string, value interface{}) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	r.Data[key] = value
}

func (r *RedisSD) Load(sessionID string) (err error) {
	value, err := r.client.Get(sessionID).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value), &r.Data)
	if err != nil {
		panic(err)
		return
	}
	// 2. 根据sessionID找到对应的数据
	// 3. 把数据取出来反序列化到r.data
	return
}

func (r *RedisSD) Del(key string) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	delete(r.Data, key)
}

func (r *RedisSD) Save() {
	value, err := json.Marshal(r.Data)
	if err != nil {
		fmt.Printf("marshal session data failed,err:%v\n", err)
		return
	}
	r.client.Set(r.ID, value, time.Second*time.Duration(r.expired))
}

func (r *RedisSD) SetExpire(expired int) {
	r.expired = expired
}

type RedisMgr struct {
	Session map[string]SessionData
	rwLock  sync.RWMutex
	client  *redis.Client
}

func NewRedisMgr() Mgr {
	return &RedisMgr{
		Session: make(map[string]SessionData, 1024),
	}
}
func (r *RedisMgr) Init(addr string, option ...string) {
	var (
		password string
		db       string
	)
	if len(option) == 1 {
		password = option[0]
	} else if len(option) == 2 {
		password = option[0]
		db = option[1]
	}
	dbValue, err := strconv.Atoi(db)
	if err != nil {
		dbValue = 0
	}
	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbValue,
	})
	_, err = r.client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

// GetSessionData 获取sessionID对应的sessionData
func (r *RedisMgr) GetSessionData(sessionID string) (sd SessionData, err error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	sd, ok := r.Session[sessionID]
	if !ok {
		err = fmt.Errorf("invalid session id")
		return
	}
	err = sd.Load(sessionID)
	return

}

func (r *RedisMgr) CreateSession() (sd SessionData) {
	uuidObj := uuid.NewV4()
	sd = NewRedisSessionData(uuidObj.String(), r.client)
	r.Session[sd.GetID()] = sd
	return
}
