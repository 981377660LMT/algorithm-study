from collections import defaultdict, deque


class Logger:
    def __init__(self):
        self.queues = defaultdict(deque)

    # 每条 `不重复` 的消息最多只能每 10 秒打印一次
    # 每个 timestamp 都将按非递减顺序（时间顺序）传递
    def shouldPrintMessage(self, timestamp: int, message: str) -> bool:
        queue = self.queues[message]
        if not queue:
            queue.append((timestamp, message))
            return True
        if timestamp - queue[-1][0] >= 10:
            queue.append((timestamp, message))
            while queue and timestamp - queue[0][0] >= 10:
                queue.popleft()
            return True

        return False

