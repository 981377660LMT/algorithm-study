from math import ceil
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 给定一个队列，初始时队列为空。

# 首先，依次从队尾插入 a、b、c、d、e 五个元素。

# 随后，不断重复以下操作：

# 设当前队头元素为 x。
# 依次从队尾插入两个 x。
# 将队头元素 x 弹出队列。
# 例如，第 1 轮操作过后，队列变为 b、c、d、e、a、a，第 1 个被弹出队列的元素为 a；第 2 轮操作过后，队列变为 c、d、e、a、a、b、b，第 2 个被弹出队列的元素为 b；第 3 轮操作过后，队列变为 d、e、a、a、b、b、c、c，第 3 个被弹出队列的元素为 c......

# 请你计算并输出第 n 个被弹出队列的元素。

if __name__ == "__main__":
    n = int(input()) - 1
    count = 1
    while n >= count * 5:
        n -= count * 5
        count *= 2
    div = n // count
    print(["a", "b", "c", "d", "e"][div])
