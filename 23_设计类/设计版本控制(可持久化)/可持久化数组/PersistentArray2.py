"""完全可持久化数组 动态开点"""


from typing import List, Optional, Union


class _Node:
    __slots__ = ("left", "right", "leftChild", "rightChild", "value")

    @staticmethod
    def create(left: int, right: int) -> "_Node":
        """创建管理[left, right)的节点"""
        if left == right - 1:
            return _Node(left, right, None, None, 0)
        mid = (left + right) // 2
        return _Node(left, right, _Node.create(left, mid), _Node.create(mid, right), None)

    def __init__(
        self,
        left: int,
        right: int,
        leftChild: Optional["_Node"],
        rightChild: Optional["_Node"],
        value: Optional[int],
    ) -> None:
        self.left = left
        self.right = right
        self.leftChild = leftChild
        self.rightChild = rightChild
        self.value = value


class PersistentArray:
    __slots__ = ("_root", "_length")

    @staticmethod
    def create(lengthOrArray: Union[int, List[int]]) -> "PersistentArray":
        """Create a PersistentArray from length or array."""
        isArray = isinstance(lengthOrArray, list)
        n = len(lengthOrArray) if isArray else lengthOrArray
        assert n > 0, f"length must be positive, but {n} received"
        root = _Node.create(0, n)
        res = PersistentArray(root, n)
        if isArray:
            PersistentArray._build(res._root, lengthOrArray)
        return res

    @staticmethod
    def _get(node: "_Node", index: int) -> int:
        if node.value is not None:
            return node.value

        mid = (node.left + node.right) // 2
        if index < mid:
            return PersistentArray._get(node.leftChild, index)  # type: ignore
        return PersistentArray._get(node.rightChild, index)  # type: ignore

    @staticmethod
    def _update(node: "_Node", index: int, value: int) -> "_Node":
        left, right = node.left, node.right
        if left == right - 1:
            return _Node(left, right, None, None, value)

        mid = (left + right) // 2
        if index < mid:
            return _Node(
                left,
                right,
                PersistentArray._update(node.leftChild, index, value),  # type: ignore
                node.rightChild,
                None,
            )
        return _Node(
            left,
            right,
            node.leftChild,
            PersistentArray._update(node.rightChild, index, value),  # type: ignore
            None,
        )

    @staticmethod
    def _build(node: "_Node", array: List[int]) -> None:
        left, right = node.left, node.right
        if left == right - 1:
            node.value = array[left]
            return
        PersistentArray._build(node.leftChild, array)  # type: ignore
        PersistentArray._build(node.rightChild, array)  # type: ignore

    def __init__(self, root: "_Node", length: int) -> None:
        self._root = root
        self._length = length

    def get(self, index: int) -> int:
        """Get the value at index."""
        assert 0 <= index < self._length, f"index out of range: {index}"
        return PersistentArray._get(self._root, index)

    def update(self, index: int, value: int) -> "PersistentArray":
        """Update the value at index and return a new PersistentArray."""
        assert 0 <= index < self._length, f"index out of range: {index}"
        node = PersistentArray._update(self._root, index, value)
        return PersistentArray(node, self._length)

    def __repr__(self) -> str:
        def inOrder(node: Optional["_Node"]) -> None:
            if node is None:
                return
            inOrder(node.leftChild)
            if node.value is not None:
                res.append(node.value)
                return
            inOrder(node.rightChild)

        res = []
        inOrder(self._root)
        return f"{self.__class__.__name__}({res})"

    def __len__(self) -> int:
        return self._length

    def __getitem__(self, index: int) -> int:
        raise NotImplementedError(f"use {self.__class__.__name__}.get instead")

    def __setitem__(self, index: int, value: int) -> None:
        raise NotImplementedError(f"use {self.__class__.__name__}.update instead")


if __name__ == "__main__":
    arr = PersistentArray.create([1, 2, 3, 4, 5, 6])
    print(arr)
    arr = arr.update(1, 4)
    print(arr)
    print(arr.get(3))
