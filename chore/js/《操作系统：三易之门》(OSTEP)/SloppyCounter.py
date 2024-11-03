from threading import Lock, Thread
from collections import defaultdict


class LazyCounter:
    __slots__ = "global_count", "threshold", "lock", "local_counts"

    def __init__(self, threshold=10):
        self.global_count = 0
        self.threshold = threshold
        self.lock = Lock()
        self.local_counts = defaultdict(int)

    def increment(self, thread_id):
        self.local_counts[thread_id] += 1
        if self.local_counts[thread_id] >= self.threshold:
            with self.lock:
                self.global_count += self.local_counts[thread_id]
                self.local_counts[thread_id] = 0

    def get_value(self):
        with self.lock:
            total = self.global_count + sum(self.local_counts.values())
        return total


# 使用示例
def worker(counter, thread_id, increments):
    for _ in range(increments):
        counter.increment(thread_id)


if __name__ == "__main__":
    counter = LazyCounter(threshold=5)
    threads = []
    for i in range(4):
        t = Thread(target=worker, args=(counter, i, 20))
        threads.append(t)
        t.start()
    for t in threads:
        t.join()
    print("最终计数:", counter.get_value())
