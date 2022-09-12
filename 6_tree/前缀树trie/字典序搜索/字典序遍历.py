from collections import deque
from itertools import product
from typing import Generator


def genString() -> Generator[str, None, None]:
    """按照长度和字典序生成字符串"""
    len_ = 1
    while True:
        for cur in product(range(10), repeat=len_):
            yield "".join(map(str, cur))
        len_ += 1


gen_ = genString()
for i in range(20):
    # 0 1 2 3 4 5 6 7 8 9 00 01 02 03 04 05 06 07 08 09
    print(next(gen_), end=" ")


print()

#####################################################################
def genLexicalOrder1(n: int) -> Generator[int, None, None]:
    """字典序dfs遍历十叉树,生成[0,n]的整数"""

    def dfs(cur: int) -> Generator[int, None, None]:
        for i in range(10):
            next = cur * 10 + i
            if next == 0:
                continue
            if next > n:
                continue
            yield next
            yield from dfs(next)

    yield 0
    yield from dfs(0)


# !这个就是 for i in range(n)
def genLexicalOrder2(n: int) -> Generator[int, None, None]:
    """字典序bfs遍历十叉树,生成[0,n]的整数"""

    def bfs(start: int) -> Generator[int, None, None]:
        queue = deque([start])
        while queue:
            nextQueue = deque()
            step = len(queue)
            for _ in range(step):
                cur = queue.popleft()
                for i in range(10):
                    next = cur * 10 + i
                    if next == 0:
                        continue
                    if next > n:
                        continue
                    yield next
                    nextQueue.append(next)
            queue = nextQueue

    yield 0
    yield from bfs(0)


print(*genLexicalOrder1(10))
print(*genLexicalOrder2(10))
