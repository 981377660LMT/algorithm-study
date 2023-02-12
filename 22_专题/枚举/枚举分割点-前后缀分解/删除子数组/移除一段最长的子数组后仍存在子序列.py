class Solution:
    def solve(self, s: str, t: str) -> int:
        """移除一段子数组后,s中仍然存在子序列t;求可移除的数组的最大长度"""
        # 子数组：前后缀/滑窗
        # 找前后缀+双指针找最长
        if not t:
            return len(s)

        n1, n2 = len(s), len(t)
        pre, suffix = [0] * n1, [0] * n1

        hit = 0
        for i in range(n1):
            if hit < n2 and s[i] == t[hit]:
                hit += 1
            pre[i] = hit

        hit = 0
        for i in range(n1 - 1, -1, -1):
            if hit < n2 and s[i] == t[n2 - 1 - hit]:
                hit += 1
            suffix[i] = hit

        # 滑窗里的部分可以删去
        l, r = 0, n1 - 1
        while suffix[r] < n2:
            r -= 1

        res = r
        while l < n1:
            while r < n1 and pre[l] + suffix[r] >= n2:
                r += 1
            res = max(res, r - l - 2)
            l += 1
        return res


print(Solution().solve(s="abcabac", t="bc"))
