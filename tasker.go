package tasker

//Function is a function that can be executed by an Executor
type Function interface {
	Initialize()
	Execute(params []interface{}) interface{}
	Close()
}

type task struct {
	id    int64
	close bool
	args  []interface{}
	res   interface{}
}

//Executor executes functions
type Executor struct {
	taskID int64
	out    chan task
	in     chan task
}

//Initialize inits the Executor (Channels)
func (e *Executor) Initialize() {
	e.out = make(chan task)
	e.in = make(chan task)
}

// AddFunction adds a Function to the executed functions.
// Don't add different Functions, or the results will become quite.. funny
func (e Executor) AddFunction(fn Function) {
	go run(e.in, e.out, fn)
}

//Execute the function with the given arguments
func (e *Executor) Execute(args ...interface{}) interface{} {
	var t = task{args: args}

	t = e.executeTask(e.in, e.out, t)

	return t.res
}

func run(in <-chan task, out chan<- task, fn Function) {
	fn.Initialize()

	t, cont := <-in //read first

	for cont && !t.close { //check if channel is closed or the current task is to close the current goroutine
		t.res = fn.Execute(t.args) // Execute action
		out <- t                   // return result

		t, cont = <-in //next
	}

	//this goroutine is done, close the function
	fn.Close()
}

func (e *Executor) executeTask(in chan<- task, out <-chan task, t task) task {
	e.taskID++
	t.id = e.taskID
	in <- t

	var tdone = <-out
	for tdone.id != t.id {
		tdone = <-out
	}
	return tdone
}

//TODO Auto open and close functions based on channel size
