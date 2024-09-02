package jobs

import (
	"errors"
	"fmt"
)

type batchContext struct {
	kvs map[string]interface{}
}

// BatchContext contains properties during a job or step execution
type BatchContext struct {
	ctx batchContext
}

// NewBatchContext new instance
func NewBatchContext() *BatchContext {
	c := &BatchContext{
		ctx: batchContext{
			kvs: map[string]interface{}{},
		},
	}
	return c
}

func (ctx *BatchContext) Put(key string, value interface{}) {
	ctx.ctx.kvs[key] = value
}

func (ctx *BatchContext) Exists(key string) bool {
	val := ctx.ctx.kvs[key]
	return val != nil
}

func (ctx *BatchContext) Remove(key string) {
	delete(ctx.ctx.kvs, key)
}

func (ctx *BatchContext) Get(key string, def ...interface{}) interface{} {
	val := ctx.ctx.kvs[key]
	if val == nil && len(def) > 0 {
		val = def[0]
	}
	return val
}

func (ctx *BatchContext) GetInt(key string, def ...int) (int, error) {
	v := ctx.ctx.kvs[key]
	if v == nil && len(def) > 0 {
		return def[0], nil
	}
	if v != nil {
		switch r := v.(type) {
		case int:
			return r, nil
		case int8:
			return int(r), nil
		case int16:
			return int(r), nil
		case int32:
			return int(r), nil
		case int64:
			return int(r), nil
		case uint:
			return int(r), nil
		case uint8:
			return int(r), nil
		case uint16:
			return int(r), nil
		case uint32:
			return int(r), nil
		case uint64:
			return int(r), nil
		case float32:
			return int(r), nil
		case float64:
			return int(r), nil
		}
	}
	return 0, errors.New(fmt.Sprintf("value is nil or not int: %v", v))
}

func (ctx *BatchContext) GetInt64(key string, def ...int64) (int64, error) {
	v := ctx.ctx.kvs[key]
	if v == nil && len(def) > 0 {
		return def[0], nil
	}
	if v != nil {
		switch r := v.(type) {
		case int:
			return int64(r), nil
		case int8:
			return int64(r), nil
		case int16:
			return int64(r), nil
		case int32:
			return int64(r), nil
		case int64:
			return int64(r), nil
		case uint:
			return int64(r), nil
		case uint8:
			return int64(r), nil
		case uint16:
			return int64(r), nil
		case uint32:
			return int64(r), nil
		case uint64:
			return int64(r), nil
		case float32:
			return int64(r), nil
		case float64:
			return int64(r), nil
		}
	}
	return 0, errors.New(fmt.Sprintf("value is nil or not int64: %v", v))
}

func (ctx *BatchContext) GetString(key string, def ...string) (string, error) {
	v := ctx.ctx.kvs[key]
	if v == nil && len(def) > 0 {
		return def[0], nil
	}
	if v != nil {
		if r, ok := v.(string); ok {
			return r, nil
		}
	}
	return "", errors.New(fmt.Sprintf("value is nil or not string: %v", v))
}

func (ctx *BatchContext) GetBool(key string, def ...bool) (bool, error) {
	v := ctx.ctx.kvs[key]
	if v == nil && len(def) > 0 {
		return def[0], nil
	}
	if v != nil {
		if r, ok := v.(bool); ok {
			return r, nil
		}
	}
	return false, errors.New(fmt.Sprintf("value is nil or not bool: %v", v))
}
