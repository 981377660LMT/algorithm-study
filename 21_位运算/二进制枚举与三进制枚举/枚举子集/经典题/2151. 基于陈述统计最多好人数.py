from typing import List


# 2 <= n <= 15
# 好人不能说假话(说好人是坏人，说坏人是好人)
# 验证好人之间的评价是否自洽
# !0:坏人 1:好人 2:不知道
class Solution:
    def maximumGood(self, statements: List[List[int]]) -> int:
        def check(goodMask: int) -> bool:
            for i in range(n):
                for j in range(n):
                    if i == j or not goodMask & (1 << i):
                        continue
                    if statements[i][j] == 0 and goodMask & (1 << j):
                        return False  # 说好人是坏人
                    if statements[i][j] == 1 and not goodMask & (1 << j):
                        if goodMask & (1 << i):
                            return False  # 说坏人是好人
            return True

        n = len(statements)
        res = 0
        for state in range(1, 1 << n):
            if check(state):
                res = max(res, state.bit_count())
        return res


# 2 1 2 1
print(Solution().maximumGood(statements=[[2, 1, 2], [1, 2, 2], [2, 0, 2]]))
print(Solution().maximumGood(statements=[[2, 0], [0, 2]]))
print(
    Solution().maximumGood(
        statements=[
            [2, 0, 2, 2, 0],
            [2, 2, 2, 1, 2],
            [2, 2, 2, 1, 2],
            [1, 2, 0, 2, 2],
            [1, 0, 2, 1, 2],
        ]
    )
)
print(Solution().maximumGood(statements=[[2, 2, 2, 2], [1, 2, 1, 0], [0, 2, 2, 2], [0, 0, 0, 2]]))


# 变形
# https://codeforces.com/problemset/problem/156/B

# 有 n(<=1e5) 个人，编号从 1 到 n。其中恰好有一个人是罪犯。
# 同时还有 n 条陈述，每条陈述要么是 +x，表示 x 是罪犯；要么是 -x，表示 x 不是罪犯。(1<=x<=n)
# 已知这 n 条陈述中恰好有 m(<=n) 条是实话，有 n-m 条是假话。
# 对于每条陈述，如果这条陈述一定是实话，输出 "Truth"；如果一定是假话，输出 "Lie"；如果不确定是真是假，输出 "Not defined"。


# https://codeforces.com/problemset/submission/156/159122474

# 枚举。
# 假设 i 是罪犯，如果 “i 是罪犯” 的陈述数和 “j 不是罪犯”(j≠i) 的陈述数之和等于 m，则 i 可能是罪犯，否则 i 一定不是罪犯。
# 然后遍历所有陈述，按照上面统计的结果来输出对应的字符串。
