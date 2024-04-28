# from titan_pylib.data_structures.splay_tree.splay_tree_list_array import SplayTreeListArray
from array import array
from typing import Generic, List, TypeVar, Tuple, Iterable, Union
from __pypy__ import newlist_hint
T = TypeVar('T')

class SplayTreeListArrayData(Generic[T]):

  def __init__(self, e: T=0):
    self.keys: List[T] = [e]
    self.e: T = e
    self.arr: array[int] = array('I', bytes(16))
    # left:  arr[node<<2]
    # right: arr[node<<2|1]
    # size:  arr[node<<2|2]
    # rev:   arr[node<<2|3]
    self.end: int = 1

  def reserve(self, n: int) -> None:
    if n <= 0:
      return
    self.keys += [self.e] * (2 * n)
    self.arr += array('I', bytes(16 * n))

  def _update_triple(self, x: int, y: int, z: int) -> None:
    arr = self.arr
    arr[z<<2|2] = arr[x<<2|2]
    arr[x<<2|2] = 1 + arr[arr[x<<2]<<2|2] + arr[arr[x<<2|1]<<2|2]
    arr[y<<2|2] = 1 + arr[arr[y<<2]<<2|2] + arr[arr[y<<2|1]<<2|2]

  def _update_double(self, x: int, y: int) -> None:
    arr = self.arr
    arr[y<<2|2] = arr[x<<2|2]
    arr[x<<2|2] = 1 + arr[arr[x<<2]<<2|2] + arr[arr[x<<2|1]<<2|2]

  def _update(self, node: int) -> None:
    arr = self.arr
    arr[node<<2|2] = 1 + arr[arr[node<<2]<<2|2] + arr[arr[node<<2|1]<<2|2]

  def _make_node(self, key: T) -> int:
    if self.end >= len(self.arr)//4:
      self.keys.append(key)
      self.arr.append(0)
      self.arr.append(0)
      self.arr.append(1)
      self.arr.append(0)
    else:
      self.keys[self.end] = key
    self.end += 1
    return self.end - 1

  def _splay(self, path: List[int], d: int) -> None:
    arr = self.arr
    g = d & 1
    while len(path) > 1:
      pnode = path.pop()
      gnode = path.pop()
      f = d >> 1 & 1
      node = arr[pnode<<2|g^1]
      nnode = (pnode if g == f else node) << 2 | f
      arr[pnode<<2|g^1] = arr[node<<2|g]
      arr[node<<2|g] = pnode
      arr[gnode<<2|f^1] = arr[nnode]
      arr[nnode] = gnode
      self._update_triple(gnode, pnode, node)
      if not path:
        return
      d >>= 2
      g = d & 1
      arr[path[-1]<<2|g^1] = node
    pnode = path.pop()
    node = arr[pnode<<2|g^1]
    arr[pnode<<2|g^1] = arr[node<<2|g]
    arr[node<<2|g] = pnode
    self._update_double(pnode, node)

class SplayTreeListArray(Generic[T]):

  def __init__(self,
               data: SplayTreeListArrayData,
               n_or_a: Union[int, Iterable[T]]=0,
               _root: int=0
               ) -> None:
    self.data: SplayTreeListArrayData = data
    self.root: int = _root
    if not n_or_a:
      return
    if isinstance(n_or_a, int):
      a = [data.e for _ in range(n_or_a)]
    else:
      a = list(n_or_a)
    if a:
      self._build(a)

  def _build(self, a: List[T]) -> None:
    def rec(l: int, r: int) -> int:
      mid = (l + r) >> 1
      if l != mid:
        arr[mid<<2] = rec(l, mid)
      if mid + 1 != r:
        arr[mid<<2|1] = rec(mid+1, r)
      self.data._update(mid)
      return mid
    n = len(a)
    keys, arr = self.data.keys, self.data.arr
    end = self.data.end
    self.data.reserve(n+end-len(keys)//2+1)
    self.data.end += n
    keys[end:end+n] = a
    self.root = rec(end, n+end)

  def _kth_elm_splay(self, node: int, k: int) -> int:
    arr = self.data.arr
    if k < 0: k += arr[node<<2|2]
    d = 0
    path = []
    while True:
      t = arr[arr[node<<2]<<2|2]
      if t == k:
        if path:
          self.data._splay(path, d)
        return node
      d = d << 1 | (t > k)
      path.append(node)
      node = arr[node<<2|(t<k)]
      if t < k:
        k -= t + 1

  def _left_splay(self, node: int) -> int:
    if not node: return 0
    arr = self.data.arr
    if not arr[node<<2]: return node
    path = []
    while arr[node<<2]:
      path.append(node)
      node = arr[node<<2]
    self.data._splay(path, (1<<len(path))-1)
    return node

  def _right_splay(self, node: int) -> int:
    if not node: return 0
    arr = self.data.arr
    if not arr[node<<2|1]: return node
    path = []
    while arr[node<<2|1]:
      path.append(node)
      node = arr[node<<2|1]
    self.data._splay(path, 0)
    return node

  def reserve(self, n: int) -> None:
    self.data.reserve(n)

  def merge(self, other: 'SplayTreeListArray') -> None:
    assert self.data is other.data
    if not other.root: return
    if not self.root:
      self.root = other.root
      return
    self.root = self._right_splay(self.root)
    self.data.arr[self.root<<2|1] = other.root
    self.data._update(self.root)

  def split(self, k: int) -> Tuple['SplayTreeListArray', 'SplayTreeListArray']:
    assert -len(self) < k <= len(self), \
        f'IndexError: SplayTreeListArray.split({k}), len={len(self)}'
    if k < 0: k += len(self)
    if k >= self.data.arr[self.root<<2|2]:
      return self, SplayTreeListArray(self.data, _root=0)
    self.root = self._kth_elm_splay(self.root, k)
    left = SplayTreeListArray(self.data, _root=self.data.arr[self.root<<2])
    self.data.arr[self.root<<2] = 0
    self.data._update(self.root)
    return left, self

  def _internal_split(self, k: int) -> Tuple[int, int]:
    if k >= self.data.arr[self.root<<2|2]:
      return self.root, 0
    self.root = self._kth_elm_splay(self.root, k)
    left = self.data.arr[self.root<<2]
    self.data.arr[self.root<<2] = 0
    self.data._update(self.root)
    return left, self.root

  def insert(self, k: int, key: T) -> None:
    assert -len(self) <= k <= len(self), \
        f'IndexError: SplayTreeListArray.insert({k}, {key}), len={len(self)}'
    if k < 0: k += len(self)
    data = self.data
    node = self.data._make_node(key)
    if not self.root:
      self.data._update(node)
      self.root = node
      return
    arr = data.arr
    if k == data.arr[self.root<<2|2]:
      arr[node<<2] = self._right_splay(self.root)
    else:
      node_ = self._kth_elm_splay(self.root, k)
      if arr[node_<<2]:
        arr[node<<2] = arr[node_<<2]
        arr[node_<<2] = 0
        self.data._update(node_)
      arr[node<<2|1] = node_
    self.data._update(node)
    self.root = node

  def append(self, key: T) -> None:
    data = self.data
    node = self._right_splay(self.root)
    self.root = self.data._make_node(key)
    data.arr[self.root<<2] = node
    self.data._update(self.root)

  def appendleft(self, key: T) -> None:
    node = self._left_splay(self.root)
    self.root = self.data._make_node(key)
    self.data.arr[self.root<<2|1] = node
    self.data._update(self.root)

  def pop(self, k: int=-1) -> T:
    assert -len(self) <= k < len(self), \
        f'IndexError: SplayTreeListArray.pop({k})'
    data = self.data
    if k == -1:
      node = self._right_splay(self.root)
      self.root = data.arr[node<<2]
      return data.keys[node]
    self.root = self._kth_elm_splay(self.root, k)
    res = data.keys[self.root]
    if not data.arr[self.root<<2]:
      self.root = data.arr[self.root<<2|1]
    elif not data.arr[self.root<<2|1]:
      self.root = data.arr[self.root<<2]
    else:
      node = self._right_splay(data.arr[self.root<<2])
      data.arr[node<<2|1] = data.arr[self.root<<2|1]
      self.root = node
      self.data._update(self.root)
    return res

  def popleft(self) -> T:
    assert self, 'IndexError: SplayTreeListArray.popleft()'
    node = self._left_splay(self.root)
    self.root = self.data.arr[node<<2|1]
    return self.data.keys[node]

  def rotate(self, x: int) -> None:
    # 「末尾をを削除し先頭に挿入」をx回
    n = self.data.arr[self.root<<2|2]
    l, self = self.split(n-(x%n))
    self.merge(l)

  def tolist(self) -> List[T]:
    node = self.root
    arr, keys = self.data.arr, self.data.keys
    stack = newlist_hint(len(self))
    res = newlist_hint(len(self))
    while stack or node:
      if node:
        stack.append(node)
        node = arr[node<<2]
      else:
        node = stack.pop()
        res.append(keys[node])
        node = arr[node<<2|1]
    return res

  def clear(self) -> None:
    self.root = 0

  def __setitem__(self, k: int, key: T):
    assert -len(self) <= k < len(self), f'IndexError: SplayTreeListArray.__setitem__({k})'
    self.root = self._kth_elm_splay(self.root, k)
    self.data.keys[self.root] = key
    self.data._update(self.root)

  def __getitem__(self, k: int) -> T:
    assert -len(self) <= k < len(self), f'IndexError: SplayTreeListArray.__getitem__({k})'
    self.root = self._kth_elm_splay(self.root, k)
    return self.data.keys[self.root]

  def __iter__(self):
    self.__iter = 0
    return self

  def __next__(self):
    if self.__iter == self.data.arr[self.root<<2|2]:
      raise StopIteration
    res = self.__getitem__(self.__iter)
    self.__iter += 1
    return res

  def __reversed__(self):
    for i in range(len(self)):
      yield self.__getitem__(-i-1)

  def __len__(self):
    return self.data.arr[self.root<<2|2]

  def __str__(self):
    return str(self.tolist())

  def __bool__(self):
    return self.root != 0

  def __repr__(self):
    return f'SplayTreeListArray({self})'

