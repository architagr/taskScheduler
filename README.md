# taskScheduler

This library can be used at places where you want to create tasks that have to be executed after some particular time duration.

To use this you have to 1st call the AddChannel function sending in the channel name and a function that have a read only channel as a input (to listen to this channel).

Then call a AddTask function with channel name in which channel to add the task, duration in which you want the event and data that you want to receive once the duration ends.

for example if we are configuring this in a project then -

In main.go I have created a function configureChannel() that is run as a goroutine in the main function 

func <b>configureChannel()</b> {<br/>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;taskscheduler.AddChannel("channel1", <i>receiveDataChannel1</i>)<br/>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;ok, ID, err := taskscheduler.AddTask("channel1", time.Second*time.Duration(10), "data1")<br/>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;if !ok {<br/>
		&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;fmt.Println("error adding task:",err)<br/>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;} else {<br/>
		&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;fmt.Println("no error adding task", ID)<br/>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;}<br/>
}<br/><br/>

func <b>receiveDataChannel1(ch <-chan ts.ChannelData)</b> {<br/>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;for channelTask := range ch {<br/>
		&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;fmt.Println("received data:",channelTask.Data,"for task id :", channelTask.ID)<br/>
	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;}<br/>
}<br/>

func <b>main()</b>{<br/>
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;fmt.Println("run")<br/>
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<i>go configureChannel()</i><br/>
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;time.Sleep(time.Second*time.Duration(60))<br/>
}<br/>
