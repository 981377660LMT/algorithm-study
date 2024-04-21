# from titan_pylib.data_structures.b_tree.b_tree_set import BTreeSet
# from titan_pylib.my_class.ordered_set_interface import OrderedSetInterface
# from titan_pylib.my_class.supports_less_than import SupportsLessThan
from typing import Protocol

class SupportsLessThan(Protocol):

  def __lt__(self, other) -> bool: ...

from abc import ABC, abstractmethod
from typing import Iterable, Optional, Iterator, TypeVar, Generic, List
T = TypeVar('T', bound=SupportsLessThan)

class OrderedSetInterface(ABC, Generic[T]):

  @abstractmethod
  def __init__(self, a: Iterable[T]) -> None:
    raise NotImplementedError

  @abstractmethod
  def add(self, key: T) -> bool:
    raise NotImplementedError

  @abstractmethod
  def discard(self, key: T) -> bool:
    raise NotImplementedError

  @abstractmethod
  def remove(self, key: T) -> None:
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

from collections import deque
from bisect import bisect_left, bisect_right, insort
from typing import Deque, Generic, Tuple, TypeVar, List, Optional, Iterable
T = TypeVar('T', bound=SupportsLessThan)

class BTreeSet(OrderedSetInterface, Generic[T]):

  class _Node():

    def __init__(self):
      self.key: List = []
      self.child: List['BTreeSet._Node'] = []

    def is_leaf(self) -> bool:
      return not self.child

    def split(self, i: int) -> 'BTreeSet._Node':
      right = BTreeSet._Node()
      self.key, right.key = self.key[:i], self.key[i:]
      self.child, right.child = self.child[:i+1], self.child[i+1:]
      return right

    def insert_key(self, i: int, key: T) -> None:
      self.key.insert(i, key)

    def insert_child(self, i: int, node: 'BTreeSet._Node') -> None:
      self.child.insert(i, node)

    def append_key(self, key: T) -> None:
      self.key.append(key)

    def append_child(self, node: 'BTreeSet._Node') -> None:
      self.child.append(node)

    def pop_key(self, i: int=-1) -> T:
      return self.key.pop(i)

    def len_key(self) -> int:
      return len(self.key)

    def insort_key(self, key: T) -> None:
      insort(self.key, key)

    def pop_child(self, i: int=-1) -> 'BTreeSet._Node':
      return self.child.pop(i)

    def extend_key(self, keys: List[T]) -> None:
      self.key += keys

    def extend_child(self, children: List['BTreeSet._Node']) -> None:
      self.child += children

    def __str__(self):
      return str(str(self.key))

    __repr__ = __str__

  def __init__(self, a: Iterable[T]=[]):
    self._m: int = 1000
    self._root: 'BTreeSet._Node' = BTreeSet._Node()
    self._len: int = 0
    self._build(a)

  def _build(self, a: Iterable[T]):
    for e in a:
      self.add(e)

  def _is_over(self, node: 'BTreeSet._Node') -> bool:
    return node.len_key() > self._m

  def add(self, key: T) -> bool:
    node = self._root
    stack = []
    while True:
      i = bisect_left(node.key, key)
      if i < node.len_key() and node.key[i] == key:
        return False
      if i >= len(node.child):
        break
      stack.append(node)
      node = node.child[i]
    self._len += 1
    node.insort_key(key)
    while stack:
      if not self._is_over(node):
        break
      pnode = stack.pop()
      i = node.len_key() // 2
      center = node.pop_key(i)
      right = node.split(i)
      indx = bisect_left(pnode.key, center)
      pnode.insert_key(indx, center)
      pnode.insert_child(indx+1, right)
      node = pnode
    if self._is_over(node):
      pnode = BTreeSet._Node()
      i = node.len_key() // 2
      center = node.pop_key(i)
      right = node.split(i)
      pnode.append_key(center)
      pnode.append_child(node)
      pnode.append_child(right)
      self._root = pnode
    return True

  def __contains__(self, key: T) -> bool:
    node = self._root
    while True:
      i = bisect_left(node.key, key)
      if i < node.len_key() and node.key[i] == key:
        return True
      if node.is_leaf():
        break
      node = node.child[i]
    return False

  def _discard_right(self, node: 'BTreeSet._Node') -> T:
    while not node.is_leaf():
      if node.child[-1].len_key() == self._m//2:
        if node.child[-2].len_key() > self._m//2:
          cnode = node.child[-2]
          node.child[-1].insert_key(0, node.key[-1])
          node.key[-1] = cnode.pop_key()
          if cnode.child:
            node.child[-1].insert_child(0, cnode.pop_child())
          node = node.child[-1]
          continue
        cnode = self._merge(node, node.len_key()-1)
        if node is self._root and not node.key:
          self._root = cnode
        node = cnode
        continue
      node = node.child[-1]
    return node.pop_key()

  def _discard_left(self, node: 'BTreeSet._Node') -> T:
    while not node.is_leaf():
      if node.child[0].len_key() == self._m//2:
        if node.child[1].len_key() > self._m//2:
          cnode = node.child[1]
          node.child[0].append_key(node.key[0])
          node.key[0] = cnode.pop_key(0)
          if cnode.child:
            node.child[0].append_child(cnode.pop_child(0))
          node = node.child[0]
          continue
        cnode = self._merge(node, 0)
        if node is self._root and not node.key:
          self._root = cnode
        node = cnode
        continue
      node = node.child[0]
    return node.pop_key(0)

  def _merge(self, node: 'BTreeSet._Node', i: int) -> 'BTreeSet._Node':
    y = node.child[i]
    z = node.pop_child(i+1)
    y.append_key(node.pop_key(i))
    y.extend_key(z.key)
    y.extend_child(z.child)
    return y

  def _merge_key(self, key: T, node: 'BTreeSet._Node', i: int) -> None:
    if node.child[i].len_key() > self._m//2:
      node.key[i] = self._discard_right(node.child[i])
      return
    if node.child[i+1].len_key() > self._m//2:
      node.key[i] = self._discard_left(node.child[i+1])
      return
    y = self._merge(node, i)
    self._discard(key, y)
    if node is self._root and not node.key:
      self._root = y

  def _discard(self, key: T, node: Optional['BTreeSet._Node']=None) -> bool:
    if node is None:
      node = self._root
    if not node.key:
      return False
    while True:
      i = bisect_left(node.key, key)
      if node.is_leaf():
        if i < node.len_key() and node.key[i] == key:
          node.pop_key(i)
          return True
        return False
      if i < node.len_key() and node.key[i] == key:
        assert i+1 < len(node.child)
        self._merge_key(key, node, i)
        return True
      if node.child[i].len_key() == self._m//2:
        if i+1 < len(node.child) and node.child[i+1].len_key() > self._m//2:
          cnode = node.child[i+1]
          node.child[i].append_key(node.key[i])
          node.key[i] = cnode.pop_key(0)
          if cnode.child:
            node.child[i].append_child(cnode.pop_child(0))
          node = node.child[i]
          continue
        if i-1 >= 0 and node.child[i-1].len_key() > self._m//2:
          cnode = node.child[i-1]
          node.child[i].insert_key(0, node.key[i-1])
          node.key[i-1] = cnode.pop_key()
          if cnode.child:
            node.child[i].insert_child(0, cnode.pop_child())
          node = node.child[i]
          continue
        if i+1 >= len(node.child):
          i -= 1
        cnode = self._merge(node, i)
        if node is self._root and not node.key:
          self._root = cnode
        node = cnode
        continue
      node = node.child[i]

  def discard(self, key: T) -> bool:
    if self._discard(key):
      self._len -= 1
      return True
    return False

  def remove(self, key: T) -> None:
    if self.discard(key):
      return
    raise ValueError

  def tolist(self) -> List[T]:
    a = []
    def dfs(node):
      if not node.child:
        a.extend(node.key)
        return
      dfs(node.child[0])
      for i in range(node.len_key()):
        a.append(node.key[i])
        dfs(node.child[i+1])
    dfs(self._root)
    return a

  def get_max(self) -> Optional[T]:
    node = self._root
    while True:
      if not node.child:
        return node.key[-1] if node.key else None
      node = node.child[-1]

  def get_min(self) -> Optional[T]:
    node = self._root
    while True:
      if not node.child:
        return node.key[0] if node.key else None
      node = node.child[0]

  def debug(self) -> None:
    dep = [[] for _ in range(10)]
    dq: Deque[Tuple['BTreeSet._Node', int]] = deque([(self._root, 0)])
    while dq:
      node, d = dq.popleft()
      dep[d].append(node.key)
      if node.child:
        print(node, 'child=', node.child)
      for e in node.child:
        if e:
          dq.append((e, d+1))
    for i in range(10):
      if not dep[i]: break
      for e in dep[i]:
        print(e, end='  ')
      print()

  def pop_max(self) -> T:
    res = self.get_max()
    assert res is not None, f'IndexError: pop_max from empty {self.__class__.__name__}.'
    self.discard(res)
    return res

  def pop_min(self) -> T:
    res = self.get_min()
    assert res is not None, f'IndexError: pop_min from empty {self.__class__.__name__}.'
    self.discard(res)
    return res

  def ge(self, key: T) -> Optional[T]:
    res, node = None, self._root
    while node.key:
      i = bisect_left(node.key, key)
      if i < node.len_key() and node.key[i] == key:
        return node.key[i]
      if i < node.len_key():
        res = node.key[i]
      if not node.child:
        break
      node = node.child[i]
    return res

  def gt(self, key: T) -> Optional[T]:
    res, node = None, self._root
    while node.key:
      i = bisect_right(node.key, key)
      if i < node.len_key():
        res = node.key[i]
      if not node.child:
        break
      node = node.child[i]
    return res

  def le(self, key: T) -> Optional[T]:
    res, node = None, self._root
    while node.key:
      i = bisect_left(node.key, key)
      if i < node.len_key() and node.key[i] == key:
        return node.key[i]
      if i-1 >= 0:
        res = node.key[i-1]
      if not node.child:
        break
      node = node.child[i]
    return res

  def lt(self, key: T) -> Optional[T]:
    res, node = None, self._root
    while node.key:
      i = bisect_left(node.key, key)
      if i-1 >= 0:
        res = node.key[i-1]
      if not node.child:
        break
      node = node.child[i]
    return res

  def clear(self) -> None:
    self._root = BTreeSet._Node()

  def __iter__(self):
    self._iter_val = self.get_min()
    return self

  def __next__(self):
    if self._iter_val is None:
      raise StopIteration
    p = self._iter_val
    self._iter_val = self.gt(self._iter_val)
    return p

  def __bool__(self):
    return self._len > 0

  def __len__(self):
    return self._len

  def __str__(self):
    return '{' + ', '.join(map(str, self.tolist())) + '}'

  def __repr__(self):
    return f'{self.__class__.__name__}({self.tolist()})'
