package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
)

type Content interface {
	CalculateHash() ([]byte, error) // 计算元数据的hash
	Equals(other Content) (bool, error) // 判断两个元数据内容是否相等
}

type MerkleTree struct {
	Root         *Node // 根节点
	merkleRoot   []byte // 根Hash
	Leafs        []*Node // 叶子节点列表
	hashStrategy func() hash.Hash //hash策略
}

type Node struct {
	Tree   *MerkleTree // 所属merkle树
	Parent *Node // 父节点
	Left   *Node // 左节点
	Right  *Node // 右节点
	leaf   bool // 是否是叶子节点
	dup    bool // 是否是副本（叶子节点为单数时需要）
	Hash   []byte // hash值
	C      Content // 元数据操作方法接口
}

// 重新计算整棵树的根hash，用于校验
func (n *Node) verifyNode() ([]byte, error) {
	// 计算叶子节点hash
	if n.leaf {
		// 叶子节点的hash，直接调用交易的hash值
		return n.C.CalculateHash()
	}
	// 中序遍历全树，获取各个子树的hash
	rightBytes, err := n.Right.verifyNode()
	if err != nil {
		return nil, err
	}
	leftBytes, err := n.Left.verifyNode()
	if err != nil {
		return nil, err
	}

	// 获取Merkle树的hash的策略，默认值是sha256
	h := n.Tree.hashStrategy()
	// 求子树的hash
	if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// 计算节点hash
func (n *Node) calculateNodeHash() ([]byte, error) {
	// 计算叶子节点hash
	if n.leaf {
		// 叶子节点的hash，直接调用交易的hash值
		return n.C.CalculateHash()
	}

	// 获取Merkle树的hash的策略，默认值是sha256
	h := n.Tree.hashStrategy()
	// 计算子树的hash
	if _, err := h.Write(append(n.Left.Hash, n.Right.Hash...)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// 创建merkle树
func NewTree(cs []Content) (*MerkleTree, error) {
	// 默认使用的hash算法为sha256
	var defaultHashStrategy = sha256.New
	t := &MerkleTree{
		hashStrategy: defaultHashStrategy,
	}

	// 利用data数据构建merkle树
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash

	// 返回MerkleTree结构
	return t, nil
}

// 使用自定义hash算法创建merkle树
func NewTreeWithHashStrategy(cs []Content, hashStrategy func() hash.Hash) (*MerkleTree, error) {
	t := &MerkleTree{
		hashStrategy: hashStrategy,
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

// 查找某数据的merkle树路径
func (m *MerkleTree) GetMerklePath(content Content) ([][]byte, []int64, error) {
	for _, current := range m.Leafs {
		// 遍历叶子节点查找对应的元数据
		ok, err := current.C.Equals(content)
		if err != nil {
			return nil, nil, err
		}
		// 查到节点后
		if ok {
			currentParent := current.Parent
			var merklePath [][]byte
			var index []int64
			// 向树顶遍历，将路径上的父节点的hash和索引集合返回
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

// 利用数据创建merkle树
func buildWithContent(cs []Content, t *MerkleTree) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: cannot construct tree with no content")
	}

	var leafs []*Node
	// 遍历数据集合
	for _, c := range cs {
		// 计算元数据hash
		hash, err := c.CalculateHash()
		if err != nil {
			return nil, nil, err
		}

		// 添加到叶子节点中
		leafs = append(leafs, &Node{
			Hash: hash,
			C:    c,
			leaf: true,
			Tree: t,
		})
	}

	// 如果叶子节点是单数，就要多复制一个叶子节点
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

	// 递归生成merkle树
	root, err := buildIntermediate(leafs, t)
	if err != nil {
		return nil, nil, err
	}

	// 返回merkle树根节点，叶子节点集合
	return root, leafs, nil
}

// 生成上一层节点集合
func buildIntermediate(nl []*Node, t *MerkleTree) (*Node, error) {
	var nodes []*Node
	// 两个为一组，进行遍历
	for i := 0; i < len(nl); i += 2 {
		h := t.hashStrategy()
		var left, right int = i, i + 1
		if i+1 == len(nl) {
			right = i
		}
		// 计算两个叶子节点的hash，并生成父节点
		chash := append(nl[left].Hash, nl[right].Hash...)
		if _, err := h.Write(chash); err != nil {
			return nil, err
		}
		n := &Node{
			Left:  nl[left],
			Right: nl[right],
			Hash:  h.Sum(nil),
			Tree:  t,
		}
		
		// 生成上一层的节点集合
		nodes = append(nodes, n)
		nl[left].Parent = n
		nl[right].Parent = n

		// 当完成最顶层时直接返回
		if len(nl) == 2 {
			return n, nil
		}
	}

	// 递归向上生成merkle树节点
	return buildIntermediate(nodes, t)
}

// 返回根hash
func (m *MerkleTree) MerkleRoot() []byte {
	return m.merkleRoot
}

// 重建新树
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

// 利用元数据集合重建新树
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

// 校验树的根hash
func (m *MerkleTree) VerifyTree() (bool, error) {
	// 计算根hash
	calculatedMerkleRoot, err := m.Root.verifyNode()
	if err != nil {
		return false, err
	}
	// 比对根hash值
	if bytes.Compare(m.merkleRoot, calculatedMerkleRoot) == 0 {
		return true, nil
	}
	return false, nil
}

// 检查元数据是否合法存在于树中
func (m *MerkleTree) VerifyContent(content Content) (bool, error) {
	for _, l := range m.Leafs {
		// 遍历叶子节点查找对应的元数据
		ok, err := l.C.Equals(content)
		if err != nil {
			return false, err
		}

		// 查找到节点后
		if ok {
			currentParent := l.Parent
			for currentParent != nil {
				h := m.hashStrategy()
				// 计算父节点的左右节点hash
				rightBytes, err := currentParent.Right.calculateNodeHash()
				if err != nil {
					return false, err
				}

				leftBytes, err := currentParent.Left.calculateNodeHash()
				if err != nil {
					return false, err
				}

				if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
					return false, err
				}
				// 比对计算出来的父节点hash是否正确
				if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
					return false, nil
				}
				// 递归向上检查父节点
				currentParent = currentParent.Parent
			}
			return true, nil
		}
	}
	return false, nil
}

// 打印节点信息
func (n *Node) String() string {
	return fmt.Sprintf("%t %t %v %s", n.leaf, n.dup, n.Hash, n.C)
}

// 打印树信息
func (m *MerkleTree) String() string {
	s := ""
	for _, l := range m.Leafs {
		s += fmt.Sprint(l)
		s += "\n"
	}
	return s
}