# !类似 Java BlockingQueue 源码实现
#
# enqueue 生产者队列
# 加锁对资源进行保护，判断队列 full 状态，进入wait
#
# dequeue 消费者队列
# 加锁对资源进行保护，判断队列 empty 状态，进入wait


from collections import deque
from threading import Condition, Lock


class BoundedBlockingQueue:

    def __init__(self, capacity: int):
        self.maxsize = capacity
        self.data = deque()
        self.mutex = Lock()
        self.not_empty = Condition(self.mutex)
        self.not_full = Condition(self.mutex)

    def enqueue(self, element: int) -> None:
        with self.not_full:
            while len(self.data) == self.maxsize:
                self.not_full.wait()
            self.data.append(element)
            self.not_empty.notify()

    def dequeue(self) -> int:
        with self.not_empty:
            while len(self.data) == 0:
                self.not_empty.wait()
            res = self.data.popleft()
            self.not_full.notify()
            return res

    def size(self) -> int:
        return len(self.data)
