package promise
import (
	"fmt"
	"testing"
	"strings"
	"errors"
)


func TestBasicFlow1(t *testing.T) {
	var sb strings.Builder
    p := NewPromise(func(resolve func(interface{}), reject func(er error)) {
    	sb.WriteString("new->")
		resolve("wow")
	})
	p.then(func(val interface{}) interface{} {
		sb.WriteString("then->")
		return val.(string) + "wow";
	}).then(func(val interface{}) interface{} {
		sb.WriteString("then->")
		return val.(string) + "wow";
	})
	result := await(p).(string)
	if sb.String() != "new->then->then->" {
		t.Errorf("Got = %v; want new->then->then->", sb.String())
	}
	if result != "wowwowwow" {
		t.Errorf("Got = %v; want wowwowwow", result)
	}
}

func TestBasicFlow2(t *testing.T) {
	var sb strings.Builder
    p := NewPromise(func(resolve func(interface{}), reject func(er error)) {
    	sb.WriteString("new->")
		resolve("wow")
	})
	p.then(func(val interface{}) interface{} {
		sb.WriteString("then->")
		return val.(string) + "wow";
	}).then(func(val interface{}) interface{} {
		sb.WriteString("then->")
		return val.(string) + "wow";
	}).catch(func(err error) interface{} {
		// This should not get called
		sb.WriteString("catch->")
		fmt.Println("MAJAK CATCH")
		return "wow";
	})
	result := await(p).(string)
	if sb.String() != "new->then->then->" {
		t.Errorf("Got = %v; want new->then->then->", sb.String())
	}
	if result != "wowwowwow" {
		t.Errorf("Got = %v; want wowwowwow", result)
	}
}


func TestCatchFlow1(t *testing.T) {
	var sb strings.Builder
    p := NewPromise(func(resolve func(interface{}), reject func(er error)) {
    	sb.WriteString("new->")
		resolve("wow")
	})
	p.then(func(val interface{}) interface{} {
		sb.WriteString("then->")
		return errors.New("ohla ohla");
	}).then(func(val interface{}) interface{} {
		// This should not get called
		sb.WriteString("then->")
		return val.(string) + "wow";
	}).catch(func(err error) interface{} {
		sb.WriteString("catch->")
		return "wow";
	})
	result := await(p)
	if sb.String() != "new->then->catch->" {
		t.Errorf("Got = %v; want new->then->catch->", sb.String())
	}
	if result != "wow" {
		t.Errorf("Got = %v; want wow", result)
	}
}

func TestCatchFlow2(t *testing.T) {
	var sb strings.Builder
    p := NewPromise(func(resolve func(interface{}), reject func(er error)) {
    	sb.WriteString("new->")
		reject(errors.New("ohla ohla"));
	})
	p.then(func(val interface{}) interface{} {
		// This should not be called
		sb.WriteString("then->")
		return errors.New("ohla ohla");
	}).then(func(val interface{}) interface{} {
		// This should not get called
		sb.WriteString("then->")
		return val.(string) + "wow";
	}).catch(func(err error) interface{} {
		sb.WriteString("catch->")
		return "wow";
	})
	result := await(p).(string)
	if sb.String() != "new->catch->" {
		t.Errorf("Got = %v; want new->catch->", sb.String())
	}
	if result != "wow" {
		t.Errorf("Got = %v; want wow", result)
	}
}

func TestCatchFlow3(t *testing.T) {
	var sb strings.Builder
    p := NewPromise(func(resolve func(interface{}), reject func(er error)) {
    	sb.WriteString("new->")
		reject(errors.New("ohla ohla"));
	})
	p.then(func(val interface{}) interface{} {
		// This should not be called
		sb.WriteString("then->")
		return errors.New("ohla ohla");
	}).then(func(val interface{}) interface{} {
		// This should not get called
		sb.WriteString("then->")
		return val.(string) + "wow";
	}).catch(func(err error) interface{} {
		sb.WriteString("catch->")
		return "wow";
	}).then(func(val interface{}) interface{} {
		sb.WriteString("then->")
		return val.(string) + "wow";
	})
	result := await(p).(string)
	if sb.String() != "new->catch->then->" {
		t.Errorf("Got = %v; want new->catch->", sb.String())
	}
	if result != "wowwow" {
		t.Errorf("Got = %v; want wowwow", result)
	}
}

func TestFinally(t *testing.T) {
	var sb strings.Builder
    p := NewPromise(func(resolve func(interface{}), reject func(er error)) {
    	sb.WriteString("new->")
		reject(errors.New("ohla ohla"));
	})
	p.then(func(val interface{}) interface{} {
		// This should not be called
		sb.WriteString("then->")
		return errors.New("ohla ohla");
	}).then(func(val interface{}) interface{} {
		// This should not get called
		sb.WriteString("then->")
		return val.(string) + "wow";
	}).catch(func(err error) interface{} {
		sb.WriteString("catch->")
		return "wow";
	}).then(func(val interface{}) interface{} {
		sb.WriteString("then->")
		return val.(string) + "wow";
	}).finally(func(val interface{}) interface{} {
		sb.WriteString("finally->")
		return val.(string) + "wow";
	})
	result := await(p).(string)
	if sb.String() != "new->catch->then->finally->" {
		t.Errorf("Got = %v; want new->catch->finally->", sb.String())
	}
	if result != "wowwow" {
		t.Errorf("Got = %v; want wowwow", result)
	}
}


