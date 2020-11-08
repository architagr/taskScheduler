package taskscheduler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	all     int = iota
	running int = iota
	cancled int = iota
	stopped int = iota
	completed int = iota
)

// Task this struct has data for each task,
// task can be uniquely identified by the ID field
type Task struct {
	index       int
	TriggerTime time.Duration
	EndTime     time.Time
	StartTime   time.Time
	Data        interface{}
	Status      int
}
type timmer struct {
	timmer *time.Timer
}

// ChannelData channel data emited by the channel once the task has ended
type ChannelData struct {
	ChannelName string
	Data        interface{}
	ID          uuid.UUID
}
type channel struct {
	ch       chan ChannelData
	maxIndex int
	f        func(ch <-chan ChannelData)
	tasks    map[uuid.UUID]Task
	timmers  map[uuid.UUID]timmer
}

var channels map[string]channel = make(map[string]channel)

// AddChannel Adding a new channel to the list of all channels,
// the output indicates if the channels have been added by the name
// provided as an input, also the input function is used to send the
// reference of the channel to the function in read only mode
// so that the function can implement logic to receive data
// we have a limit of 10 message in the channel
func AddChannel(name string, f func(ch <-chan ChannelData)) (bool, error) {
	fmt.Printf("All:%d , running:%d, cancled:%d, stopped:%d, completed: %d", all, running,cancled,stopped,completed)
	if _, ok := channels[name]; ok {
		return false, fmt.Errorf("%s name already exists", name)
	}

	var ch = channel{
		tasks:    make(map[uuid.UUID]Task),
		timmers:  make(map[uuid.UUID]timmer),
		maxIndex: 0,
		f:        f,
		ch:       make(chan ChannelData),
	}
	channels[name] = ch
	return true, nil
}

// RemoveChannel removes a channel that exixts in the list of all channels,
// the output indicates if the channels have been removed by the name
// provided as an input and if all task in it have been stoped
func RemoveChannel(name string) (bool, map[uuid.UUID]Task, error) {
	if channelData, ok := channels[name]; ok {
		for key, val := range channelData.timmers {
			val.timmer.Stop()
			if taskData, okT := channelData.tasks[key]; okT {
				taskData.Status = stopped
				channelData.tasks[key] = taskData
			}
		}
		close(channelData.ch)
		delete(channels, name)
		return true, channelData.tasks, nil
	}
	return false, nil, fmt.Errorf("%s name does not exists", name)
}

// AddTask adds a task in the channel, once the task is complete (duration has passed)
// the data that is been passed in this will be send over the channel
// the data along with the channel name and task ID is been send to the
// function that was passed when creating a channel
func AddTask(channelName string, duration time.Duration, data interface{}) (bool, uuid.UUID, error) {
	if channelData, ok := channels[channelName]; ok {
		channelData.maxIndex++
		var ID = uuid.New()
		var task = Task{
			index:   channelData.maxIndex,
			Data:    data,
			EndTime: time.Now().Add(duration),
			StartTime: time.Now(),
			Status: running,
			TriggerTime: duration,
		}
		var timmer = timmer{
			timmer: createTimmer(channelName, duration, data, ID, channelData, &task),
		}
		channelData.tasks[ID] = task
		channelData.timmers[ID] = timmer

		if channelData.maxIndex == 1 {
			go channelData.f(channelData.ch)
		}
		return true, ID, nil
	}
	return false, uuid.UUID{}, fmt.Errorf("channel by name %s does not exist, first add a channel by this name", channelName)
}

func createTimmer(channelName string, duration time.Duration, data interface{}, ID uuid.UUID, channelData channel, task *Task) *time.Timer {
	return time.AfterFunc(duration, func() {
		if channelData, ok:= channels[channelName];ok{
			if task, okT:= channelData.tasks[ID]; okT{
				task.Status = completed
				channelData.tasks[ID] = task
			}
		}

		channelData.ch <- ChannelData{
			ChannelName: channelName,
			Data:        data,
			ID:          ID,
		}
	})
}

// RemoveTask this stops a task by the task id
func RemoveTask(channelName string, taskID uuid.UUID) (bool, error) {
	if channelData, ok := channels[channelName]; ok {
		if timmer, okT := channelData.timmers[taskID]; okT {
			timmer.timmer.Stop()
			if task, okTask := channelData.tasks[taskID]; okTask {
				task.Status = stopped
				channelData.tasks[taskID] = task
			}
		}
		return true, nil
	}
	return false, fmt.Errorf("channel by name %s does not exist, first add a channel by this name", channelName)
}
