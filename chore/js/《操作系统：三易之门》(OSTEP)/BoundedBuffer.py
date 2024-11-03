import threading
import time
import random


class BoundedBuffer:
    def __init__(self, max_size):
        self.buffer = []
        self.max_size = max_size
        self.lock = threading.Lock()
        self.empty_slots = threading.Semaphore(max_size)
        self.filled_slots = threading.Semaphore(0)

    def produce(self, item):
        self.empty_slots.acquire()
        with self.lock:
            self.buffer.append(item)
            print(f"生产者生产了: {item}")
        self.filled_slots.release()

    def consume(self):
        self.filled_slots.acquire()
        with self.lock:
            item = self.buffer.pop(0)
            print(f"消费者消费了: {item}")
        self.empty_slots.release()
        return item


def producer(bb, id, items):
    for item in items:
        bb.produce(f"item-{id}-{item}")
        time.sleep(random.uniform(0.1, 0.5))


def consumer(bb, id, consume_count):
    for _ in range(consume_count):
        item = bb.consume()
        time.sleep(random.uniform(0.1, 0.5))


if __name__ == "__main__":
    buffer_size = 5
    bb = BoundedBuffer(buffer_size)

    producers = [
        threading.Thread(target=producer, args=(bb, 1, range(5))),
        threading.Thread(target=producer, args=(bb, 2, range(5))),
    ]

    consumers = [
        threading.Thread(target=consumer, args=(bb, 1, 5)),
        threading.Thread(target=consumer, args=(bb, 2, 5)),
    ]

    for p in producers:
        p.start()
    for c in consumers:
        c.start()

    for p in producers:
        p.join()
    for c in consumers:
        c.join()

    print("生产和消费过程完成。")
