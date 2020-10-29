package taskScheduler

import (
	"fmt"
	"time"
)

// Task sdad
type Task struct {
	T     *time.Timer
	index int
}
var maxTaskIndex int = 0
// TaskData sddf
type TaskData struct{
	Data  interface{}
	CallBack func(data interface{}, index int)
}

var timmers chan Task
var timeTask map[int]TaskData

// Init to initilize the schedule
func Init(calback func()){
	timmers = make(chan Task)
	timeTask = make(map[int]TaskData)
	go timmerChannleRead(calback)
}

// AddTask asda
func AddTask(newTask Task, taskdata TaskData){
	maxTaskIndex++
	newTask.index = maxTaskIndex
	timmers <- newTask
	timeTask[maxTaskIndex] = taskdata
}

// Close asd
func Close(){
	close(timmers)
}

func timmerChannleRead(calback func()){
	for channelTask := range <-timmers{
		go func(task Task) {
			<-task.T.C
			//fmt.Println("task : ",task.time)
			//WriteFile(f, task.time)
			if taskdata, ok := timeTask[task.index]; ok{
				taskdata.CallBack(taskdata.Data, task.index)
			}
			//DeleteTask(task)
		}(channelTask)
	}
	defer calback()
	// for {
	// 	select {
	// 	case channelTask := <-timmers:
	// 		go func(task Task) {
	// 			<-task.T.C
	// 			//fmt.Println("task : ",task.time)
	// 			//WriteFile(f, task.time)
	// 			if taskdata, ok := timeTask[task.index]; ok{
	// 				taskdata.CallBack(taskdata.Data, task.index)
	// 			}
				
	// 			//DeleteTask(task)
	// 		}(channelTask)
	// 	}
	// }
}

// DeleteTask sd
func DeleteTask(index int){
	if _, ok := timeTask[index]; ok {
		delete(timeTask, index)
	}
	fmt.Println("delete :", timeTask)
}