from collections import Counter


class Solution:
    def equalDigitFrequency(self, s: str) -> int:
        # O(10*n^2)
        n = len(s)
        visited = set()

        for i in range(n):
            counter = [0] * 10
            for j in range(i, n):
                counter[ord(s[j]) - ord('0')] += 1
                if len(set(counter) - {0}) == 1:
                    visited.add(s[i : j + 1])

        return len(visited)

    def equalDigitFrequency2(self, s: str) -> int:
        """字符串哈希 O(n^2)
        
        遍历时 动态维护 各个数字出现次数的最大值的数量，如果数量 = 数字出现的个数，那么就满足条件
        """
        n = len(s)
        visited = set()

        for start in range(n):
            counter = Counter()
            maxCount, max_ = 0, 0  # 取得最大值的字符数 最大值

            for end in range(start, n):
                counter[s[end]] += 1

                if counter[s[end]] > max_:
                    maxCount = 1
                    max_ = counter[s[end]]
                elif counter[s[end]] == max_:
                    maxCount += 1

                if maxCount == len(counter):
                    visited.add(s[start : end + 1])  # 别的语言可以字符串哈希优化

        return len(visited)


print(Solution().equalDigitFrequency(s="12321"))
