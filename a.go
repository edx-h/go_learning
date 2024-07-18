package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	stack := make([]*tree.Tree, 10)

	//traverse the tree by DFS, until the stack is empty
	CurrentIndex := 0
	CurrentNode := t
	stack[CurrentIndex] = t

	for {
		//time.Sleep(5 * time.Second)
		if CurrentIndex >= 0 { // if stack is not null, then there are elements not traversed yet.
			CurrentNode = stack[CurrentIndex]
			fmt.Printf("Current index is: %d, value is: %d, channel address is: %p\n", CurrentIndex, CurrentNode.Value, &ch)

			if CurrentNode.Left != nil && CurrentNode.Left.Value != -1 {
				fmt.Printf("Current value is %d, Insert left value %d to channel %p\n", CurrentNode.Value, CurrentNode.Left.Value, &ch)
				stack[CurrentIndex+1] = CurrentNode.Left
				CurrentIndex += 1
				CurrentNode = CurrentNode.Left
			} else {
				fmt.Printf("No more left node! Channel address is: %p\n", &ch)
				fmt.Printf("Index %d, Pop out value %d for channel %p\n", CurrentIndex, stack[CurrentIndex].Value, &ch)
				ch <- stack[CurrentIndex].Value
				stack[CurrentIndex] = nil
				fmt.Println("CurrentNode value is ", CurrentNode.Value, "for channel", &ch)
				CurrentNode.Value = -1

				if CurrentNode.Right != nil {
					fmt.Printf("Index %d, Current value is %d, Insert right value %d to channel %p\n", CurrentIndex, CurrentNode.Value, CurrentNode.Right.Value, &ch)
					stack[CurrentIndex] = CurrentNode.Right
					CurrentNode = CurrentNode.Right
				} else {
					fmt.Printf("Index %d, Current value %d is finished for channel %p\n", CurrentIndex, CurrentNode.Value, &ch)
					CurrentIndex -= 1
				}
			}
		} else {
			//already empty
			close(ch)
			return
		}
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		v1, ok1 := <-ch1
		v2 := <-ch2
		//fmt.Println("ITERATION ", i, ok1, ok2)
		if !ok1 {
			break
		}
		if v1 != v2 {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	//ch1 := make(chan int)
	//Walk(tree.New(1), ch1)
	// ch := make(chan int)
	// go Walk(tree.New(2), ch)
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(<-ch)
	// }
}
