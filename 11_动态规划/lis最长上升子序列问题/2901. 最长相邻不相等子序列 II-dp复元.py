# 2901. 最长相邻不相等子序列 II-dp复元/dp还原/dp求方案


from typing import List


class Solution:
    def getWordsInLongestSubsequence(
        self, n: int, words: List[str], groups: List[int]
    ) -> List[str]:
        def check(i: int, j: int) -> int:
            if groups[i] == groups[j]:
                return False
            s1, s2 = words[i], words[j]
            if len(s1) != len(s2):
                return False
            return sum(a != b for a, b in zip(s1, s2)) == 1

        dp = [1] * n  # dp[i] 表示以 nums[i] 结尾的最长上升子序列的长度
        pre = [-1] * n  # pre[i] 表示以 nums[i] 结尾的最长上升子序列的前一个元素的下标
        for i in range(1, n):
            for j in range(i):
                if check(i, j):
                    cand = dp[j] + 1
                    if cand > dp[i]:
                        dp[i] = cand
                        pre[i] = j

        max_, max_i = -1, -1
        for i, v in enumerate(dp):
            if v > max_:
                max_ = v
                max_i = i

        cur = max_i
        path = []
        while cur != -1:
            path.append(words[cur])
            cur = pre[cur]
        return path[::-1]


# n = 3, words = ["bab","dab","cab"], groups = [1,2,2]

if __name__ == "__main__":
    print(Solution().getWordsInLongestSubsequence(3, ["bab", "dab", "cab"], [1, 2, 2]))

    print(
        Solution().getWordsInLongestSubsequence(
            9,
            ["cdb", "cdd", "cd", "dcc", "cca", "cda", "ca", "cc", "bcc"],
            [8, 5, 9, 5, 2, 7, 4, 7, 3],
        )
    )
