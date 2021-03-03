package promise
import "fmt"

const (
	PENDING = iota
	SUCCESS
	REJECT
)
type Promise struct {
	state int
	exec func(resolve func(interface{}), reject func(error))
	resolveValue chan interface{}
	rejectValue chan error
	// Needed to implement await
	last bool
	done chan interface{}
	chain *Promise
	result interface{}
}

func (p *Promise) resolve(val interface{}) {
	if p.state != PENDING {
		return
	}
	p.state = SUCCESS
	p.result = val
	p.resolveValue <-val
}

func (p *Promise) reject(e error) {
	if p.state != PENDING {
		return
	}

	p.state = REJECT
	p.result = e
	p.rejectValue <- e
}

func NewPromise(exec func(resolve func(interface{}), reject func(error))) *Promise {
	promise:= &Promise {
		state: PENDING,
		exec: exec,
		resolveValue: make(chan interface {}, 1),
		rejectValue: make(chan error, 1),
		last: true,
		done: make(chan interface{}, 1),
		chain: nil,
		result: nil,
	}
	go exec(promise.resolve, promise.reject)
	return promise
}

func (p *Promise) then(then_func func(val interface{}) interface{}) *Promise {
	var result *Promise
	result = NewPromise(func(resolve func(interface{}), reject func(error)) {
		select {
		case val := <- p.resolveValue:
			func() {
				defer func() {
					p.done <- true
				} ()
				p.chain = result
				response := then_func(val)
				err, _ := response.(error)
				if (err != nil) {
					reject(err)
				} else {
					resolve(response.(string))
				}
			} ()
		case val := <- p.rejectValue:
			func() {
				defer func() {
					p.done <- true
				} ()
				p.chain = result
				reject(val)
			} ()
		}	
	})
	p.last = false;
	return result
}

func (p *Promise) catch(catch_func func(err error) interface{}) *Promise {
	var result *Promise
	result = NewPromise(func(resolve func(interface{}), reject func(error)) {
		select {
		case val := <- p.rejectValue:
			func() {
				defer func() {
					p.done <- true
				} ()
				p.chain = result
				response := catch_func(val)
				err, _ := response.(error)
				if (err != nil) {
					reject(err)
				} else {
					resolve(response)
				}
			} ()
		case val := <- p.resolveValue:
			func() {
				defer func() {
					p.done <- true
				} ()
				p.chain = result
				resolve(val)
			} ()
		}
	})
	p.last = false;
	return result
	
}


func (p *Promise) finally(final_func func(interface{}) interface{}) *Promise {
	var result *Promise
	result = NewPromise(func(resolve func(interface{}), reject func(error)) {
		select {
		case val := <- p.rejectValue:
			func() {
				defer func() {
					p.done <- true
				} ()
				p.chain = result
				result.result = val
				final_func(val)
			} ()
		case val := <- p.resolveValue:
			func() {
				defer func() {
					p.done <- true
				} ()
				p.chain = result
				result.result = val
				final_func(val)
			} ()
		}
	})
	p.last = false;
	return result
	
}

func await(p *Promise) interface{} {
	var result interface{}
	for p != nil && !p.last {
		select{
		case <- p.done:
			p = p.chain 
			result = p.result
		}
	}
	return result
}

func main() {
	p := NewPromise(func(resolve func(interface{}), reject func(er error)) {
		resolve("wow")
	})
	p.then(func(val interface{}) interface{} {
		fmt.Println("MAJAK")
		return val.(string) + "wow";
	}).then(func(val interface{}) interface{} {
		fmt.Println("MAJAK")
		return val.(string) + "wow";
	})
	await(p)
}