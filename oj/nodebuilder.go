// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj

import "fmt"

type NodeBuilder struct {
	stack  []Node
	starts []int
}

func (b *NodeBuilder) Reset() {
	if 0 < cap(b.stack) && 0 < len(b.stack) {
		b.stack = b.stack[:0]
		b.starts = b.starts[:0]
	} else {
		b.stack = make([]Node, 0, 64)
		b.starts = make([]int, 0, 16)
	}
}

func (b *NodeBuilder) Object(key ...string) error {
	newObj := Object{}
	if 0 < len(key) {
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return fmt.Errorf("can not use a key when pushing to an array")
		}
		if obj, _ := b.stack[len(b.stack)-1].(Object); obj != nil {
			obj[key[0]] = newObj
		}
	} else if 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0 {
		return fmt.Errorf("must have a key when pushing to an object")
	}
	b.starts = append(b.starts, -1)
	b.stack = append(b.stack, newObj)

	return nil
}

func (b *NodeBuilder) Array(key ...string) error {
	if 0 < len(key) {
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return fmt.Errorf("can not use a key when pushing to an array")
		}
		b.stack = append(b.stack, Key(key[0]))
	} else if 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0 {
		return fmt.Errorf("must have a key when pushing to an object")
	}
	b.starts = append(b.starts, len(b.stack))
	b.stack = append(b.stack, EmptyArray)

	return nil
}

func (b *NodeBuilder) Value(value Node, key ...string) error {
	if 0 < len(key) {
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return fmt.Errorf("can not use a key when pushing to an array")
		}
		if obj, _ := b.stack[len(b.stack)-1].(Object); obj != nil {
			obj[key[0]] = value
		}
	} else if 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0 {
		return fmt.Errorf("must have a key when pushing to an object")
	} else {
		b.stack = append(b.stack, value)
	}
	return nil
}

func (b *NodeBuilder) Pop() {
	if 0 < len(b.starts) {
		start := b.starts[len(b.starts)-1]
		if 0 <= start { // array
			start++
			size := len(b.stack) - start
			a := Array(make([]Node, size))
			copy(a, b.stack[start:len(b.stack)])
			b.stack = b.stack[:start]
			b.stack[start-1] = a
			if 2 < len(b.stack) {
				if k, ok := b.stack[len(b.stack)-2].(Key); ok {
					if obj, _ := b.stack[len(b.stack)-3].(Object); obj != nil {
						obj[string(k)] = a
						b.stack = b.stack[:len(b.stack)-2]
					}
				}
			}
		}
		b.starts = b.starts[:len(b.starts)-1]
	}
}

func (b *NodeBuilder) PopAll() {
	for 0 < len(b.starts) {
		b.Pop()
	}
}

func (b *NodeBuilder) Result() (result interface{}) {
	if 0 < len(b.stack) {
		result = b.stack[0]
	}
	return
}
