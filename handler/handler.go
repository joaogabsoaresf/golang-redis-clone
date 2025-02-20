package handler

import (
	"strconv"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"COMMAND": command,
	"TEST":    command,
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"HLEN":    hlen,
	"DEL":     del,
	"EXISTS":  exists,
	"HDEL":    hdel,
	"HEXISTS": hexists,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func command(args []Value) Value {
	return Value{Typ: "string", Str: "Connected"}
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "string", Str: "PONG"}
	}

	return Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{Typ: "string", Str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{Typ: "string", Str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{Typ: "null"}
	}

	values := []Value{}
	for k, v := range value {
		values = append(values, Value{Typ: "bulk", Bulk: k})
		values = append(values, Value{Typ: "bulk", Bulk: v})
	}

	return Value{Typ: "array", Array: values}
}

func hlen(args []Value) Value {
	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hlen' command"}
	}

	hash := args[0].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{Typ: "string", Str: "0"}
	}

	return Value{Typ: "string", Str: strconv.Itoa(len(value))}
}

func del(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'del' command"}
	}

	deleted := 0

	SETsMu.Lock()
	HSETsMu.Lock()

	for _, arg := range args {
		key := arg.Bulk

		if _, ok := SETs[key]; ok {
			delete(SETs, key)
			deleted++
		}

		if _, ok := HSETs[key]; ok {
			delete(HSETs, key)
			deleted++
		}
	}

	HSETsMu.Unlock()
	SETsMu.Unlock()

	return Value{Typ: "string", Str: strconv.Itoa(deleted)}
}

func exists(args []Value) Value {
	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'exists' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	_, inSETs := SETs[key]
	SETsMu.RUnlock()

	HSETsMu.RLock()
	_, inHSETs := HSETs[key]
	HSETsMu.RUnlock()

	if inSETs || inHSETs {
		return Value{Typ: "string", Str: "1"}
	}
	return Value{Typ: "string", Str: "0"}
}

func hdel(args []Value) Value {
	if len(args) < 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hdel' command"}
	}

	hash := args[0].Bulk
	keys := args[1:]

	HSETsMu.Lock()
	defer HSETsMu.Unlock()

	hashMap, ok := HSETs[hash]
	if !ok {
		return Value{Typ: "string", Str: "0"}
	}

	deleted := 0
	for _, key := range keys {
		if _, exists := hashMap[key.Bulk]; exists {
			delete(hashMap, key.Bulk)
			deleted++
		}
	}

	if len(hashMap) == 0 {
		delete(HSETs, hash)
	}

	return Value{Typ: "string", Str: strconv.Itoa(deleted)}
}

func hexists(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hexists' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	defer HSETsMu.RUnlock()

	if hashMap, ok := HSETs[hash]; ok {
		if _, exists := hashMap[key]; exists {
			return Value{Typ: "string", Str: "1"}
		}
	}

	return Value{Typ: "string", Str: "0"}
}
