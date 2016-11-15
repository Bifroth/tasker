package tasker

//Function is a function that can be executed by an Executor
type Function interface {
	Initialize()
	Execute(params []interface{}) interface{}
	Close()
}

type task struct {
	close bool
	args  []interface{}
	res   interface{}
	out   chan task
}

//Executor executes functions
type Executor struct {
	in    chan task
	count int
}

//Initialize inits the Executor (Channels)
func (e *Executor) Initialize() {
	e.in = make(chan task)
}

// AddFunction adds a Function to the executed functions.
// Don't add different Functions, or the results will become quite.. funny
func (e *Executor) AddFunction(fn Function) {
	go run(e.in, fn)
	e.count++
}

//CloseOne dismisses one Function instance
func (e *Executor) CloseOne() {
	e.executeTask(e.in, task{close: true})
	e.count--
}

//Execute the function with the given arguments
func (e *Executor) Execute(args ...interface{}) interface{} {
	var t = task{args: args}
	t = e.executeTask(e.in, t)
	return t.res
}

func run(in <-chan task, fn Function) {
	fn.Initialize()

	t, cont := <-in //read first

	for cont && !t.close { //check if channel is closed or the current task is to close the current goroutine
		t.res = fn.Execute(t.args) // Execute action
		t.out <- t                 // return result

		t, cont = <-in //next
	}

	t.out <- t

	//this goroutine is done, close the function
	fn.Close()
}

func (e *Executor) executeTask(in chan<- task, t task) task {
	out := make(chan task)
	t.out = out
	in <- t

	var tdone = <-out
	close(out)
	return tdone
}

//ThreadCount gets the number of currently running goroutines on this Executor
func (e Executor) ThreadCount() int {
	return e.count
}

//ChanSize the number of currently waiting tasks
func (e Executor) ChanSize() int {
	return len(e.in)
}

//TODO Auto open and close functions based on channel size
