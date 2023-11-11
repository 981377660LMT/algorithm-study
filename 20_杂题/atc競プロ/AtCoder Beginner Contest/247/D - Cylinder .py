# 有一个队列,有两种操作
# 1.往队尾插入c个分数为a的小球
# 2.取出队首的c个球,并求出这些球的分数之和(保证球数足够c个)
# 分析
# 用队列模拟就可以了,注意队列模拟不要一个个往里放球,
# 而是一次插入操作只插入其分数和个数取球时如果队首球不够就 pop继续遍历队首,
# !球够就修改队首,因为进队和出队操作和修改队首操作至多执行n次,所以复杂度是O(n)的.
from collections import deque
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    q = int(input())
    queue = deque()
    for _ in range(q):
        qt, *rest = map(int, input().split())
        if qt == 1:
            score, count = rest
            queue.append((score, count))
        elif qt == 2:
            remain = rest[0]
            res = 0
            while remain:
                score, count = queue.popleft()
                take = min(remain, count)
                res += take * score
                remain -= take
                count -= take
                if count:
                    queue.appendleft((score, count))
            print(res)
        else:
            raise ValueError("Invalid input")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
