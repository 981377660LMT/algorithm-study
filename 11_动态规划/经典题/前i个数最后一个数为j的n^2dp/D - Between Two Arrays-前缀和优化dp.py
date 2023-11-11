# 求数组C的个数
# 1.C数组单调不减(広義単調増加)
# 2.ai<=ci<=bi

# !1<=n<=3000
# !0<=ai<=bi<=3000

# !dp[i][val] 前缀和优化dp


from functools import reduce
from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))

    dp = [0] * 3005
    for num in range(nums1[0], nums2[0] + 1):
        dp[num] = 1

    for i in range(1, n):
        ndp = [0] * 3005
        dpSum = [0] + list(accumulate(dp))
        for num in range(nums1[i], nums2[i] + 1):
            ndp[num] = (ndp[num] + dpSum[num + 1]) % MOD
        dp = ndp

    res = reduce(lambda x, y: (x + y) % MOD, dp, 0)
    print(res)
