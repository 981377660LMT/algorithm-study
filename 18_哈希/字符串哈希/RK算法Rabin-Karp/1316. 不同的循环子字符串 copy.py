# O(n^3) 比较字符串切片最坏时间复杂度也为O(n)
class Solution:
    def distinctEchoSubstrings(self, text: str) -> int:
        n = len(text)
        res = 0
        visited = set()
        for i in range(n):
            for j in range(i + 1, n):
                # 这里可替换为hasher 但是python切片比hasher快
                slice = text[i:j]
                slice_len = j - i
                if j + slice_len <= n and slice == text[j : j + slice_len] and slice not in visited:
                    visited.add(slice)
                    res += 1

        return res
