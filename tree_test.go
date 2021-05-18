package tree

import (
	"sync"
	"testing"
)

func TestInsert(t *testing.T) {
	tree := NewTree("test")
	var wg sync.WaitGroup
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() {
			tree.Insert(NewPeer())
			wg.Done()
		}()
	}
	wg.Wait()
	tree.DumpTree(nil)
	tree.Delete(tree.Find(1))
	tree.Delete(tree.Find(2))
	tree.Delete(tree.Find(3))
	tree.DumpTree(nil)
	//p := tree.Find(0)
	//fmt.Println(p.id, p.parent.GetAddr())
}
