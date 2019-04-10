package go_concurrency

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConfuse(t *testing.T) {
	for i := 0; i < 1000; i++ {
		var c Confuse
		fmt.Println(i)
		confuse(i, &c)
	}
}

type Confuse struct {
	A string
	B string
	C string
}

func confuse(i int, c *Confuse) {
	var resultMap sync.Map
	keyset := make(map[string]bool)
	keyset["a"] = true
	keyset["b"] = true
	keyset["c"] = true

	extendInfos := make(map[string]interface{}, 0)
	extendFlags := make([]string, 0)

	go func() {
		c.A = "aaa"
		extendInfos["a"] = "aaa"
		extendFlags = append(extendFlags, "aaa")
		resultMap.Store("a", true)
	}()

	go func() {
		c.B = "bbb"
		extendInfos["b"] = "bbb"
		extendFlags = append(extendFlags, "bbb")
		resultMap.Store("b", true)
	}()

	go func() {
		c.C = "ccc"
		extendInfos["c"] = "ccc"
		extendFlags = append(extendFlags, "ccc")
		resultMap.Store("c", true)
	}()

	waitResults(keyset, &resultMap, 2*time.Second)
	if c.A != "aaa" || c.B != "bbb" || c.C != "ccc" {
		panic("confuse")
	}
	fmt.Println(extendFlags)
	if len(extendFlags) != 3 {
		panic(fmt.Sprintf("%d extendFlags", i))
	}

	return
}

func waitResults(keyset map[string]bool, resultmap *sync.Map, timeOut time.Duration) (istimeout bool, err error) {
	after := time.After(timeOut)
	for {
		select {
		case <-after:
			// 超时跳出
			return true, nil
		default:
			for k := range keyset {
				if v, ok := resultmap.Load(k); ok {
					switch t := v.(type) {
					case bool:
						delete(keyset, k)
					case error:
						// 错误跳出
						return false, t
					}
				}
			}

			// 正常跳出
			if len(keyset) == 0 {
				return false, nil
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}
