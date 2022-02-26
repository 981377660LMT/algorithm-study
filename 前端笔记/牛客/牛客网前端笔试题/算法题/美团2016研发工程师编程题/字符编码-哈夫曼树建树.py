# 请设计一个算法，给一个字符串进行二进制编码，使得编码后字符串的长度最短。
# 一行输出最短的编码后长度。
# n<=1000

# 解法一：直接模拟建立哈夫曼树
from collections import Counter
from heapq import heapify, heappop, heappush


class Node:
    def __init__(self, weight: int, value: str = None, left: 'Node' = None, right: 'Node' = None):
        self.weight = weight
        self.value = value
        self.left = left
        self.right = right

    def __lt__(self, other: 'Node') -> bool:
        return self.weight < other.weight

    def __eq__(self, other: 'Node') -> bool:
        return self.weight == other.weight


def main1(string: str) -> None:
    def dfs(root: Node, depth: int) -> None:
        nonlocal res
        if not root:
            return
        if root.value is not None:
            res += root.weight * depth
        root.left and dfs(root.left, depth + 1)
        root.right and dfs(root.right, depth + 1)

    chars = list(string)
    counter = Counter(chars)

    pq = []
    for value, weight in counter.items():
        pq.append((weight, Node(weight, value)))
    heapify(pq)

    while len(pq) >= 2:
        _, left = heappop(pq)
        _, right = heappop(pq)
        parent = Node(left.weight + right.weight, None, left, right)
        heappush(pq, (parent.weight, parent))

    root = pq[0][1]
    res = 0
    dfs(root, 0)
    print(res)


# 只需要push字符数即可，不必实际建树
def main2(string: str) -> None:
    chars = list(string)
    counter = Counter(chars)

    pq = []
    for _, weight in counter.items():
        pq.append((weight))
    heapify(pq)

    res = 0
    while len(pq) >= 2:
        left = heappop(pq)
        right = heappop(pq)
        res += left + right
        heappush(pq, left + right)

    print(res)


while True:
    try:
        string = input()
        main2(string)
    except EOFError:
        break

