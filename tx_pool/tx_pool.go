package tx_pool

import (
	"container/heap"
	"sort"
	"go-blockchain/blc"
)

type PriceIndex struct {
	price uint64
	index uint64
}

type priceHeap []PriceIndex
func (h priceHeap) Len() int { return len(h) }
func (h priceHeap) Less(i, j int) bool { return h[i].price < h[j].price }
func (h priceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *priceHeap) Push(x interface{}) {
	*h = append(*h, x.(PriceIndex))
}
func (h *priceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type txSortedMap struct {
	items map[uint64]*blc.Transaction // 交易列表
	index *priceHeap // 价格堆
	cache []*blc.Transaction // 已排序交易
}

// 创建新堆空间
func newTxSortedMap() *txSortedMap {
	return &txSortedMap{
		items: make(map[uint64]*blc.Transaction),
		index: new(priceHeap),
	}
}

// 获取对应索引的交易
func (m *txSortedMap) Get(index uint64) *blc.Transaction {
	return m.items[index]
}

// 添加交易
func (m *txSortedMap) Put(tx *blc.Transaction, index uint64) {
	if m.items[index] == nil {
		pi := PriceIndex{
			// TODO：此处需要Transaction提供price方法获取矿工的手续费
			// price: tx.price(),
			index: index,
		}
		heap.Push(m.index, pi)
	}
	m.items[index], m.cache = tx, nil
}

// 删除低于价格阈值的交易
func (m *txSortedMap) Forward(threshold uint64) []*blc.Transaction {
	var removed []*blc.Transaction

	for m.index.Len() > 0 && (*m.index)[0].price < threshold {
		priceIndex := heap.Pop(m.index).(PriceIndex)
		removed = append(removed, m.items[priceIndex.index])
		delete(m.items, priceIndex.index)
	}
	if m.cache != nil {
		m.cache = m.cache[len(removed):]
	}
	return removed
}

// 根据筛选函数进行筛选，再重建堆
func (m *txSortedMap) Filter(filter func(*blc.Transaction) bool) []*blc.Transaction {
	removed := m.filter(filter)
	if len(removed) > 0 {
		m.reheap()
	}
	return removed
}

// 重建堆
func (m *txSortedMap) reheap() {
	*m.index = make([]PriceIndex, 0, len(m.items))
	// for index, tx := range m.items {
	for index := range m.items {
		*m.index = append(*m.index, PriceIndex{
			// TODO：此处需要Transaction提供price方法获取矿工的手续费
			// price: tx.price(),
			index: index,
		})
	}
	heap.Init(m.index)
	m.cache = nil
}

// 根据筛选函数进行筛选交易
func (m *txSortedMap) filter(filter func(*blc.Transaction) bool) []*blc.Transaction {
	var removed []*blc.Transaction
	for index, tx := range m.items {
		if filter(tx) {
			removed = append(removed, tx)
			delete(m.items, index)
		}
	}
	if len(removed) > 0 {
		m.cache = nil
	}
	return removed
}

// 设置交易池上限，将超过上限的交易删除，并重建堆
func (m *txSortedMap) Cap(threshold int) []*blc.Transaction {
	if len(m.items) <= threshold {
		return nil
	}
	var drops []*blc.Transaction
	sort.Sort(*m.index)
	for size := len(m.items); size > threshold; size-- {
		drops = append(drops, m.items[(*m.index)[size-1].index])
		delete(m.items, (*m.index)[size-1].index)
	}
	*m.index = (*m.index)[:threshold]
	heap.Init(m.index)
	if m.cache != nil {
		m.cache = m.cache[:len(m.cache)-len(drops)]
	}
	return drops
}

// 删除指定索引的交易
func (m *txSortedMap) Remove(index uint64) bool {
	_, ok := m.items[index]
	if !ok {
		return false
	}
	// 从堆中删除交易
	for i := 0; i < m.index.Len(); i++ {
		if (*m.index)[i].index == index {
			heap.Remove(m.index, i)
			break
		}
	}
	delete(m.items, index)
	m.cache = nil

	return true
}

// 获取即将处理的交易
func (m *txSortedMap) Ready(start uint64) []*blc.Transaction {
	if m.index.Len() == 0 || (*m.index)[0].price < start {
		return nil
	}
	var ready []*blc.Transaction
	for next := (*m.index)[0].index; m.index.Len() > 0 && (*m.index)[0].index == next; next++ {
		ready = append(ready, m.items[next])
		delete(m.items, next)
		heap.Pop(m.index)
	}
	m.cache = nil

	return ready
}

// 返回交易列表长度
func (m *txSortedMap) Len() int {
	return len(m.items)
}

// 返回最大价格的交易
func (m *txSortedMap) LastElement() *blc.Transaction {
	return m.cache[len(m.cache)-1]
}
