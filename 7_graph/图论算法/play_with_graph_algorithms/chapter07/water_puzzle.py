from collections import deque


class WaterPuzzle:

    def __init__(self):
        queue = deque()
        visited = [False] * 100
        self._pre = [0] * 100
        self._end = -1

        queue.append(0)
        visited[0] = True
        while queue:
            cur = queue.popleft()
            a, b = cur // 10, cur % 10
            # max a = 5, max b = 3
            nexts = []
            nexts.append(5 * 10 + b)
            nexts.append(a * 10 + 3)
            nexts.append(0 * 10 + b)
            nexts.append(a * 10 + 0)

            # from a, minus x
            x = min(a, 3 - b)
            nexts.append((a - x) * 10 + (b + x))
            # from b, minux y
            y = min(5 - a, b)
            nexts.append((a + y) * 10 + (b - y))

            for next_ in nexts:
                if not visited[next_]:
                    queue.append(next_)
                    visited[next_] = True
                    self._pre[next_] = cur
                    if next_ // 10 == 4 or next_ % 10 == 4:
                        self._end = next_

    def result(self):
        res = []
        if self._end == -1:
            return []
        cur = self._end
        while cur != 0:
            res.append(cur)
            cur = self._pre[cur]
        res.append(0)
        return res[::-1]

if __name__ == '__main__':
    prob = WaterPuzzle()
    print(prob.result())