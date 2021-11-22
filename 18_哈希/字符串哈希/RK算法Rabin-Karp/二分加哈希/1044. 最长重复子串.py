class Solution:
    def longestDupSubstring(self, s: str) -> str:
        n = len(s)

        # ------------------------------------------- rabin-karp编码（滚动哈细）
        BASE = 13331  # 经验值
        MOD = 10 ** 11 + 7  # 经验值

        # ---- 预计算BASE的幂
        mul = [0 for _ in range(n + 1)]
        mul[0] = 1
        for i in range(1, n + 1):
            mul[i] = mul[i - 1] * BASE % MOD

        # ---- 前缀和
        presum = [0 for _ in range(n + 1)]
        for i in range(n):
            presum[i + 1] = presum[i] * BASE + ord(s[i]) - ord('a')
            presum[i + 1] %= MOD

        # ---- 获取区间s[ll:rr]的哈希值
        def get_zone_val(l: int, r: int) -> int:
            res = presum[r + 1] - presum[l] * mul[r + 1 - l] % MOD
            return (res + MOD) % MOD

        # ------------------------------------------ 二分查找
        # -------- check 函数，找是否存在长度为mid的重复子串
        def search_start_idx(mid: int) -> int:
            visited = set()
            for rr in range(mid - 1, n):
                ll = rr - mid + 1
                tt = get_zone_val(ll, rr)
                if tt in visited:
                    return ll
                else:
                    visited.add(tt)
            return -1

        lo = 0
        hi = n
        while lo < hi:
            mid = (lo + hi + 1) // 2
            if search_start_idx(mid) != -1:
                lo = mid
            else:
                hi = mid - 1

        if lo == 0:
            return ""
        idx = search_start_idx(lo)
        return s[idx : idx + lo]

