from typing import Generic, Iterable, Optional, Tuple, TypeVar


_V = TypeVar("_V", int, str)


class _SplayNode(Generic[_V]):
    __slots__ = "value", "left", "right", "parent", "size"

    def __init__(self, value: _V):
        self.value = value
        self.left, self.right, self.parent = None, None, None
        self.size = 1

    def update(self) -> None:
        self.size = 1
        if self.left is not None:
            self.size += self.left.size
        if self.right is not None:
            self.size += self.right.size

    def getState(self) -> int:
        if self.parent is None:
            return 0
        if self.parent.left == self:
            return 1
        else:
            return -1  # self.parent.right == self

    def splay(self) -> "_SplayNode[_V]":
        while self.parent is not None:
            curState, parentState = self.getState(), self.parent.getState()
            if parentState == 0:
                if curState == 1:
                    self.parent._rotateRight()
                else:
                    self.parent._rotateLeft()
            elif parentState == 1:
                if curState == 1:
                    self.parent.parent._rotateRight()
                    self.parent._rotateRight()
                else:
                    self.parent._rotateLeft()
                    self.parent._rotateRight()
            else:
                if curState == 1:
                    self.parent._rotateRight()
                    self.parent._rotateLeft()
                else:
                    self.parent.parent._rotateLeft()
                    self.parent._rotateLeft()
        return self

    def _rotateLeft(self) -> None:
        node = self.right
        node.parent = self.parent
        if node.parent is not None:
            if node.parent.left is self:
                node.parent.left = node
            else:
                node.parent.right = node
        self.right = node.left
        if self.right is not None:
            self.right.parent = self
        self.parent = node
        node.left = self
        self.update()
        node.update()

    def _rotateRight(self) -> None:
        node = self.left
        node.parent = self.parent
        if node.parent is not None:
            if node.parent.right is self:
                node.parent.right = node
            else:
                node.parent.left = node
        self.left = node.right
        if self.left is not None:
            self.left.parent = self
        self.parent = node
        node.right = self
        self.update()
        node.update()


class SplayArray(Generic[_V]):
    """增删改查的时间复杂度都是O(logn)的数组"""

    __slots__ = "_root"

    def __init__(self, iterable: Optional[Iterable[_V]] = None):
        self._root = None
        if iterable is not None:
            for item in iterable:
                self.insert(len(self), item)

    def insert(self, index: int, value: _V) -> None:
        """Insert object before index."""
        left, right = self._split(index)
        node = _SplayNode(value)
        self._root = self._merge(self._merge(left, node), right)

    def append(self, value: _V) -> None:
        """Append object to the end of the list."""
        self.insert(len(self), value)

    def pop(self, index: int) -> None:
        """Remove and return item at index."""
        if index < 0:
            index += len(self)
        self._splay(index)
        left = self._root.left
        right = self._root.right
        if left is not None:
            left.parent = None
        if right is not None:
            right.parent = None
        self._root = self._merge(left, right)

    def index(self, value: _V) -> int:
        """Return first index of value.

        Raises ValueError if the value is not present.
        """
        node = self._root
        index = 0
        try:
            while True:
                if value < node.value:
                    node = node.left
                elif value > node.value:
                    index += node.left.size + 1
                    node = node.right
                else:
                    index += node.left.size
                    return index
        except AttributeError:
            raise ValueError(f"{value} is not in list")

    def remove(self, value: _V) -> None:
        """Remove first occurrence of value.

        Raises ValueError if the value is not present.
        """
        index = self.index(value)
        self.pop(index)

    def _splay(self, index: int) -> None:
        """Splay the node at index to the root."""
        node = self._root
        while True:
            leftSize = node.left.size if node.left is not None else 0
            if index < leftSize:
                node = node.left
            elif index > leftSize:
                node = node.right
                index -= leftSize + 1
            else:
                self._root = node.splay()
                break

    def _merge(self, left: "_SplayNode[_V]", right: "_SplayNode[_V]") -> "_SplayNode[_V]":
        if left is None:
            return right
        if right is None:
            return left
        while left.right is not None:
            left = left.right
        left = left.splay()
        left.right = right
        right.parent = left
        left.update()
        return left

    def _split(
        self, leftCount: int
    ) -> Tuple[Optional["_SplayNode[_V]"], Optional["_SplayNode[_V]"]]:
        if leftCount == 0:
            return None, self._root
        if leftCount == self._root.size:
            return self._root, None
        self._splay(leftCount)
        left = self._root.left
        right = self._root
        left.parent = None
        right.left = None
        right.update()
        return left, right

    def _inorder(self, node: Optional["_SplayNode[_V]"]):
        """splay的中序遍历是当前维护的序列"""
        if node is None:
            return
        yield from self._inorder(node.left)
        yield node.value
        yield from self._inorder(node.right)

    def __getitem__(self, index: int):
        if index < 0:
            index += len(self)
        self._splay(index)
        return self._root.value

    def __setitem__(self, index: int, value: _V):
        if index < 0:
            index += len(self)
        self._splay(index)
        self._root.value = value

    def __delitem__(self, index: int):
        if index < 0:
            index += len(self)
        self.pop(index)

    def __contains__(self, value: _V):
        try:
            self.index(value)
        except ValueError:
            return False
        return True

    def __iter__(self):
        return self._inorder(self._root)

    def __len__(self):
        return self._root.size if self._root is not None else 0

    def __repr__(self) -> str:
        return f"SplayArray({list(self)})"


if __name__ == "__main__":
    import time

    step = int(2e5)  # !数据量超过1e5时，splay的效率才能体现出来
    print(f"step = {step}")

    print()
    print(f"测试SplayArray,初始长度{step}")
    start = time.time()
    nums = SplayArray(range(step))
    delta = time.time() - start
    print(f"初始化耗时 {(delta*1000):.2f}ms")

    # 插入删除step次
    start = time.time()
    for i in range(step):
        nums.insert(i, i)
        nums.pop(i)
    delta = time.time() - start
    print(f"插入删除耗时 {(delta*1000):.2f}ms")

    # 访问step次
    start = time.time()
    for i in range(step):
        nums[i]
    delta = time.time() - start
    print(f"随机访问耗时 {(delta*1000):.2f}ms")
    ##################################################
    print()
    print(f"测试普通数组,初始长度{step}")
    start = time.time()
    nums = list(range(step))
    delta = time.time() - start
    print(f"初始化耗时 {(delta*1000):.2f}ms")

    # 插入删除step次
    start = time.time()
    for i in range(step):
        nums[i:i] = [i]
        nums[i : i + 1] = []
    delta = time.time() - start
    print(f"插入删除耗时 {(delta*1000):.2f}ms")

    # 访问step次
    start = time.time()
    for i in range(step):
        nums[i]
    delta = time.time() - start
    print(f"随机访问耗时 {(delta*1000):.2f}ms")
