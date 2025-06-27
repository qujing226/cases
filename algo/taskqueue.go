package main

import (
	"container/heap"
	"fmt"
	"time"
)

/*
假设你现在需要开发一个简单的任务调度器，它要根据一定的规则来安排任务的执行顺序。
比如说，某些任务可能需要按优先级执行，某些任务可能需要按时间安排。

我们要处理的任务调度器问题就是：
	有一些任务需要在给定的时间点执行，每个任务有一个执行时间，并且任务可能有不同的优先级。
	我们需要设计一个调度器，确保每个任务按时执行，且符合优先级要求。

我们来想象一下，如果我们用一种简单的方式来调度任务，最直观的思路就是一个队列。
队列按顺序处理任务，简单直接，但是问题也很明显：
	如果某些任务需要优先执行，或者任务之间有时间的要求，普通的队列就不能满足需求了。
	这时候，我们就需要用到优先级队列了。优先级队列的特点是，每次从队列中取出的元素是当前优先级最高的任务。
*/

type task struct {
	ID       int
	Priority int       // 优先级
	ExecTime time.Time // 执行时间
}

type TaskQueue []*task

func (tq *TaskQueue) Len() int { return len(*tq) }
func (tq *TaskQueue) Less(i, j int) bool {
	// 根据优先级决定执行顺序
	t := *tq
	if t[i].Priority == t[j].Priority {
		// 如果优先级相同，则根据执行时间决定执行顺序
		return t[i].ExecTime.Before(t[j].ExecTime)
	}
	return t[i].Priority > t[j].Priority
}
func (tq *TaskQueue) Swap(i, j int) {
	t := *tq
	t[i], t[j] = t[j], t[i]
}

func (tq *TaskQueue) Push(x any) {
	*tq = append(*tq, x.(*task))
}

func (tq *TaskQueue) Pop() any {
	old := *tq
	n := len(old)
	item := old[n-1]
	*tq = old[:n-1]
	return item
}

func (tq *TaskQueue) Schedule() {
	for tq.Len() > 0 {
		task := heap.Pop(tq).(*task)
		// 执行任务
		time.Sleep(time.Until(task.ExecTime))
		fmt.Printf("Task %d executed at %v\n", task.ID, task.ExecTime)
	}
}
func main() {
	// 创建一个任务队列
	taskQueue := &TaskQueue{
		&task{ID: 1, Priority: 2, ExecTime: time.Now().Add(1 * time.Second)},
		&task{ID: 2, Priority: 1, ExecTime: time.Now().Add(2 * time.Second)},
		&task{ID: 3, Priority: 3, ExecTime: time.Now().Add(3 * time.Second)},
	}

	// 创建堆
	heap.Init(taskQueue)

	// 模拟任务调度
	taskQueue.Schedule()
}
