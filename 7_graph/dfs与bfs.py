# !输出各个点的值


from collections import defaultdict, deque
from typing import Deque

DAG = defaultdict(set)
edges = [[0, 1], [0, 2], [1, 3], [1, 4], [2, 5], [2, 6]]
for u, v in edges:
    DAG[u].add(v)


def solve1(queue: Deque[int] = deque([0])):
    """dfs模拟bfs 结束时使用生成器中断函数执行,阻止继续回溯"""
    len_ = len(queue)
    if len(queue) == 0:
        print()
        yield 'end'
    for _ in range(len_):  # !逐个元素作为当前队头看一遍 (取出 popleft 回溯时 append)
        cur = queue.popleft()
        print(cur, end=' ')
        for next in DAG[cur]:
            queue.append(next)
        yield from solve1(queue)
        queue.append(cur)


def solve2(queue: Deque[int] = deque([0])) -> None:
    """bfs模拟dfs 队列当作栈用即可"""
    while queue:
        cur = queue.pop()
        print(cur, end=' ')
        for next in DAG[cur]:
            queue.append(next)
    print()


if __name__ == '__main__':
    next(solve1())
    solve2()

