from typing import List

# 1 <= S.length <= 200
class Solution:
    def splitIntoFibonacci(self, S: str) -> List[int]:
        def backtrack(index: int, path: List[int]):
            if index == n:  # 退出条件
                if len(path) >= 3:
                    self.res = path
                    return

            for i in range(index, n):
                if S[index] == "0" and i > index:  # 当两位数数字以0开头时,死路
                    continue

                if int(S[index : i + 1]) > 2 ** 31 - 1 or int(S[index : i + 1]) < 0:  # 剪枝
                    continue

                if len(path) < 2:
                    backtrack(i + 1, path + [int(S[index : i + 1])])
                elif int(S[index : i + 1]) == path[-1] + path[-2]:
                    backtrack(i + 1, path + [int(S[index : i + 1])])

        n = len(S)
        self.res = []
        backtrack(0, [])
        return self.res


print(Solution().splitIntoFibonacci("11235813"))

