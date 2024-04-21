# from titan_pylib.data_structures.scapegoat_tree.scapegoat_tree_multiset import ScapegoatTreeMultiset
# from titan_pylib.my_class.ordered_multiset_interface import OrderedMultisetInterface
# from titan_pylib.my_class.supports_less_than import SupportsLessThan
from typing import Protocol

class SupportsLessThan(Protocol):

  def __lt__(self, other) -> bool: ...

from abc import ABC, abstractmethod
from typing import Iterable, Optional, Iterator, TypeVar, Generic, List
T = TypeVar('T', bound=SupportsLessThan)

class OrderedMultisetInterface(ABC, Generic[T]):

  @abstractmethod
  def __init__(self, a: Iterable[T]) -> None:
    raise NotImplementedError

  @abstractmethod
  def add(self, key: T, cnt: int) -> None:
    raise NotImplementedError

  @abstractmethod
  def discard(self, key: T, cnt: int) -> bool:
    raise NotImplementedError

  @abstractmethod
  def discard_all(self, key: T) -> bool:
    raise NotImplementedError

  @abstractmethod
  def count(self, key: T) -> int:
    raise NotImplementedError

  @abstractmethod
  def remove(self, key: T, cnt: int) -> None:
    raise NotImplementedError

  @abstractmethod
  def le(self, key: T) -> Optional[T]:
    raise NotImplementedError

  @abstractmethod
  def lt(self, key: T) -> Optional[T]:
    raise NotImplementedError

  @abstractmethod
  def ge(self, key: T) -> Optional[T]:
    raise NotImplementedError

  @abstractmethod
  def gt(self, key: T) -> Optional[T]:
    raise NotImplementedError

  @abstractmethod
  def get_max(self) -> Optional[T]:
    raise NotImplementedError

  @abstractmethod
  def get_min(self) -> Optional[T]:
    raise NotImplementedError

  @abstractmethod
  def pop_max(self) -> T:
    raise NotImplementedError

  @abstractmethod
  def pop_min(self) -> T:
    raise NotImplementedError

  @abstractmethod
  def clear(self) -> None:
    raise NotImplementedError

  @abstractmethod
  def tolist(self) -> List[T]:
    raise NotImplementedError

  @abstractmethod
  def __iter__(self) -> Iterator:
    raise NotImplementedError

  @abstractmethod
  def __next__(self) -> T:
    raise NotImplementedError

  @abstractmethod
  def __contains__(self, key: T) -> bool:
    raise NotImplementedError

  @abstractmethod
  def __len__(self) -> int:
    raise NotImplementedError

  @abstractmethod
  def __bool__(self) -> bool:
    raise NotImplementedError

  @abstractmethod
  def __str__(self) -> str:
    raise NotImplementedError

  @abstractmethod
  def __repr__(self) -> str:
    raise NotImplementedError

# from titan_pylib.data_structures.bst_base.bst_multiset_node_base import BSTMultisetNodeBase
from __pypy__ import newlist_hint
from typing import List, Tuple, TypeVar, Generic, Optional
T = TypeVar('T')
Node = TypeVar('Node')
# protcolで、key,val,left,right を規定

class BSTMultisetNodeBase(Generic[T, Node]):

  @staticmethod
  def count(node, key: T) -> int:
    while node is not None:
      if node.key == key:
        return node.val
      node = node.left if key < node.key else node.right
    return 0

  @staticmethod
  def get_min(node: Node) -> Optional[T]:
    if not node:
      return None
    while node.left:
      node = node.left
    return node.key

  @staticmethod
  def get_max(node: Node) -> Optional[T]:
    if not node:
      return None
    while node.right:
      node = node.right
    return node.key

  @staticmethod
  def contains(node: Node, key: T) -> bool:
    while node:
      if key == node.key:
        return True
      node = node.left if key < node.key else node.right
    return False

  @staticmethod
  def tolist(node: Node, _len: int=0) -> List[T]:
    stack = []
    a = newlist_hint(_len)
    while stack or node:
      if node:
        stack.append(node)
        node = node.left
      else:
        node = stack.pop()
        for _ in range(node.val):
          a.append(node.key)
        node = node.right
    return a

  @staticmethod
  def tolist_items(node: Node, _len: int=0) -> List[Tuple[T, int]]:
    stack = newlist_hint(_len)
    a = newlist_hint(_len)
    while stack or node:
      if node:
        stack.append(node)
        node = node.left
      else:
        node = stack.pop()
        a.append((node.key, node.val))
        node = node.right
    return a

  @staticmethod
  def _rle(a: List[T]) -> Tuple[List[T], List[int]]:
    keys, vals = newlist_hint(len(a)), newlist_hint(len(a))
    keys.append(a[0])
    vals.append(1)
    for i, elm in enumerate(a):
      if i == 0:
        continue
      if elm == keys[-1]:
        vals[-1] += 1
        continue
      keys.append(elm)
      vals.append(1)
    return keys, vals

  @staticmethod
  def le(node: Node, key: T) -> Optional[T]:
    res = None
    while node is not None:
      if key == node.key:
        res = key
        break
      if key < node.key:
        node = node.left
      else:
        res = node.key
        node = node.right
    return res

  @staticmethod
  def lt(node: Node, key: T) -> Optional[T]:
    res = None
    while node is not None:
      if key <= node.key:
        node = node.left
      else:
        res = node.key
        node = node.right
    return res

  @staticmethod
  def ge(node: Node, key: T) -> Optional[T]:
    res = None
    while node is not None:
      if key == node.key:
        res = key
        break
      if key < node.key:
        res = node.key
        node = node.left
      else:
        node = node.right
    return res

  @staticmethod
  def gt(node: Node, key: T) -> Optional[T]:
    res = None
    while node is not None:
      if key < node.key:
        res = node.key
        node = node.left
      else:
        node = node.right
    return res

  @staticmethod
  def index(node: Node, key: T) -> int:
    k = 0
    while node:
      if key == node.key:
        if node.left:
          k += node.left.valsize
        break
      if key < node.key:
        node = node.left
      else:
        k += node.val if node.left is None else node.left.valsize + node.val
        node = node.right
    return k

  @staticmethod
  def index_right(node: Node, key: T) -> int:
    k = 0
    while node:
      if key == node.key:
        k += node.val if node.left is None else node.left.valsize + node.val
        break
      if key < node.key:
        node = node.left
      else:
        k += node.val if node.left is None else node.left.valsize + node.val
        node = node.right
    return k

import math
from __pypy__ import newlist_hint
from typing import Final, List, TypeVar, Generic, Iterable, Tuple, Optional, Iterator
T = TypeVar('T', bound=SupportsLessThan)

class ScapegoatTreeMultiset(OrderedMultisetInterface, Generic[T]):

  ALPHA: Final[float] = 0.75
  BETA: Final[float] = math.log2(1 / ALPHA)

  class Node():

    def __init__(self, key: T, val: int):
      self.key: T = key
      self.val: int = val
      self.size: int = 1
      self.valsize: int = val
      self.left: Optional[ScapegoatTreeMultiset.Node] = None
      self.right: Optional[ScapegoatTreeMultiset.Node] = None

    def __str__(self):
      if self.left is None and self.right is None:
        return f'key:{self.key, self.val, self.size, self.valsize}\n'
      return f'key:{self.key, self.val, self.size, self.valsize},\n left:{self.left},\n right:{self.right}\n'

  def __init__(self, a: Iterable[T]=[]):
    self.root = None
    if not isinstance(a, list):
      a = list(a)
    self._build(a)

  def _build(self, a: List[T]) -> None:
    Node = ScapegoatTreeMultiset.Node
    def rec(l: int, r: int) -> 'Node':
      mid = (l + r) >> 1
      node = Node(x[mid], y[mid])
      if l != mid:
        node.left = rec(l, mid)
        node.size += node.left.size
        node.valsize += node.left.valsize
      if mid+1 != r:
        node.right = rec(mid+1, r)
        node.size += node.right.size
        node.valsize += node.right.valsize
      return node
    if not all(a[i] <= a[i+1] for i in range(len(a)-1)):
      a = sorted(a)
    if not a:
      return
    x, y = BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node]._rle(a)
    self.root = rec(0, len(x))

  def _rebuild(self, node: Node) -> Node:
    def rec(l: int, r: int) -> ScapegoatTreeMultiset.Node:
      mid = (l + r) >> 1
      node = a[mid]
      node.size = 1
      node.valsize = node.val
      if l != mid:
        node.left = rec(l, mid)
        node.size += node.left.size
        node.valsize += node.left.valsize
      else:
        node.left = None
      if mid+1 != r:
        node.right = rec(mid+1, r)
        node.size += node.right.size
        node.valsize += node.right.valsize
      else:
        node.right = None
      return node
    a = newlist_hint(node.size)
    stack = []
    while stack or node:
      if node:
        stack.append(node)
        node = node.left
      else:
        node = stack.pop()
        a.append(node)
        node = node.right
    return rec(0, len(a))

  def _kth_elm(self, k: int) -> Tuple[T, int]:
    if k < 0:
      k += len(self)
    node = self.root
    while node:
      t = (node.val + node.left.valsize) if node.left else node.val
      if t-node.val <= k and k < t:
        return node.key, node.val
      elif t > k:
        node = node.left
      else:
        node = node.right
        k -= t

  def _kth_elm_tree(self, k: int) -> Tuple[T, int]:
    if k < 0:
      k += self.len_elm()
    node = self.root
    while node:
      t = node.left.size if node.left else 0
      if t == k:
        return node.key, node.val
      if t > k:
        node = node.left
      else:
        node = node.right
        k -= t + 1
    assert False, 'IndexError'

  def add(self, key: T, val: int=1) -> None:
    if val <= 0: return
    if not self.root:
      self.root = ScapegoatTreeMultiset.Node(key, val)
      return
    node = self.root
    path = []
    while node:
      path.append(node)
      if key == node.key:
        node.val += val
        for p in path:
          p.valsize += val
        return
      node = node.left if key < node.key else node.right
    if key < path[-1].key:
      path[-1].left = ScapegoatTreeMultiset.Node(key, val)
    else:
      path[-1].right = ScapegoatTreeMultiset.Node(key, val)
    if len(path)*ScapegoatTreeMultiset.BETA > math.log(self.len_elm()):
      node_size = 1
      while path:
        pnode = path.pop()
        pnode_size = pnode.size + 1
        if ScapegoatTreeMultiset.ALPHA * pnode_size < node_size:
          break
        node_size = pnode_size
      new_node = self._rebuild(pnode)
      if not path:
        self.root = new_node
        return
      if new_node.key < path[-1].key:
        path[-1].left = new_node
      else:
        path[-1].right = new_node
    for p in path:
      p.size += 1
      p.valsize += val

  def _discard(self, key: T) -> bool:
    path = []
    node = self.root
    di, cnt = 1, 0
    while node:
      if key == node.key:
        break
      path.append(node)
      di = key < node.key
      node = node.left if di else node.right
    if node.left and node.right:
      path.append(node)
      lmax = node.left
      di = 0 if lmax.right else 1
      while lmax.right:
        cnt += 1
        path.append(lmax)
        lmax = lmax.right
      lmax_val = lmax.val
      node.key = lmax.key
      node.val = lmax_val
      node = lmax
    cnode = node.left if node.left else node.right
    if path:
      if di == 1:
        path[-1].left = cnode
      else:
        path[-1].right = cnode
    else:
      self.root = cnode
      return True
    for _ in range(cnt):
      p = path.pop()
      p.size -= 1
      p.valsize -= lmax_val
    for p in path:
      p.size -= 1
      p.valsize -= 1
    return True

  def discard(self, key: T, val=1) -> bool:
    if val <= 0:
      return True
    path = []
    node = self.root
    while node:
      path.append(node)
      if key == node.key:
        break
      node = node.left if key < node.key else node.right
    else:
      return False
    if val > node.val:
      val = node.val - 1
      if val > 0:
        node.val -= val
        while path:
          path.pop().valsize -= val
    if node.val == 1:
      self._discard(key)
    else:
      node.val -= val
      while path:
        path.pop().valsize -= val
    return True

  def remove(self, key: T, val: int=1) -> None:
    c = self.count(key)
    if c > val:
      raise KeyError(key)
    self.discard(key, val)

  def count(self, key: T) -> int:
    node = self.root
    while node:
      if key == node.key:
        return node.val
      node = node.left if key < node.key else node.right
    return 0

  def discard_all(self, key: T) -> bool:
    return self.discard(key, self.count(key))

  def le(self, key: T) -> Optional[T]:
    return BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node].le(self.root, key)

  def lt(self, key: T) -> Optional[T]:
    return BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node].lt(self.root, key)

  def ge(self, key: T) -> Optional[T]:
    return BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node].ge(self.root, key)

  def gt(self, key: T) -> Optional[T]:
    return BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node].gt(self.root, key)

  def index(self, key: T) -> int:
    return BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node].index(self.root, key)

  def index_right(self, key: T) -> int:
    return BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node].index_right(self.root, key)

  def index_keys(self, key: T) -> int:
    k = 0
    node = self.root
    while node:
      if key == node.key:
        if node.left:
          k += node.left.size
        break
      elif key < node.key:
        node = node.left
      else:
        k += node.val if node.left is None else node.left.size + node.val
        node = node.right
    return k

  def index_right_keys(self, key: T) -> int:
    k = 0
    node = self.root
    while node:
      if key == node.key:
        k += node.val if node.left is None else node.left.size + node.val
        break
      if key < node.key:
        node = node.left
      else:
        k += node.val if node.left is None else node.left.size + node.val
        node = node.right
    return k

  def pop(self, k: int=-1) -> T:
    if k < 0: k += self.root.valsize
    x = self[k]
    self.discard(x)
    return x

  def pop_min(self) -> T:
    return self.pop(0)

  def pop_max(self) -> T:
    return self.pop(-1)

  def items(self) -> Iterator[Tuple[T, int]]:
    for i in range(self.len_elm()):
      yield self._kth_elm_tree(i)

  def keys(self) -> Iterator[T]:
    for i in range(self.len_elm()):
      yield self._kth_elm_tree(i)[0]

  def values(self) -> Iterator[int]:
    for i in range(self.len_elm()):
      yield self._kth_elm_tree(i)[1]

  def show(self) -> None:
    print('{' + ', '.join(map(lambda x: f'{x[0]}: {x[1]}', self.tolist_items())) + '}')

  def get_elm(self, k: int) -> T:
    assert -self.len_elm() <= k < self.len_elm(), \
        f'IndexError: {self.__class__.__name__}.get_elm({k}), len_elm=({self.len_elm()})'
    return self._kth_elm_tree(k)[0]

  def len_elm(self) -> int:
    return self.root.size if self.root else 0

  def tolist(self) -> List[T]:
    return BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node].tolist(self.root, len(self))

  def tolist_items(self) -> List[Tuple[T, int]]:
    return BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node].tolist_items(self.root, self.len_elm())

  def clear(self) -> None:
    self.root = None

  def get_max(self) -> T:
    return self._kth_elm_tree(-1)[0]

  def get_min(self) -> T:
    return self._kth_elm_tree(0)[0]

  def __contains__(self, key: T):
    return BSTMultisetNodeBase[T, ScapegoatTreeMultiset.Node].contains(self.root, key)

  def __getitem__(self, k: int) -> T:
    assert -len(self) <= k < len(self), \
        f'IndexError: {self.__class__.__name__}[{k}], len={len(self)}'
    return self._kth_elm(k)[0]

  def __iter__(self):
    self.__iter = 0
    return self

  def __next__(self):
    if self.__iter == len(self):
      raise StopIteration
    res = self._kth_elm(self.__iter)[0]
    self.__iter += 1
    return res

  def __reversed__(self):
    for i in range(len(self)):
      yield self._kth_elm(-i-1)[0]

  def __len__(self):
    return self.root.valsize if self.root else 0

  def __bool__(self):
    return self.root is not None

  def __str__(self):
    return '{' + ', '.join(map(str, self.tolist())) + '}'

  def __repr__(self):
    return f'{self.__class__.__name__}({self.tolist})'

