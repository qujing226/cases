package main

func main() {
	FindAbb()
}

//
//func main() {
//	// 创建一个任务队列
//	taskQueue := &TaskQueue{
//		&task{ID: 1, Priority: 2, ExecTime: time.Now().Add(1 * time.Second)},
//		&task{ID: 2, Priority: 1, ExecTime: time.Now().Add(2 * time.Second)},
//		&task{ID: 3, Priority: 3, ExecTime: time.Now().Add(3 * time.Second)},
//	}
//
//	// 创建堆
//	heap.Init(taskQueue)
//
//	// 模拟任务调度
//	taskQueue.Schedule()
//}
//
