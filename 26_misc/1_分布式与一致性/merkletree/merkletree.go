package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"log"
	"math/big"
)

// TestContent implements the Content interface provided by merkletree and represents the content stored in the tree.
type TestContent struct {
	x string
}

// CalculateHash hashes the values of a TestContent
func (t TestContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t TestContent) Equals(other Content) (bool, error) {
	otherTC, ok := other.(TestContent)
	if !ok {
		return false, errors.New("value is not of type TestContent")
	}
	return t.x == otherTC.x, nil
}

func main() {
	//Build list of Content to build tree
	var list []Content
	list = append(list, TestContent{x: "Hello"})
	list = append(list, TestContent{x: "Hi"})
	list = append(list, TestContent{x: "Hey"})
	list = append(list, TestContent{x: "Hola"})

	//Create a new Merkle Tree from the list of Content
	t, err := NewTree(list)
	if err != nil {
		log.Fatal(err)
	}

	//Get the Merkle Root of the tree
	mr := t.MerkleRoot()
	log.Println(mr)

	//Verify the entire tree (hashes for each node) is valid
	vt, err := t.VerifyTree()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Verify Tree: ", vt)

	//Verify a specific content in in the tree
	vc, err := t.VerifyContent(list[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Verify Content: ", vc)

	//String representation
	log.Println(t)
}

// Content represents the data that is stored and verified by the tree. A type that
// implements this interface can be used as an item in the tree.
type Content interface {
	CalculateHash() ([]byte, error)
	Equals(other Content) (bool, error)
}

// MerkleTree is the container for the tree. It holds a pointer to the root of the tree,
// a list of pointers to the leaf nodes, and the merkle root.
type MerkleTree struct {
	Root         *Node            // 根节点
	merkleRoot   []byte           // 根节点Hash
	Leafs        []*Node          // 所有叶子节点的列表
	hashStrategy func() hash.Hash // 用于哈希的函数(缺省sha256)
	sort         bool             // 是否对左右子节点做排序后再拼接, 用于兼容 OpenZeppelin 的排序
}

// Node represents a node, root, or leaf in the tree. It stores pointers to its immediate
// relationships, a hash, the content stored if it is a leaf, and other metadata.
type Node struct {
	Tree   *MerkleTree
	Parent *Node
	Left   *Node
	Right  *Node
	leaf   bool
	dup    bool // 是否是由于奇数个叶子时, 复制了最后一个叶子
	Hash   []byte
	C      Content // 只有leaf节点才会存实际Content
	sort   bool    // 是否需要sortAppend
}

// sortAppend sort and append the nodes to be compatible with OpenZepplin libraries
// https://github.com/OpenZeppelin/openzeppelin-contracts-ethereum-package/blob/master/contracts/cryptography/MerkleProof.sol
func sortAppend(sort bool, a, b []byte) []byte {
	if !sort {
		return append(a, b...)
	}
	var aBig, bBig big.Int
	aBig.SetBytes(a)
	bBig.SetBytes(b)
	if aBig.Cmp(&bBig) == -1 {
		return append(a, b...)
	}
	return append(b, a...)
}

// verifyNode walks down the tree until hitting a leaf, calculating the hash at each level
// and returning the resulting hash of Node n.
func (n *Node) verifyNode(sort bool) ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}
	rightBytes, err := n.Right.verifyNode(sort)
	if err != nil {
		return nil, err
	}

	leftBytes, err := n.Left.verifyNode(sort)
	if err != nil {
		return nil, err
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(sortAppend(sort, leftBytes, rightBytes)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// calculateNodeHash is a helper function that calculates the hash of the node.
func (n *Node) calculateNodeHash(sort bool) ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(sortAppend(sort, n.Left.Hash, n.Right.Hash)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// NewTree creates a new Merkle Tree using the content cs.
func NewTree(cs []Content) (*MerkleTree, error) {
	var defaultHashStrategy = sha256.New
	t := &MerkleTree{
		hashStrategy: defaultHashStrategy,
		sort:         false,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

// NewTreeWithHashStrategy creates a new Merkle Tree using the content cs using the provided hash
// strategy. Note that the hash type used in the type that implements the Content interface must
// match the hash type provided to the tree.
func NewTreeWithHashStrategy(cs []Content, hashStrategy func() hash.Hash) (*MerkleTree, error) {
	t := &MerkleTree{
		hashStrategy: hashStrategy,
		sort:         false,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

// NewTreeWithHashStrategySorted just like NewTreeWithHashStrategy
// but sorts the siblings before hashing, mostly to follow the OpenZepplin Merkle implementation
// https://github.com/OpenZeppelin/openzeppelin-contracts-ethereum-package/blob/master/contracts/cryptography/MerkleProof.sol
func NewTreeWithHashStrategySorted(cs []Content, hashStrategy func() hash.Hash, sort bool) (*MerkleTree, error) {
	t := &MerkleTree{
		hashStrategy: hashStrategy,
		sort:         sort,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

// GetMerklePath: Get Merkle path and indexes(left leaf or right leaf)
func (m *MerkleTree) GetMerklePath(content Content) ([][]byte, []int64, error) {
	for _, current := range m.Leafs {
		ok, err := current.C.Equals(content)
		if err != nil {
			return nil, nil, err
		}

		if ok {
			currentParent := current.Parent
			var merklePath [][]byte
			var index []int64
			for currentParent != nil {
				if bytes.Equal(currentParent.Left.Hash, current.Hash) {
					merklePath = append(merklePath, currentParent.Right.Hash)
					index = append(index, 1) // right leaf
				} else {
					merklePath = append(merklePath, currentParent.Left.Hash)
					index = append(index, 0) // left leaf
				}
				current = currentParent
				currentParent = currentParent.Parent
			}
			return merklePath, index, nil
		}
	}
	return nil, nil, nil
}

// buildWithContent is a helper function that for a given set of Contents, generates a
// corresponding tree and returns the root node, a list of leaf nodes, and a possible error.
// Returns an error if cs contains no Contents.
func buildWithContent(cs []Content, t *MerkleTree) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: cannot construct tree with no content")
	}
	var leafs []*Node
	for _, c := range cs {
		hash, err := c.CalculateHash()
		if err != nil {
			return nil, nil, err
		}

		leafs = append(leafs, &Node{
			Hash: hash,
			C:    c,
			leaf: true,
			Tree: t,
		})
	}
	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			C:    leafs[len(leafs)-1].C,
			leaf: true,
			dup:  true,
			Tree: t,
		}
		leafs = append(leafs, duplicate)
	}
	root, err := buildIntermediate(leafs, t)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}

// buildIntermediate is a helper function that for a given list of leaf nodes, constructs
// the intermediate and root levels of the tree. Returns the resulting root node of the tree.
func buildIntermediate(nl []*Node, t *MerkleTree) (*Node, error) {
	var nodes []*Node
	for i := 0; i < len(nl); i += 2 {
		h := t.hashStrategy()
		var left, right int = i, i + 1
		if i+1 == len(nl) {
			right = i
		}
		chash := sortAppend(t.sort, nl[left].Hash, nl[right].Hash)
		if _, err := h.Write(chash); err != nil {
			return nil, err
		}
		n := &Node{
			Left:  nl[left],
			Right: nl[right],
			Hash:  h.Sum(nil),
			Tree:  t,
		}
		nodes = append(nodes, n)
		nl[left].Parent = n
		nl[right].Parent = n
		if len(nl) == 2 {
			return n, nil
		}
	}
	return buildIntermediate(nodes, t)
}

// MerkleRoot returns the unverified Merkle Root (hash of the root node) of the tree.
func (m *MerkleTree) MerkleRoot() []byte {
	return m.merkleRoot
}

// RebuildTree is a helper function that will rebuild the tree reusing only the content that
// it holds in the leaves.
func (m *MerkleTree) RebuildTree() error {
	var cs []Content
	for _, c := range m.Leafs {
		cs = append(cs, c.C)
	}
	root, leafs, err := buildWithContent(cs, m)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

// RebuildTreeWith replaces the content of the tree and does a complete rebuild; while the root of
// the tree will be replaced the MerkleTree completely survives this operation. Returns an error if the
// list of content cs contains no entries.
func (m *MerkleTree) RebuildTreeWith(cs []Content) error {
	root, leafs, err := buildWithContent(cs, m)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

// VerifyTree() 递归计算得到的根哈希是否和 merkleRoot 相等.
func (m *MerkleTree) VerifyTree() (bool, error) {
	calculatedMerkleRoot, err := m.Root.verifyNode(m.sort)
	if err != nil {
		return false, err
	}

	if bytes.Compare(m.merkleRoot, calculatedMerkleRoot) == 0 {
		return true, nil
	}
	return false, nil
}

// VerifyContent(content) 找到对应叶子，一路往上校验哈希.
func (m *MerkleTree) VerifyContent(content Content) (bool, error) {
	for _, l := range m.Leafs {
		ok, err := l.C.Equals(content)
		if err != nil {
			return false, err
		}

		if ok {
			currentParent := l.Parent
			for currentParent != nil {
				h := m.hashStrategy()
				rightBytes, err := currentParent.Right.calculateNodeHash(m.sort)
				if err != nil {
					return false, err
				}

				leftBytes, err := currentParent.Left.calculateNodeHash(m.sort)
				if err != nil {
					return false, err
				}

				if _, err := h.Write(sortAppend(m.sort, leftBytes, rightBytes)); err != nil {
					return false, err
				}
				if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
					return false, nil
				}
				currentParent = currentParent.Parent
			}
			return true, nil
		}
	}
	return false, nil
}

// String returns a string representation of the node.
func (n *Node) String() string {
	return fmt.Sprintf("%t %t %v %s", n.leaf, n.dup, n.Hash, n.C)
}

// String returns a string representation of the tree. Only leaf nodes are included
// in the output.
func (m *MerkleTree) String() string {
	s := ""
	for _, l := range m.Leafs {
		s += fmt.Sprint(l)
		s += "\n"
	}
	return s
}
