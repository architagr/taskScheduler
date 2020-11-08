# taskScheduler

This library cna be used at places where you want to create task that have to be executed after some perticular time duration.

To use this you have to 1st call the AddChannel function sending in the channel name and a function that have a read only channel as a input (to listen to this channel).

Then call a AddTask function with channel name in which channel to add the task, duration in which you want the event and data that you want to receive once the duration ends.

for example if we are configuring this in a project then -

In main.go I have created a function configureChannel() that is run as a goroutine in the main function 

func configureChannel() {
	taskscheduler.AddChannel("channel1", receiveDataChannel1)
	ok, ID, err := taskscheduler.AddTask("channel1", time.Second*time.Duration(10), "data1")
	if !ok {
		fmt.Println("error adding task:",err)
	} else {
		fmt.Println("no error adding task", ID)
	}
}

func receiveDataChannel1(ch <-chan ts.ChannelData) {
	for channelTask := range ch {
		fmt.Println("received data:",channelTask.Data,"for task id :", channelTask.ID)
	}
}

func main(){
    fmt.Println("run")
}