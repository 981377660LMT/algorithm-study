# // 还有nlogn的解法 hu_tucker算法
# https://atcoder.jp/contests/atc002/submissions/18244148


from heapq import heappop, heappush, heapify

INF = int(1e18)


class Heap:
    __slots__ = ("val", "lt", "rt")

    def __init__(self, val):
        self.val = val
        self.lt = None
        self.rt = None


def meld(a, b):
    if a is None:
        return b
    if b is None:
        return a
    if a.val > b.val:
        a, b = b, a
    a.rt = meld(a.rt, b)
    a.lt, a.rt = a.rt, a.lt
    return a


def top(a):
    return a.val


def pop(a):
    return meld(a.lt, a.rt)


def push(a, x):
    b = Heap(x)
    return meld(a, b)


def hu_tucker(n, arr):
    w = list(arr)
    lt = [0] * n
    rt = [0] * n
    cost = [0] * (n - 1)
    heap = [None for _ in range(n - 1)]
    queue = []
    for i in range(n - 1):
        lt[i] = i - 1
        rt[i] = i + 1
        cost[i] = w[i] + w[i + 1]
        queue.append(cost[i] * n + i)
    heapify(queue)
    res = 0
    for _ in range(n - 1):
        while True:
            p = heappop(queue)
            c, i = divmod(p, n)
            if cost[i] == c and rt[i] >= 0:
                break
        ml = mr = False
        if heap[i] is not None and w[i] + heap[i].val == c:
            heap[i] = pop(heap[i])
            ml = True
        elif w[i] + w[rt[i]] == c:
            ml = mr = True
        else:
            t = top(heap[i])
            heap[i] = pop(heap[i])
            if heap[i] is not None and top(heap[i]) + t == c:
                heap[i] = pop(heap[i])
            else:
                mr = True
        res += c
        heap[i] = push(heap[i], c)
        if ml:
            w[i] = INF
        if mr:
            w[rt[i]] = INF
        if ml and i > 0:
            j = lt[i]
            heap[j] = meld(heap[i], heap[j])
            rt[j] = rt[i]
            rt[i] = -1
            lt[rt[j]] = j
            i = j
        if mr and rt[i] + 1 < n:
            j = rt[i]
            heap[i] = meld(heap[i], heap[j])
            rt[i] = rt[j]
            rt[j] = -1
            lt[rt[i]] = i
        cost[i] = w[i] + w[rt[i]]
        if heap[i] is not None:
            t = top(heap[i])
            heap[i] = pop(heap[i])
            cost[i] = min(cost[i], w[i] + t, w[rt[i]] + t)
            if heap[i] is not None:
                cost[i] = min(cost[i], top(heap[i]) + t)
            heap[i] = push(heap[i], t)
        heappush(queue, cost[i] * n + i)
    return res


import sys

input = sys.stdin.buffer.readline
print(hu_tucker(int(input()), map(int, input().split())))
