package cache

import (
	"sync"
	"time"
)

type HashCacheInfo struct{
	Key string
	Data interface{}
	Expire int64
}

const(
	NO_REMOVE int = 0
	LAZY_REMOEV int = 1
	EACH_REMOVE int = 2
)

type HashCache struct{
	data map[string]HashCacheInfo
	mutex sync.Mutex
	expire int64
	expireAdd int64
	removeMode int
}

func (c *HashCache)Get(key string)interface{}{
	if v, ok := c.data[key]; ok{
		v.Expire += c.expireAdd
		c.data[key] = v
		switch c.removeMode {
		case NO_REMOVE:
			break
		case LAZY_REMOEV:
			c.LazyRemove(&v)
			break
		case EACH_REMOVE:
			c.EachRemove()
			break
		}
		return v.Data
	}
	return nil
}

func (c *HashCache)SaveGet(key string)interface{}{
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Get(key)
}

func (c *HashCache)Set(key string, data interface{}, expire int64){
	if v, ok := c.data[key]; ok{
		v.Data = data
		v.Expire += c.expireAdd
		c.data[key] = v
	}else{
		info := HashCacheInfo{Key: key, Data:data, Expire:expire}
		c.data[key] = info
	}
}

func (c *HashCache)SaveSet(key string, data interface{}, expire int64){
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Set(key, data, expire)
}

func (c *HashCache)Remove(key string)bool{
	if _, ok := c.data[key]; ok{
		delete(c.data, key)
		return true
	}
	return false
}

func (c *HashCache)SaveRemove(key string)bool{
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Remove(key)
}

func (c *HashCache)SetConfig(Expire int64, ExpireAdd int64){
	c.expire = Expire
	c.expireAdd = ExpireAdd
	if c.data == nil{
		c.data = make(map[string]HashCacheInfo)
	}
}

func (c *HashCache)LazyRemove(info *HashCacheInfo){
	if info.Expire >= time.Now().Unix(){
		c.Remove(info.Key)
	}
}

func (c *HashCache)EachRemove(){
	now := time.Now().Unix()
	for k, v := range c.data{
		if v.Expire >= now{
			delete(c.data, k)
		}
	}
}