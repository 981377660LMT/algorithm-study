# 这里的队列（只是加了锁）通常不能完全满足这种程序的需求。
# 更完善的有界队列，在队列空或者满时，能让线程等待。

import threading


class Node:
    def __init__(self, value=None):
        self.value = value
        self.next = None


class ConcurrentQueue:
    def __init__(self):
        dummy = Node()
        self.head = dummy
        self.tail = dummy
        # !两个锁分别保护头部和尾部
        self.head_lock = threading.Lock()
        self.tail_lock = threading.Lock()

    def enqueue(self, value):
        new_node = Node(value)
        with self.tail_lock:
            self.tail.next = new_node  # type: ignore
            self.tail = new_node

    def dequeue(self):
        with self.head_lock:
            if self.head.next is None:
                return None  # 队列为空
            value = self.head.next.value
            self.head = self.head.next
            return value


if __name__ == "__main__":
    # 使用示例

    def producer(queue, items):
        for item in items:
            queue.enqueue(item)
            print(f"Enqueued: {item}")

    def consumer(queue, consume_count):
        for _ in range(consume_count):
            item = None
            while item is None:
                item = queue.dequeue()
            print(f"Dequeued: {item}")

    q = ConcurrentQueue()
    producer_thread = threading.Thread(target=producer, args=(q, range(10)))
    consumer_thread = threading.Thread(target=consumer, args=(q, 10))

    producer_thread.start()
    consumer_thread.start()

    producer_thread.join()
    consumer_thread.join()
