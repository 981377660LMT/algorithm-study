# C - Snake Queue
# https://atcoder.jp/contests/abc389/tasks/abc389_c
# 题意：三种操作:
#
# 1. 每次再结尾插入一个长度为m的蛇。
# 2. 删除第一条蛇。
# 3. 问当前第k条蛇蛇头的位置。
#
# 用数组累加前缀和，模拟即可。


from typing import List, Tuple


def snakeQueue(operations: List[Tuple[int, int]]) -> List[int]:
    presum = []
    curSum = 0
    ptr = 0
    res = []
    for op in operations:
        if op[0] == 1:
            m = op[1]
            presum.append(curSum)
            curSum += m
        elif op[0] == 2:
            ptr += 1
        else:
            k = op[1]
            res.append(presum[ptr + k] - presum[ptr])
    return res


if __name__ == "__main__":
    q = int(input())
    operations = []
    for _ in range(q):
        op, *args = list(map(int, input().split()))
        if op == 1:
            m = args[0]
            operations.append((op, m))
        elif op == 2:
            operations.append((op,))
        else:
            k = args[0]
            k -= 1
            operations.append((op, k))
    res = snakeQueue(operations)
    print("\n".join(map(str, res)))
