# https://leetcode.cn/problems/rearrange-string-k-distance-apart/solutions/854882/fu-za-du-wei-onde-zhi-jie-gou-zao-fa-by-jlhsx/

from collections import Counter


class Solution:
    def rearrangeString(self, s: str, k: int) -> str:
        if k <= 1:
            return s
        c = Counter(s)
        chars = sorted(c.items(), key=lambda x: -x[1])
        n = len(s)
        res = [""] * n
        i = 0
        j = k - 1
        curr = 0
        less_size = n // k
        more_size = less_size + 1
        more_track = n % k
        less_track = k - more_track
        for c, v in chars:
            if v > more_size:
                return ""
            elif v == more_size:
                if not more_track:
                    return ""
                # Use full track
                res[i::k] = [c] * v
                i += 1
                more_track -= 1
                curr = i
            elif v == less_size and less_track:
                # Use full track
                res[j::k] = [c] * v
                j -= 1
                less_track -= 1
            else:
                # Fill in order
                for _ in range(v):
                    res[curr] = c
                    curr += k
                    if curr >= n:
                        i += 1
                        curr = i
        return "".join(res)
