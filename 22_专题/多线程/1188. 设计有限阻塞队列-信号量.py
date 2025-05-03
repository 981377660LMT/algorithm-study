from collections import deque
from threading import Semaphore


class BoundedBlockingQueue:

    def __init__(self, capacity: int):
        self.left_capacity = Semaphore(capacity)
        self.curr_capacity = Semaphore(0)
        self.data = deque()

    def enqueue(self, element: int) -> None:
        self.left_capacity.acquire()
        self.curr_capacity.release()
        self.data.append(element)

    def dequeue(self) -> int:
        self.left_capacity.release()
        self.curr_capacity.acquire()
        return self.data.popleft()

    def size(self) -> int:
        return len(self.data)
