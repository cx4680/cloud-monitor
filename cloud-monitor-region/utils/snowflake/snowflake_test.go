package snowflake

import (
	"fmt"
	"testing"
)

func TestGetWorker(t *testing.T) {
	worker := GetWorker()
	ch := make(chan int64)
	count := 100
	// 并发 goroutine ID生成
	for i := 0; i < count; i++ {
		go func() {
			id := worker.nextId()
			ch <- id
		}()
	}
	defer close(ch)
	m := make(map[int64]int)
	for i := 0; i < count; i++ {
		id := <-ch
		// map中存在为id的key,说明生成的 ID有重复
		_, ok := m[id]
		if ok {
			fmt.Println("ID is not unique!")
		}
		// id作为key存入map
		m[id] = i
		fmt.Println(id)
		fmt.Println(convertToBin(id))
	}
}
