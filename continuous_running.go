package continuous_daemon

//ContinuousRunning defines a struct for mantain parallel enven handleable panic
type ContinuousRunning struct {
	bufferSize uint
	parallel   uint

	msgList chan interface{}

	running RunningInterface
}

// RunningInterface defines interface about daemon
type RunningInterface interface {
	Run(msg interface{})
	Error(err interface{})
}

// NewContinuousRunning create a continuous runable daemon
func NewContinuousRunning(bufferSize, parallel uint, run RunningInterface) *ContinuousRunning {
	msgList := make(chan interface{}, bufferSize)

	c := &ContinuousRunning{
		bufferSize: bufferSize,
		parallel:   parallel,

		msgList: msgList,

		running: run,
	}

	//continuous run
	go c.run()

	return c
}

// Put push a new msg into msg list at list's tail
func (c *ContinuousRunning) Put(msg interface{}) {
	c.msgList <- msg
}

// run provide a continuous daemon
func (c *ContinuousRunning) run() {
	ch := make(chan struct{}, c.parallel)
	for i := uint(0); i < c.parallel; i++ {
		ch <- struct{}{}
	}

	for {
		<-ch
		go func() {
			defer func() {
				if err := recover(); err != nil {
					c.running.Error(err)
					ch <- struct{}{}
				}
			}()
			for msg := range c.msgList {
				c.running.Run(msg)
			}
		}()
	}
}
