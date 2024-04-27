# from titan_pylib.data_structures.b_tree.b_tree_list import BTreeList
from collections import deque
from typing import Generic, TypeVar, List, Iterable
try:
  from __pypy__ import newlist_hint
except ImportError:
  pass
T = TypeVar('T')

class BTreeList(Generic[T]):

  class _Node():

    def __init__(self):
      self.key: List[T] = []
      self.bit_len: List[int] = []
      self.key_len: int = 0
      self.child: List['BTreeList._Node'] = []
      self.size: int = 0
      self.sum: int = 0

    def is_leaf(self) -> bool:
      return not self.child

    def _add_size(self, s):
      self.size += s

    def split(self, i: int) -> 'BTreeList._Node':
      right = BTreeList._Node()
      self.key, right.key = self.key[:i], self.key[i:]
      self.child, right.child = self.child[:i+1], self.child[i+1:]
      size = len(self.key) + sum(cnode.size for cnode in self.child)
      s = sum(self.key) + sum(cnode.sum for cnode in self.child)
      right.sum = self.sum - s
      right.size = self.size - size
      self.sum = s
      self.size = size
      return right

    def insert_key(self, i: int, key: T, size: int=-1) -> None:
      self.size += 1
      self.sum += key if size == -1 else size
      self.key.insert(i, key)

    def insert_child(self, i: int, node: 'BTreeList._Node', size=-1) -> None:
      self.size += node.size if size == -1 else size
      self.sum += node.sum if size == -1 else size
      self.child.insert(i, node)

    def append_key(self, key: T) -> None:
      self.size += 1
      self.sum += key
      self.key.append(key)

    def append_child(self, node: 'BTreeList._Node', size: int=-1) -> None:
      self.size += node.size if size == -1 else size
      self.sum += node.sum if size == -1 else size
      self.child.append(node)

    def pop_key(self, i: int=-1) -> T:
      self.size -= 1
      v = self.key.pop(i)
      self.sum -= v
      return v

    def len_key(self) -> int:
      return len(self.key)

    def pop_child(self, i: int=-1, size: int=-1) -> 'BTreeList._Node':
      cnode = self.child.pop(i)
      self.sum -= cnode.sum if size == -1 else size
      self.size -= cnode.size if size == -1 else size
      return cnode

    def extend_key(self, keys: List[T]) -> None:
      self.size += len(keys)
      self.sum += sum(keys)
      self.key += keys

    def extend_child(self, children: List['BTreeList._Node']) -> None:
      self.size += sum(cnode.size for cnode in children)
      self.sum += sum(cnode.sum for cnode in children)
      self.child += children

    def __str__(self):
      return str((str(self.key), self.size, self.sum))

    __repr__ = __str__

  def __init__(self, a: Iterable[T]=[]):
    self._m = 300
    self._root = BTreeList._Node()
    self._len = 0
    self._build(a)

  def _build(self, a: Iterable[T]):
    for e in a:
      self.append(e)

  def __getitem__(self, k: int) -> T:
    node = self._root
    while True:
      if node.is_leaf():
        return node.key[k]
      for i in range(node.len_key()+1):
        if k < node.child[i].size:
          node = node.child[i]
          break
        k -= node.child[i].size
        if k == 0 and i < node.len_key():
          return node.key[i]
        k -= 1

  def pref(self, r: int) -> int:
    node = self._root
    s = 0
    while True:
      if node.is_leaf():
        s += sum(node.key[:r])
        break
      for i in range(node.len_key()+1):
        if r < node.child[i].size:
          node = node.child[i]
          break
        s += node.child[i].sum
        r -= node.child[i].size
        if r == 0 and i < node.len_key():
          break
        if i < node.len_key():
          s += node.key[i]
          r -= 1
    return s

  def prod(self, l: int, r: int) -> int:
    return self.pref(r) - self.pref(l)

  def _is_over(self, node: 'BTreeList._Node') -> bool:
    return node.len_key() > self._m

  def insert(self, k: int, key: T) -> bool:
    node = self._root
    stack = []
    while True:
      if node.is_leaf():
        node.insert_key(k, key)
        break
      for i in range(node.len_key()+1):
        if k < node.child[i].size:
          break
        k -= node.child[i].size
        if k == 0:
          k = node.child[i].size
          break
        k -= 1
      stack.append((node, i))
      node = node.child[i]
    self._len += 1
    while stack:
      if not self._is_over(node):
        break
      pnode, indx = stack.pop()
      i = node.len_key() // 2
      pre = node.sum
      center = node.pop_key(i)
      right = node.split(i)
      pnode.insert_key(indx, center, size=0)
      pnode.insert_child(indx+1, right, size=0)
      pnode.sum += key
      node = pnode
    while stack:
      pnode, _ = stack.pop()
      pnode._add_size(1)
      pnode.sum += key
    if self._is_over(node):
      pnode = BTreeList._Node()
      i = node.len_key() // 2
      center = node.pop_key(i)
      right = node.split(i)
      pnode.append_key(center)
      pnode.append_child(node)
      pnode.append_child(right)
      self._root = pnode
    return True

  def append(self, key: T) -> None:
    node = self._root
    stack = []
    while True:
      if node.is_leaf():
        node.append_key(key)
        break
      stack.append((node, node.len_key()))
      node = node.child[-1]
    self._len += 1
    while stack:
      if not self._is_over(node):
        break
      pnode, indx = stack.pop()
      i = node.len_key() // 2
      pre = node.sum
      center = node.pop_key(i)
      right = node.split(i)
      pnode.insert_key(indx, center, size=0)
      pnode.insert_child(indx+1, right, size=0)
      pnode.sum += key
      node = pnode
    while stack:
      pnode, _ = stack.pop()
      pnode._add_size(1)
      pnode.sum += key
    if self._is_over(node):
      pnode = BTreeList._Node()
      i = node.len_key() // 2
      center = node.pop_key(i)
      right = node.split(i)
      pnode.append_key(center)
      pnode.append_child(node)
      pnode.append_child(right)
      self._root = pnode
    return True

  def _discard_right(self, node: 'BTreeList._Node') -> T:
    stack = []
    while not node.is_leaf():
      stack.append(node)
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
    v = node.pop_key()
    self._update_stack(stack, v)
    return v

  def _discard_left(self, node: 'BTreeList._Node') -> T:
    stack = []
    while not node.is_leaf():
      stack.append(node)
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
    v = node.pop_key(0)
    self._update_stack(stack, v)
    return v

  def _merge(self, node: 'BTreeList._Node', i: int) -> 'BTreeList._Node':
    y = node.child[i]
    z = node.pop_child(i+1, size=0)
    v = node.pop_key(i)
    y.append_key(v)
    node.sum += v
    node.size += 1
    y.extend_key(z.key)
    y.extend_child(z.child)
    return y

  def _merge_key(self, node: 'BTreeList._Node', i: int) -> None:
    if node.child[i].len_key() > self._m//2:
      node.key[i] = self._discard_right(node.child[i])
      return
    if node.child[i+1].len_key() > self._m//2:
      node.key[i] = self._discard_left(node.child[i+1])
      return
    indx = node.child[i if i+1 < len(node.child) else (i-1)].size
    y = self._merge(node, i)
    self._pop(indx, y)
    if node is self._root and not node.key:
      self._root = y

  def _update_stack(self, stack, key):
    for s in stack:
      s.size -= 1
      s.sum -= key

  def tolist(self) -> List[T]:
    a = newlist_hint(len(self))
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

  def debug(self):
    dep = [[] for _ in range(10)]
    dq = deque([(self._root, 0)])
    while dq:
      node, d = dq.popleft()
      dep[d].append((node.key, node.size, node.sum))
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

  def pop(self, k: int) -> T:
    assert 0 <= k < len(self), f'IndexError'
    self._len -= 1
    return self._pop(k, self._root)

  def _pop(self, k: int, node: 'BTreeList._Node') -> T:
    stack = []
    while True:
      if node.is_leaf():
        v = node.pop_key(k)
        self._update_stack(stack, v)
        return v
      stack.append(node)
      for i in range(node.len_key()+1):
        if k < node.child[i].size:
          break
        k -= node.child[i].size
        if k == 0 and i < node.len_key():
          v = node.key[i]
          self._merge_key(node, i)
          self._update_stack(stack, v)
          return v
        k -= 1
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
          k += 1
          if cnode.child:
            k += cnode.child[-1].size
          node.child[i].insert_key(0, node.key[i-1])
          node.key[i-1] = cnode.pop_key()
          if cnode.child:
            node.child[i].insert_child(0, cnode.pop_child())
          node = node.child[i]
          continue
        if i+1 >= len(node.child):
          i -= 1
          k += node.child[i].size + 1
        cnode = self._merge(node, i)
        if node is self._root and not node.key:
          self._root = cnode
        node = cnode
        continue
      node = node.child[i]

  def __len__(self):
    return self._len

  def __str__(self):
    return str(self.tolist())
