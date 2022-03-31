from typing import List, Tuple


class Solution:
    def minWindow(self, s1: str, s2: str) -> str:
        """
        找出 s1 中最短的（连续）子串 W ，使得 s2 是 W 的 子序列 。
        如果有不止一个最短长度的窗口，返回开始位置最靠左的那个。
        """
        n = len(s1)
        nexts: List[Tuple[int, ...]] = [()] * n
        last = [-1] * 26
        for i in range(n - 1, -1, -1):
            nexts[i] = tuple(last)
            last[ord(s1[i]) - 97] = i

        # 假设窗口的起点为 S[i]，S[i] = T[0]。
        # 那么要拓展窗口就需要在 S[i+1:] 中找到最近的 S[j]，
        # 使得 S[j] = T[1]。同样的道理，再从 S[j+1:] 中找到最近的 S[k]，
        # 使得 S[k] = T[2]。按照这种方式，就可以找到包含整个 T 的窗口。
        res = None
        starts = [i for i, char in enumerate(s1) if char == s2[0]]
        for start in starts:
            cur = start
            for char in s2[1:]:
                cur = nexts[cur][ord(char) - 97]
                if cur == -1:
                    break
            else:
                if res is None or cur - start + 1 < res[1] - res[0] + 1:
                    res = (start, cur)

        return s1[res[0] : res[1] + 1] if res is not None else ''


print(Solution().minWindow("abcdebdde", "bde"))
