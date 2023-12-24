from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 高橋君は
# N 組の靴下を持っており、
# i 番目の組は色
# i の靴下
# 2 枚からなります。 ある日タンスの中を整理した高橋君は、色
# A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# K
# ​
#   の靴下を
# 1 枚ずつなくしてしまったことに気づいたので、残っている
# 2N−K 枚の靴下を使って、靴下
# 2 枚ずつからなる
# ⌊
# 2
# 2N−K
# ​
#  ⌋ 個の組を新たに作り直すことにしました。 色
# i の靴下と色
# j の靴下からなる組の奇妙さは
# ∣i−j∣ として定義され、高橋君は奇妙さの総和をできるだけ小さくしたいです。

# 残っている靴下をうまく組み合わせて
# ⌊
# 2
# 2N−K
# ​
#  ⌋ 個の組を作ったとき、奇妙さの総和が最小でいくつになるか求めてください。 なお、
# 2N−K が奇数のとき、どの組にも含まれない靴下が
# 1 枚存在することに注意してください。
if __name__ == "__main__":
    N, K = map(int, input().split())
    A = list(map(int, input().split()))
    counter = [2] * (N + 1)
    counter[0] = 0
    for a in A:
        counter[a] -= 1

    nums = []
    for i in range(N + 1):
        nums.extend([i] * counter[i])

    if len(nums) % 2 == 1:
        # 跳过一个数不选
        @lru_cache(None)
        def dfs(index: int, jumped: bool) -> int:
            if index == len(nums):
                return 0
            if index == len(nums) - 1:
                return 0 if not jumped else INF
            res = dfs(index + 2, jumped) + abs(nums[index] - nums[index + 1])
            if not jumped:
                res = min(res, dfs(index + 1, True))
            return res

        res = min(dfs(0, False), dfs(1, True))
        print(res)
    else:
        print(sum(abs(nums[i] - nums[i + 1]) for i in range(0, len(nums), 2)))
