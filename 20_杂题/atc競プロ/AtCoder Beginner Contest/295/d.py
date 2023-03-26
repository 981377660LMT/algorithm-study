from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 20230322 は並べ替えると 02320232 となり、これは 0232 を
# 2 度繰り返しています。
# このように、数字のみからなる文字列であって、適切に文字を並び替える (そのままでもよい) ことによって同じ列を
# 2 度繰り返すようにできるものを 嬉しい列 と呼びます。
# 数字のみからなる文字列
# S が与えられるので、以下の条件を全て満たす整数の組
# (l,r) はいくつあるか求めてください。

# 1≤l≤r≤∣S∣ (
# ∣S∣ は
# S の長さ)
# S の
# l 文字目から
# r 文字目までの (連続する) 部分文字列は嬉しい列である。
if __name__ == "__main__":
    s = input()
    nums = [int(i) for i in s]
    # !多少个子数组,每种元素的个数都是偶数
    # 异或+状态压缩
    preSum = defaultdict(int, {0: 1})  # 如果记录索引就是{0: -1}
    res, curSum = 0, 0
    for i, num in enumerate(nums):
        curSum ^= 1 << num
        res += preSum[curSum]
        preSum[curSum] += 1
    print(res)
