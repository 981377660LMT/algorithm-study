from typing import Generic, TypeVar, List

T = TypeVar("T")


class RingBuffer(Generic[T]):
    """
    固定大小的环形缓冲区，可用于双端队列、固定大小滑动窗口等场景。
    """

    __slots__ = ("_values", "_start", "_end", "_maxSize", "_size")

    def __init__(self, maxSize: int):
        if maxSize < 1:
            raise ValueError("Invalid maxSize, should be at least 1")
        self._values: List[T] = [None] * maxSize  # type: ignore
        self._start = 0
        self._end = 0
        self._maxSize = maxSize
        self._size = 0

    def append(self, value: T) -> None:
        """
        在队列尾部添加元素。如果队列已满，则自动移除头部元素。
        """
        if self.full():
            self.popleft()
        self._values[self._end] = value
        self._end += 1
        if self._end >= self._maxSize:
            self._end = 0
        self._size += 1

    def appendleft(self, value: T) -> None:
        """
        在队列头部添加元素。如果队列已满，则自动移除尾部元素。
        """
        if self.full():
            self.pop()
        self._start -= 1
        if self._start < 0:
            self._start = self._maxSize - 1
        self._values[self._start] = value
        self._size += 1

    def pop(self) -> T:
        """
        移除并返回队列尾部的元素。如果队列为空，则引发异常。
        """
        if self.empty():
            raise IndexError("RingBuffer is empty")
        self._end -= 1
        if self._end < 0:
            self._end = self._maxSize - 1
        value = self._values[self._end]
        self._size -= 1
        return value

    def popleft(self) -> T:
        """
        移除并返回队列头部的元素。如果队列为空，则引发异常。
        """
        if self.empty():
            raise IndexError("RingBuffer is empty")
        value = self._values[self._start]
        self._start += 1
        if self._start >= self._maxSize:
            self._start = 0
        self._size -= 1
        return value

    def head(self) -> T:
        if self.empty():
            raise IndexError("RingBuffer is empty")
        return self._values[self._start]

    def tail(self) -> T:
        if self.empty():
            raise IndexError("RingBuffer is empty")
        index = self._end - 1
        if index < 0:
            index = self._maxSize - 1
        return self._values[index]

    def empty(self) -> bool:
        return self._size == 0

    def full(self) -> bool:
        return self._size == self._maxSize

    def clear(self) -> None:
        self._start = 0
        self._end = 0
        self._size = 0

    def __len__(self) -> int:
        return self._size

    def __iter__(self):
        ptr = self._start
        for _ in range(self._size):
            yield self._values[ptr]
            ptr += 1
            if ptr >= self._maxSize:
                ptr = 0

    def __getitem__(self, index: int) -> T:
        size = self._size
        if index < 0:
            index += size
        if index < 0 or index >= size:
            raise IndexError("Index out of range")
        index += self._start
        if index >= self._maxSize:
            index -= self._maxSize
        return self._values[index]

    def __setitem__(self, index: int, value: T) -> None:
        size = self._size
        if index < 0:
            index += size
        if index < 0 or index >= size:
            raise IndexError("Index out of range")
        index += self._start
        if index >= self._maxSize:
            index -= self._maxSize
        self._values[index] = value

    def __repr__(self) -> str:
        return f"RingBuffer: {', '.join(map(str, self))}"


if __name__ == "__main__":
    #  641. 设计循环双端队列
    #  https://leetcode.cn/problems/design-circular-deque/description/
    class MyCircularDeque:
        __slots__ = ("_k", "_queue")

        def __init__(self, k: int):
            self._k = k
            self._queue = RingBuffer(k)

        def insertFront(self, value: int) -> bool:
            if self.isFull():
                return False
            self._queue.appendleft(value)
            return True

        def insertLast(self, value: int) -> bool:
            if self.isFull():
                return False
            self._queue.append(value)
            return True

        def deleteFront(self) -> bool:
            if self.isEmpty():
                return False
            self._queue.popleft()
            return True

        def deleteLast(self) -> bool:
            if self.isEmpty():
                return False
            self._queue.pop()
            return True

        def getFront(self) -> int:
            if self.isEmpty():
                return -1
            return self._queue.head()

        def getRear(self) -> int:
            if self.isEmpty():
                return -1
            return self._queue.tail()

        def isEmpty(self) -> bool:
            return self._queue.empty()

        def isFull(self) -> bool:
            return self._queue.full()
