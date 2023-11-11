# 求隐藏山脉数组的最大值
# 交互式问题

# n<=1500
# 0<=ai<=1e9
# T组测试数据 T<=50

# !凸函数最值:三分法 => 黄金分割法(加速)

import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline

fib = [1, 1]
while fib[-1] < 1500:
    fib.append(fib[-1] + fib[-2])


def query(index: int, n: int) -> int:
    """1 <= index <= n"""
    if index <= 0 or index > n:
        return 0
    print(f'? {index}', flush=True)
    return int(input())


def output(val: int) -> None:
    # !注意要flush=True 阻止函数对输出数据进行缓冲、强行刷新
    # !如果将 flush 参数设置为 True ，则 print() 函数将不会对数据进行缓冲以提高效率，而是在每次调用时不断地对数据进行刷新
    print(f'! {val}', flush=True)


T = int(input())
for _ in range(T):
    n = int(input())
    res = 0
    cand1, cand2 = -1, -1
    left, right = 0, fib[-1]
    for i in range(3, len(fib))[::-1]:
        if cand1 == -1:
            cand1 = query(left + fib[i - 2], n)
        if cand2 == -1:
            cand2 = query(right - fib[i - 2], n)
        if cand1 < cand2:
            left += fib[i - 2]
            cand1 = cand2
            cand2 = -1
        else:
            right -= fib[i - 2]
            cand2 = cand1
            cand1 = -1
    output(max(cand1, cand2))

