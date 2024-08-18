# 100409. 找出最大的 N 位 K 回文数
# https://leetcode.cn/problems/find-the-largest-palindrome-divisible-by-k/solutions/2884548/tong-yong-zuo-fa-jian-tu-dfsshu-chu-ju-t-m3pu/# 找出最大的N位回文数，使得该回文数可以被 K 整除。
# 求最大的被k整除的n位回文数。
# 1 <= n <= 1e5
# 1 <= k <= 9
# !dp求方案.
# 由于我们只关心回文数模 k 的值是否为 0，可以把：
# 当前从右到左填到第 i 位。
# 已填入的数字模 k 的值为 j。
# 作为状态 (i,j)。
# 填入数字后，回文模k的值变成了
# !j2 = (j+d*(10^i+10^(n-i-1)))%k
# !注意特判n为奇数且i=mid-1的情况(mid=ceil(n/2))，此时模值变成了 (j+d*10^mid)%k。
# 因为要求最大，所以这里用dfs简化dp.


class Solution:
    def largestPalindrome(self, n: int, k: int) -> str:
        pow10 = [1] * n
        for i in range(1, n):
            pow10[i] = pow10[i - 1] * 10 % k

        mid = (n + 1) // 2
        dp = [-1] * (mid + 1) * k
        pre = [-1] * (mid + 1) * k
        preValue = [-1] * (mid + 1) * k

        def dfs(pos: int, mod: int) -> int:
            if pos == mid:
                return 1 if mod == 0 else 0
            hash = pos * k + mod
            if dp[hash] != -1:
                return dp[hash]
            for d in range(9, -1, -1):
                if n & 1 and pos == mid - 1:
                    nMod = (mod + d * pow10[pos]) % k
                else:
                    nMod = (mod + d * (pow10[pos] + pow10[n - pos - 1])) % k
                if dfs(pos + 1, nMod) == 1:
                    dp[hash] = 1
                    nHash = (pos + 1) * k + nMod
                    pre[nHash] = hash
                    preValue[nHash] = d
                    return 1
            dp[hash] = 0
            return 0

        dfs(0, 0)

        res = []
        cur = mid * k
        while pre[cur] != -1:
            res.append(preValue[cur])
            cur = pre[cur]
        res.reverse()
        if n & 1:
            res += res[:-1][::-1]
        else:
            res += res[::-1]
        return "".join(map(str, res))


#  n = 3, k = 5
print(Solution().largestPalindrome(3, 5))  # 666
#  n = 5, k = 6

print(Solution().largestPalindrome(5, 6))  # 66666
