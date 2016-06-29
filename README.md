# tasker
A minimal concurrency helper for golang

tasker makes threadsafe calls to functions easier.

Advantages:

- call a function and get a result while it is being processed on another goroutine
- add multiple functions to execute Requests parallel

Disadvantages:

- need to cast arguments in function
- need to cast result

Work to be done:

- add auto opening / closing of Functions for maximum efficiency

to call functions threadsafe, you must:
1. put the function into a struct. Methods:
  - Initialize: set up your db connection or whatever
  - Close: close your db connection or whatever
  - Execute: do your task (gets parameters) and return a value
2. Create and Initialize a `tasker.Executor` Object
3. Add one or more Instances of your Function to that Executor
4. call ex.Execute from any goroutine with any arguments
5. ex.Execute returns the Result of your Function.Execute()
3. ???
4. Profit!


Example:

    type CalcFunction struct { //needs to implement tasker.Function
    }

    func (c CalcFunction) Initialize() { //is called before this instance is used

    }

    func (c CalcFunction) Execute(args []interface{}) interface{} { // executes a single task
        // Gets alll parameters as a []interface{} and returns a interface{}
        //downside: you have to cast all parameters (at least for now, i need to improve that)
    	method, _ := args[0].(string)
    	num1, _ := args[1].(int)
    	num2, _ := args[2].(int)
    	switch method {
    	case "+":
    		return num1 + num2
    	case "*":
    		return num1 * num2
    	}
    	return 0
    }

    func (c CalcFunction) Close() { //is executed after the last task executed by this Function struct

    }

    var ex tasker.Executor

    func main() {
    	ex.Initialize() //must be called first

    	ex.AddFunction(CalcFunction{}) //add executors. could be 2 db connections or whatever
    	ex.AddFunction(CalcFunction{})

    	go run("+", 1, 2) //run from different goroutines (not necessary, of course)
    	go run("*", 1, 2)
    	go run("+", 10, 20)
    	go run("*", 10, 20)

        //wait for exit of those go routines. (quick and dirty, don't do that)
    	time.Sleep(1 * time.Second)
    }

    func run(operator string, num1 int, num2 int) {
        // to create, execute and retrieve the result of a task, just call Executor.Execute()
    	fmt.Println(ex.Execute(operator, num1, num2))
    }

## how does it work?

tasker creates a goroutine for every Function instance passed to it. This goroutine
initializes the function and listens for tasks on a channel. Upon receiving a task,
it calls `Function.Execute()` with the given arguments and sends the task with the set
result to another channel.

`Executor.Execute()` creates a task and sends it down the channel where it is picked up.
It then listens for done tasks until it found the right one (tasks have an id) and returns
that tasks result.

to dismiss one Function, a closing task is sent down the channel. This triggers the goroutine which receives it to exit, thus calling `Function.close`
