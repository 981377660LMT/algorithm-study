// # from titan_pylib.data_structures.splay_tree.splay_tree_list_array import SplayTreeListArray
// from array import array
// from typing import Generic, List, TypeVar, Tuple, Iterable, Union
// from __pypy__ import newlist_hint
// T = TypeVar('T')

// class SplayTreeListArrayData(Generic[T]):

//   def __init__(self, e: T=0):
//     self.keys: List[T] = [e]
//     self.e: T = e
//     self.arr: array[int] = array('I', bytes(16))
//     # left:  arr[node<<2]
//     # right: arr[node<<2|1]
//     # size:  arr[node<<2|2]
//     # rev:   arr[node<<2|3]
//     self.end: int = 1

//   def reserve(self, n: int) -> None:
//     if n <= 0:
//       return
//     self.keys += [self.e] * (2 * n)
//     self.arr += array('I', bytes(16 * n))

//   def _update_triple(self, x: int, y: int, z: int) -> None:
//     arr = self.arr
//     arr[z<<2|2] = arr[x<<2|2]
//     arr[x<<2|2] = 1 + arr[arr[x<<2]<<2|2] + arr[arr[x<<2|1]<<2|2]
//     arr[y<<2|2] = 1 + arr[arr[y<<2]<<2|2] + arr[arr[y<<2|1]<<2|2]

//   def _update_double(self, x: int, y: int) -> None:
//     arr = self.arr
//     arr[y<<2|2] = arr[x<<2|2]
//     arr[x<<2|2] = 1 + arr[arr[x<<2]<<2|2] + arr[arr[x<<2|1]<<2|2]

//   def _update(self, node: int) -> None:
//     arr = self.arr
//     arr[node<<2|2] = 1 + arr[arr[node<<2]<<2|2] + arr[arr[node<<2|1]<<2|2]

//   def _make_node(self, key: T) -> int:
//     if self.end >= len(self.arr)//4:
//       self.keys.append(key)
//       self.arr.append(0)
//       self.arr.append(0)
//       self.arr.append(1)
//       self.arr.append(0)
//     else:
//       self.keys[self.end] = key
//     self.end += 1
//     return self.end - 1

//   def _splay(self, path: List[int], d: int) -> None:
//     arr = self.arr
//     g = d & 1
//     while len(path) > 1:
//       pnode = path.pop()
//       gnode = path.pop()
//       f = d >> 1 & 1
//       node = arr[pnode<<2|g^1]
//       nnode = (pnode if g == f else node) << 2 | f
//       arr[pnode<<2|g^1] = arr[node<<2|g]
//       arr[node<<2|g] = pnode
//       arr[gnode<<2|f^1] = arr[nnode]
//       arr[nnode] = gnode
//       self._update_triple(gnode, pnode, node)
//       if not path:
//         return
//       d >>= 2
//       g = d & 1
//       arr[path[-1]<<2|g^1] = node
//     pnode = path.pop()
//     node = arr[pnode<<2|g^1]
//     arr[pnode<<2|g^1] = arr[node<<2|g]
//     arr[node<<2|g] = pnode
//     self._update_double(pnode, node)

// class SplayTreeListArray(Generic[T]):

//   def __init__(self,
//                data: SplayTreeListArrayData,
//                n_or_a: Union[int, Iterable[T]]=0,
//                _root: int=0
//                ) -> None:
//     self.data: SplayTreeListArrayData = data
//     self.root: int = _root
//     if not n_or_a:
//       return
//     if isinstance(n_or_a, int):
//       a = [data.e for _ in range(n_or_a)]
//     else:
//       a = list(n_or_a)
//     if a:
//       self._build(a)

//   def _build(self, a: List[T]) -> None:
//     def rec(l: int, r: int) -> int:
//       mid = (l + r) >> 1
//       if l != mid:
//         arr[mid<<2] = rec(l, mid)
//       if mid + 1 != r:
//         arr[mid<<2|1] = rec(mid+1, r)
//       self.data._update(mid)
//       return mid
//     n = len(a)
//     keys, arr = self.data.keys, self.data.arr
//     end = self.data.end
//     self.data.reserve(n+end-len(keys)//2+1)
//     self.data.end += n
//     keys[end:end+n] = a
//     self.root = rec(end, n+end)

//   def _kth_elm_splay(self, node: int, k: int) -> int:
//     arr = self.data.arr
//     if k < 0: k += arr[node<<2|2]
//     d = 0
//     path = []
//     while True:
//       t = arr[arr[node<<2]<<2|2]
//       if t == k:
//         if path:
//           self.data._splay(path, d)
//         return node
//       d = d << 1 | (t > k)
//       path.append(node)
//       node = arr[node<<2|(t<k)]
//       if t < k:
//         k -= t + 1

//   def _left_splay(self, node: int) -> int:
//     if not node: return 0
//     arr = self.data.arr
//     if not arr[node<<2]: return node
//     path = []
//     while arr[node<<2]:
//       path.append(node)
//       node = arr[node<<2]
//     self.data._splay(path, (1<<len(path))-1)
//     return node

//   def _right_splay(self, node: int) -> int:
//     if not node: return 0
//     arr = self.data.arr
//     if not arr[node<<2|1]: return node
//     path = []
//     while arr[node<<2|1]:
//       path.append(node)
//       node = arr[node<<2|1]
//     self.data._splay(path, 0)
//     return node

//   def reserve(self, n: int) -> None:
//     self.data.reserve(n)

//   def merge(self, other: 'SplayTreeListArray') -> None:
//     assert self.data is other.data
//     if not other.root: return
//     if not self.root:
//       self.root = other.root
//       return
//     self.root = self._right_splay(self.root)
//     self.data.arr[self.root<<2|1] = other.root
//     self.data._update(self.root)

//   def split(self, k: int) -> Tuple['SplayTreeListArray', 'SplayTreeListArray']:
//     assert -len(self) < k <= len(self), \
//         f'IndexError: SplayTreeListArray.split({k}), len={len(self)}'
//     if k < 0: k += len(self)
//     if k >= self.data.arr[self.root<<2|2]:
//       return self, SplayTreeListArray(self.data, _root=0)
//     self.root = self._kth_elm_splay(self.root, k)
//     left = SplayTreeListArray(self.data, _root=self.data.arr[self.root<<2])
//     self.data.arr[self.root<<2] = 0
//     self.data._update(self.root)
//     return left, self

//   def _internal_split(self, k: int) -> Tuple[int, int]:
//     if k >= self.data.arr[self.root<<2|2]:
//       return self.root, 0
//     self.root = self._kth_elm_splay(self.root, k)
//     left = self.data.arr[self.root<<2]
//     self.data.arr[self.root<<2] = 0
//     self.data._update(self.root)
//     return left, self.root

//   def insert(self, k: int, key: T) -> None:
//     assert -len(self) <= k <= len(self), \
//         f'IndexError: SplayTreeListArray.insert({k}, {key}), len={len(self)}'
//     if k < 0: k += len(self)
//     data = self.data
//     node = self.data._make_node(key)
//     if not self.root:
//       self.data._update(node)
//       self.root = node
//       return
//     arr = data.arr
//     if k == data.arr[self.root<<2|2]:
//       arr[node<<2] = self._right_splay(self.root)
//     else:
//       node_ = self._kth_elm_splay(self.root, k)
//       if arr[node_<<2]:
//         arr[node<<2] = arr[node_<<2]
//         arr[node_<<2] = 0
//         self.data._update(node_)
//       arr[node<<2|1] = node_
//     self.data._update(node)
//     self.root = node

//   def append(self, key: T) -> None:
//     data = self.data
//     node = self._right_splay(self.root)
//     self.root = self.data._make_node(key)
//     data.arr[self.root<<2] = node
//     self.data._update(self.root)

//   def appendleft(self, key: T) -> None:
//     node = self._left_splay(self.root)
//     self.root = self.data._make_node(key)
//     self.data.arr[self.root<<2|1] = node
//     self.data._update(self.root)

//   def pop(self, k: int=-1) -> T:
//     assert -len(self) <= k < len(self), \
//         f'IndexError: SplayTreeListArray.pop({k})'
//     data = self.data
//     if k == -1:
//       node = self._right_splay(self.root)
//       self.root = data.arr[node<<2]
//       return data.keys[node]
//     self.root = self._kth_elm_splay(self.root, k)
//     res = data.keys[self.root]
//     if not data.arr[self.root<<2]:
//       self.root = data.arr[self.root<<2|1]
//     elif not data.arr[self.root<<2|1]:
//       self.root = data.arr[self.root<<2]
//     else:
//       node = self._right_splay(data.arr[self.root<<2])
//       data.arr[node<<2|1] = data.arr[self.root<<2|1]
//       self.root = node
//       self.data._update(self.root)
//     return res

//   def popleft(self) -> T:
//     assert self, 'IndexError: SplayTreeListArray.popleft()'
//     node = self._left_splay(self.root)
//     self.root = self.data.arr[node<<2|1]
//     return self.data.keys[node]

//   def rotate(self, x: int) -> None:
//     # 「末尾をを削除し先頭に挿入」をx回
//     n = self.data.arr[self.root<<2|2]
//     l, self = self.split(n-(x%n))
//     self.merge(l)

//   def tolist(self) -> List[T]:
//     node = self.root
//     arr, keys = self.data.arr, self.data.keys
//     stack = newlist_hint(len(self))
//     res = newlist_hint(len(self))
//     while stack or node:
//       if node:
//         stack.append(node)
//         node = arr[node<<2]
//       else:
//         node = stack.pop()
//         res.append(keys[node])
//         node = arr[node<<2|1]
//     return res

//   def clear(self) -> None:
//     self.root = 0

//   def __setitem__(self, k: int, key: T):
//     assert -len(self) <= k < len(self), f'IndexError: SplayTreeListArray.__setitem__({k})'
//     self.root = self._kth_elm_splay(self.root, k)
//     self.data.keys[self.root] = key
//     self.data._update(self.root)

//   def __getitem__(self, k: int) -> T:
//     assert -len(self) <= k < len(self), f'IndexError: SplayTreeListArray.__getitem__({k})'
//     self.root = self._kth_elm_splay(self.root, k)
//     return self.data.keys[self.root]

//   def __iter__(self):
//     self.__iter = 0
//     return self

//   def __next__(self):
//     if self.__iter == self.data.arr[self.root<<2|2]:
//       raise StopIteration
//     res = self.__getitem__(self.__iter)
//     self.__iter += 1
//     return res

//   def __reversed__(self):
//     for i in range(len(self)):
//       yield self.__getitem__(-i-1)

//   def __len__(self):
//     return self.data.arr[self.root<<2|2]

//   def __str__(self):
//     return str(self.tolist())

//   def __bool__(self):
//     return self.root != 0

//   def __repr__(self):
//     return f'SplayTreeListArray({self})'

// api:
//  Set(k, key)
//  Get(k)
//  Insert(k, key)
//  Append(key)
//  AppendLeft(key)
//  Pop(k)
//  RotateRight(x)
//  RotateLeft(x)
//  GetAll()
//  Clear()
//  Len()

package main

func main() {
	arr := []int{1, 2, 3, 4, 5}
	s := NewSplay(5, func(i int32) E { return E(arr[i]) })
	_ = s
}

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }

type Data struct {
	e    E
	keys []E
	arr  []int32
	end  int32
}

func NewData(e E) *Data {
	// left:  arr[node<<2], right: arr[node<<2|1], size: arr[node<<2|2], rev: arr[node<<2|3]
	return &Data{e: e, keys: []E{e}, arr: make([]int32, 4), end: 1}
}

func (d *Data) Reserve(n int32) {
	if n <= 0 {
		return
	}
	d.keys = append(d.keys, make([]E, 2*n)...)
	d.arr = append(d.arr, make([]int32, 4*n)...)
}

func (d *Data) _updateTriple(x, y, z int32) {
	arr := d.arr
	arr[z<<2|2] = arr[x<<2|2]
	arr[x<<2|2] = 1 + arr[arr[x<<2]<<2|2] + arr[arr[x<<2|1]<<2|2]
	arr[y<<2|2] = 1 + arr[arr[y<<2]<<2|2] + arr[arr[y<<2|1]<<2|2]
}

func (d *Data) _updateDouble(x, y int32) {
	arr := d.arr
	arr[y<<2|2] = arr[x<<2|2]
	arr[x<<2|2] = 1 + arr[arr[x<<2]<<2|2] + arr[arr[x<<2|1]<<2|2]
}

func (d *Data) _update(node int32) {
	arr := d.arr
	arr[node<<2|2] = 1 + arr[arr[node<<2]<<2|2] + arr[arr[node<<2|1]<<2|2]
}

func (d *Data) _makeNode(key E) int32 {
	if d.end >= int32(len(d.arr))>>2 {
		d.keys = append(d.keys, key)
		d.arr = append(d.arr, 0, 0, 1, 0)
	} else {
		d.keys[d.end] = key
	}
	d.end++
	return d.end - 1
}

// 这里的path可以直接push/pop.
func (d *Data) _splay(path []int32, dep int32) {
	arr := d.arr
	g := dep & 1
	for len(path) > 1 {
		pNode := path[len(path)-1]
		path = path[:len(path)-1]
		gNode := path[len(path)-1]
		path = path[:len(path)-1]
		f := dep >> 1 & 1
		node := arr[pNode<<2|g^1]
		var nnNode int32
		if g == f {
			nnNode = (pNode << 2) | f
		} else {
			nnNode = (node << 2) | f
		}
		arr[pNode<<2|g^1] = arr[node<<2|g]
		arr[node<<2|g] = pNode
		arr[gNode<<2|f^1] = arr[nnNode]
		arr[nnNode] = gNode
		d._updateTriple(gNode, pNode, node)
		if len(path) == 0 {
			return
		}
		dep >>= 2
		g = dep & 1
		arr[path[len(path)-1]<<2|g^1] = node
	}
	pNode := path[len(path)-1]
	path = path[:len(path)-1]
	node := arr[pNode<<2|g^1]
	arr[pNode<<2|g^1] = arr[node<<2|g]
	arr[node<<2|g] = pNode
	d._updateDouble(pNode, node)
}

type Splay struct {
	root    int32
	data    *Data
	tmpPath []int32
}

func NewSplay(n int32, f func(i int32) E) *Splay {
	res := &Splay{data: NewData(e()), tmpPath: make([]int32, 0, 32)}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func _new(data *Data, n int32, f func(i int32) E, root int32) *Splay {
	res := &Splay{root: root, data: data, tmpPath: make([]int32, 0, 32)}
	if n == 0 {
		return res
	}
	res._build(n, f)
	return res
}

func (s *Splay) Merge(other *Splay) {
	if other.root == 0 {
		return
	}
	if s.root == 0 {
		s.root = other.root
		return
	}
	s.root = s._rightSplay(s.root)
	s.data.arr[s.root<<2|1] = other.root
	s.data._update(s.root)
}

func (s *Splay) Split(k int32) (*Splay, *Splay) {
	if k < 0 {
		k += s.Len()
	}
	if k >= s.data.arr[s.root<<2|2] {
		return s, _new(s.data, 0, nil, 0)
	}
	s.root = s._kthElmSplay(s.root, k)
	left := _new(s.data, 0, nil, s.data.arr[s.root<<2])
	s.data.arr[s.root<<2] = 0
	s.data._update(s.root)
	return left, s
}

func (s *Splay) Insert(k int32, key E) {
	if k < 0 {
		k += s.Len()
	}
	data := s.data
	node := data._makeNode(key)
	if s.root == 0 {
		data._update(node)
		s.root = node
		return
	}
	arr := data.arr
	if k == data.arr[s.root<<2|2] {
		arr[node<<2] = s._rightSplay(s.root)
	} else {
		node_ := s._kthElmSplay(s.root, k)
		if arr[node_<<2] != 0 {
			arr[node<<2] = arr[node_<<2]
			arr[node_<<2] = 0
			data._update(node_)
		}
		arr[node<<2|1] = node_
	}
	data._update(node)
	s.root = node
}

func (s *Splay) Append(key E) {
	node := s._rightSplay(s.root)
	s.root = s.data._makeNode(key)
	s.data.arr[s.root<<2] = node
	s.data._update(s.root)
}

func (s *Splay) AppendLeft(key E) {
	node := s._leftSplay(s.root)
	s.root = s.data._makeNode(key)
	s.data.arr[s.root<<2|1] = node
	s.data._update(s.root)
}

func (s *Splay) Pop(k int32) E {
	if k == -1 {
		node := s._rightSplay(s.root)
		s.root = s.data.arr[node<<2]
		return s.data.keys[node]
	}
	s.root = s._kthElmSplay(s.root, k)
	res := s.data.keys[s.root]
	if s.data.arr[s.root<<2] == 0 {
		s.root = s.data.arr[s.root<<2|1]
	} else if s.data.arr[s.root<<2|1] == 0 {
		s.root = s.data.arr[s.root<<2]
	} else {
		mode := s._rightSplay(s.data.arr[s.root<<2])
		s.data.arr[mode<<2|1] = s.data.arr[s.root<<2|1]
		s.root = mode
		s.data._update(s.root)
	}
	return res
}

func (s *Splay) PopLeft() E {
	node := s._leftSplay(s.root)
	s.root = s.data.arr[node<<2|1]
	return s.data.keys[node]
}

func (s *Splay) RotateRight(x int32) {
	n := s.data.arr[s.root<<2|2]
	a, b := s.Split(n - (x % n))
	b.Merge(a)
}

// func (s *Splay) RotateLeft(x int32) {}

func (s *Splay) GetAll() []E {
	node := s.root
	arr, keys := s.data.arr, s.data.keys
	n := s.Len()
	stack := make([]int32, 0, n)
	res := make([]E, 0, n)
	for len(stack) > 0 || node != 0 {
		if node != 0 {
			stack = append(stack, node)
			node = arr[node<<2]
		} else {
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res = append(res, keys[node])
			node = arr[node<<2|1]
		}
	}
	return res
}

func (s *Splay) Clear() { s.root = 0 }

func (s *Splay) Set(k int32, key E) {
	s.root = s._kthElmSplay(s.root, k)
	s.data.keys[s.root] = key
	s.data._update(s.root)
}

func (s *Splay) Get(k int32) E {
	s.root = s._kthElmSplay(s.root, k)
	return s.data.keys[s.root]
}

func (s *Splay) Len() int32 { return s.data.arr[s.root<<2|2] }

func (s *Splay) _build(n int32, f func(i int32) E) {
	keys, arr := s.data.keys, s.data.arr
	end := s.data.end
	s.data.Reserve(n + end - int32(len(keys))>>1 + 1)
	s.data.end += n
	for i := end; i < end+n; i++ {
		keys[i] = f(i - end)
	}
	var dfs func(l, r int32) int32
	dfs = func(l, r int32) int32 {
		mid := (l + r) >> 1
		if l != mid {
			arr[mid<<2] = dfs(l, mid)
		}
		if mid+1 != r {
			arr[mid<<2|1] = dfs(mid+1, r)
		}
		s.data._update(mid)
		return mid
	}
	s.root = dfs(end, end+n)
}

func (s *Splay) _kthElmSplay(node int32, k int32) int32 {
	arr := s.data.arr
	if k < 0 {
		k += arr[node<<2|2]
	}
	d := int32(0)
	s.tmpPath = s.tmpPath[:0]
	path := s.tmpPath
	for {
		t := arr[arr[node<<2]<<2|2]
		if t == k {
			if len(path) > 0 {
				s.data._splay(path, d)
			}
			return node
		}
		path = append(path, node)
		d <<= 1
		if t > k {
			d |= 1
		}
		if t < k {
			node = arr[node<<2|1]
			k -= t + 1
		} else {
			node = arr[node<<2]
		}
	}
}

func (s *Splay) _leftSplay(node int32) int32 {
	if node == 0 {
		return 0
	}
	arr := s.data.arr
	if arr[node<<2] == 0 {
		return node
	}
	s.tmpPath = s.tmpPath[:0]
	path := s.tmpPath
	for arr[node<<2] != 0 {
		path = append(path, node)
		node = arr[node<<2]
	}
	s.data._splay(path, (1<<len(path))-1) // TODO: 这里是否int32可以
	return node
}

func (s *Splay) _rightSplay(node int32) int32 {
	if node == 0 {
		return 0
	}
	arr := s.data.arr
	if arr[node<<2|1] == 0 {
		return node
	}
	s.tmpPath = s.tmpPath[:0]
	path := s.tmpPath
	for arr[node<<2|1] != 0 {
		path = append(path, node)
		node = arr[node<<2|1]
	}
	s.data._splay(path, 0)
	return node
}

func (s *Splay) _reserve(n int32) {
	s.data.Reserve(n)
}

func (s *Splay) _internalSplit(k int32) (int32, int32) {
	if k >= s.data.arr[s.root<<2|2] {
		return s.root, 0
	}
	s.root = s._kthElmSplay(s.root, k)
	left := s.data.arr[s.root<<2]
	s.data.arr[s.root<<2] = 0
	s.data._update(s.root)
	return left, s.root
}
