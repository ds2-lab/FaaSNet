package tree

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wangaoone/FaaSNet/logger"
	"github.com/wangaoone/FaaSNet/util"
)

const (
	LeftDir  = 0
	RightDir = 1
	NULL     = -1
)

var (
	MaxLength = 20 // 2 ^ 20 = 1048576 peers are enough for now
	log       = &logger.ColorLogger{
		Level: logger.LOG_LEVEL_ALL,
	}
)

type Tree struct {
	funcName string
	Root     *Peer
	mu       sync.Mutex
	length   int64

	// Record last modified ts.
	lastModified int64
}

func NewTree(s string) *Tree {
	log.Info("Create New tree %v", s)
	return &Tree{
		funcName:     s,
		lastModified: time.Now().Unix(),
	}
}

func (t *Tree) updateTs() {
	t.lastModified = time.Now().Unix()
}

func (t *Tree) GetLen() int {
	return t.len()
}

func (t *Tree) len() int {
	return int(atomic.LoadInt64(&t.length))
}

func (t *Tree) getIdx() int {
	return t.len()
}

func (t *Tree) nextAvailableNode() (*Peer, int, bool) {
	root := t.Root
	queue := []*Peer{root}
	level := 1
	for len(queue) > 0 {
		l := len(queue)

		// Reach MaxLevel, return
		if level == MaxLength {
			log.Error("Reach tree's maximum height")
			return nil, NULL, false
		}
		for i := 0; i < l; i++ {
			curNode := queue[i]
			if curNode.leftChild == nil {
				return curNode, LeftDir, true
			} else {
				queue = append(queue, curNode.leftChild)
			}
			if curNode.rightChild == nil {
				return curNode, RightDir, true
			} else {
				queue = append(queue, curNode.rightChild)
			}
		}
		level += 1
		queue = queue[l:]
	}

	// Should never reach this statement, no available Peer could be Insert
	log.Error("Should never reach this")
	return nil, NULL, false
}

func (t *Tree) Insert(appendNode *Peer) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	defer t.updateTs()
	appendNode.id = t.getIdx()

	if t.Root == nil {
		log.Info("Insert to Root")
		//log.Info("Insert to Root")
		t.Root = appendNode
		appendNode.parent = nil
	} else {
		targetNode, direction, ok := t.nextAvailableNode()
		if !ok {
			return errors.New("tree is full")
		}
		switch direction {
		case LeftDir:
			log.Info("Insert to Left")
			targetNode.leftChild = appendNode
			appendNode.parent = targetNode

		case RightDir:
			log.Info("Insert to Right")
			targetNode.rightChild = appendNode
			appendNode.parent = targetNode
		}
	}

	// Update Peer height
	p := appendNode.parent
	for p != nil {
		p.height = 1 + util.Max(p.leftChild.GetHeight(), p.rightChild.GetHeight())
		p = p.parent
	}

	// Tree length++
	atomic.AddInt64(&t.length, 1)
	return nil
}

// Only suit for Peer which has both Left and right child
// Find first leaf Peer
func (t *Tree) getSuccessor(n *Peer) *Peer {
	// start BFS
	queue := []*Peer{n}
	for len(queue) > 0 {
		l := len(queue)
		for i := 0; i < l; i++ {
			curNode := queue[i]
			if curNode.leftChild == nil && curNode.rightChild == nil {
				return curNode
			}
			if curNode.leftChild != nil {
				queue = append(queue, curNode.leftChild)
			}
			if curNode.rightChild != nil {
				queue = append(queue, curNode.rightChild)
			}
		}
		queue = queue[l:]
	}
	return nil
}

// left rotate
func (t *Tree) leftRotate(x *Peer) {
	y := x.rightChild
	x.rightChild = y.leftChild

	if y.leftChild != nil {
		y.leftChild.parent = x
	}
	y.parent = x.parent
	if x.parent == nil { // x is Root Peer
		t.Root = y
	} else if x == x.parent.leftChild { // x is left child of subtree
		x.parent.leftChild = y
	} else { // x is right child of subtree
		x.parent.rightChild = y
	}
	y.leftChild = x
	x.parent = y

	// update height
	x.height = 1 + util.Max(x.leftChild.GetHeight(), x.rightChild.GetHeight())
	y.height = 1 + util.Max(y.leftChild.GetHeight(), y.rightChild.GetHeight())
}

// right rotate
func (t *Tree) rightRotate(x *Peer) {
	y := x.leftChild
	x.leftChild = y.rightChild

	if y.rightChild != nil {
		y.rightChild.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		t.Root = y
	} else if x == x.parent.rightChild {
		x.parent.rightChild = y
	} else {
		x.parent.leftChild = y
	}
	y.rightChild = x
	x.parent = y

	// update height
	x.height = 1 + util.Max(x.leftChild.GetHeight(), x.rightChild.GetHeight())
	y.height = 1 + util.Max(y.leftChild.GetHeight(), y.rightChild.GetHeight())
}

// Multiple cases need to be handled.
func (t *Tree) Delete(deleteNode *Peer) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	defer t.updateTs()

	if deleteNode == nil {
		return errors.New("delete node should not be nil")
	}

	// Record RootNode and ParentNode of DeleteNode
	// FIXME: when and why the parent is nil
	parent := deleteNode.parent
	root := t.Root

	// Case 0:
	// Delete Peer is root Peer and Root is the only Peer in this tree.
	if deleteNode == root && root.leftChild == nil && root.rightChild == nil {
		log.Info("Case0")
		t.Root = nil
		return nil
	}

	// Case 1:
	// Delete Peer is leaf Peer, also need to check balance, if L & R have been deleted together
	//          1
	//         / \
	//        2  9
	//       /\  /
	//      3 4 10    Delete Peer 10 will incur in-balancing
	//     /\ /\
	//    5 6 7 8
	if deleteNode != root && deleteNode.leftChild == nil && deleteNode.rightChild == nil {
		log.Info("Case1")
		if parent.leftChild == deleteNode {
			parent.leftChild = nil
		} else {
			parent.rightChild = nil
		}
	}

	// Case 2.1 & 2.2:
	// Delete Peer only has left child, transplant
	// Delete Peer only has right child, transplant
	if deleteNode.leftChild != nil && deleteNode.rightChild == nil {
		log.Info("Case2.1")
		if deleteNode == root {
			t.Root = deleteNode.leftChild
			return nil
		} else {
			l := deleteNode.leftChild
			l.parent = parent
			if deleteNode == parent.leftChild {
				parent.leftChild = l
			} else if deleteNode == parent.rightChild {
				parent.rightChild = l
			}
		}
	} else if deleteNode.leftChild == nil && deleteNode.rightChild != nil {
		log.Info("Case2.2")
		if deleteNode == root {
			t.Root = deleteNode.rightChild
			return nil
		} else {
			r := deleteNode.rightChild
			r.parent = parent
			if deleteNode == parent.leftChild {
				parent.leftChild = r
			} else if deleteNode == parent.rightChild {
				parent.rightChild = r
			}
		}
	}

	// Case 3:
	// Delete Peer has left and right child
	//            1
	//          /   \
	//         2     9
	//       /  \    / \
	//      3    4  10  11   Delete Peer 9 will incur partial in-balancing
	//     / \  / \     / \
	//    5  6  7  8   15  12
	//   / \
	//  13 14
	if deleteNode.leftChild != nil && deleteNode.rightChild != nil {
		log.Info("Case3")
		l, r := deleteNode.leftChild, deleteNode.rightChild
		s := t.getSuccessor(deleteNode)
		tmpParent := s.parent

		// Cut successor from its parent
		if s == s.parent.leftChild {
			s.parent.leftChild = nil
		} else if s == s.parent.rightChild {
			s.parent.rightChild = nil
		}

		// If s is the direct child of deleteNode
		if s == l {
			l = nil
		} else if s == r {
			r = nil
		}

		// Update connection between successor and parent
		// If delete Peer is root, set its parent to nil
		if deleteNode == root {
			s.parent = nil
		} else {
			if deleteNode == deleteNode.parent.leftChild {
				deleteNode.parent.leftChild = s
			} else {
				deleteNode.parent.rightChild = s
			}
			s.parent = deleteNode.parent
		}

		// Update connection between successor and child
		s.leftChild = l
		if l != nil {
			l.parent = s
		}
		s.rightChild = r
		if r != nil {
			r.parent = s
		}

		// Update correct parent for checking balance factor
		if tmpParent == deleteNode {
			parent = s
		} else {
			parent = tmpParent
		}

		// If delete Peer is root, set new root to the tree
		if deleteNode == root {
			t.Root = s
		}
	}

	// Loop for updating balance factor, from bottom-up
	for parent != nil {
		//log.Info("checked parent is ", parent.id)
		parent.height = 1 + util.Max(parent.leftChild.GetHeight(), parent.rightChild.GetHeight())
		if parent.GetBalanceFactor() > 1 && parent.leftChild.GetBalanceFactor() >= 0 {
			log.Info("R")
			t.rightRotate(parent)
		} else if parent.GetBalanceFactor() < -1 && parent.rightChild.GetBalanceFactor() <= 0 {
			log.Info("L")
			t.leftRotate(parent)
		} else if parent.GetBalanceFactor() > 1 && parent.leftChild.GetBalanceFactor() < 0 {
			log.Info("LR")
			t.leftRotate(parent.leftChild)
			t.rightRotate(parent)
		} else if parent.GetBalanceFactor() < -1 && parent.rightChild.GetBalanceFactor() > 0 {
			log.Info("RL")
			t.rightRotate(parent.rightChild)
			t.leftRotate(parent)
		}
		parent = parent.parent
	}
	return nil
}

func (t *Tree) Find(id int) *Peer {
	t.mu.Lock()
	defer t.mu.Unlock()
	root := t.Root
	if root == nil {
		log.Info("no such Peer %v", id)
		return nil
	}

	queue := []*Peer{root}
	for len(queue) > 0 {
		l := len(queue)
		for i := 0; i < l; i++ {
			curNode := queue[i]
			if curNode.id == id {
				return curNode
			}
			if curNode.leftChild != nil {
				queue = append(queue, curNode.leftChild)
			}
			if curNode.rightChild != nil {
				queue = append(queue, curNode.rightChild)
			}
		}
		queue = queue[l:]
	}
	log.Info("no such Peer %v", id)
	return nil
}

// Print Tree in level order
func (t *Tree) DumpTree(ctx context.Context) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Root == nil {
		log.Info("Tree is %v, Root is nil", t.funcName)
		return
	}

	var res [][]int
	queue := []*Peer{t.Root}
	tempLen := 0
	for len(queue) > 0 {
		l := len(queue)
		var temp []int
		for i := 0; i < l; i++ {
			curNode := queue[i]
			tempLen++
			temp = append(temp, curNode.id)
			if curNode.leftChild != nil {
				queue = append(queue, curNode.leftChild)
			}
			if curNode.rightChild != nil {
				queue = append(queue, curNode.rightChild)
			}
		}
		res = append(res, temp)
		queue = queue[l:]
	}
	log.Info("Tree is %v, Tree length is %v, level %v", t.funcName, tempLen, res)
}

func (t *Tree) getRoot() *Peer {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.Root
}

func (t *Tree) GetFuncName() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.funcName
}

func (t *Tree) GetLatestTs() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastModified
}
