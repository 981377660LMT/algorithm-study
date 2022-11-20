from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 長さ N の数列 A=(A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# N
# ​
#  ) が与えられます。

# Q 個のクエリが与えられるので、順番にすべて処理してください。 q 番目 (1≤q≤Q) のクエリは以下の 3 つのいずれかの形式で、それぞれ次のようなクエリを表します。

# 1 x
# q
# ​
#   ： A のすべての要素に x
# q
# ​
#   を代入する。
# 2 i
# q
# ​
#   x
# q
# ​
#   ： A
# i
# q
# ​

# ​
#   に x
# q
# ​
#   を加える。
# 3 i
# q
# ​
#   ： A
# i
# q
# ​

# ​
#   の値を出力する。


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    q = int(input())

    add = defaultdict(int)
    curRange = -1

    for _ in range(q):
        args = list(map(int, input().split()))
        if args[0] == 1:
            add.clear()  # !能不能开[0]*n的数组
            curRange = args[1]
        elif args[0] == 2:
            i, x = args[1:]
            add[i] += x
        else:
            i = args[1]
            if curRange == -1:
                print(nums[i - 1] + add[i])
            else:
                print(curRange + add[i])
