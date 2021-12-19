from typing import List


class Solution:
    def findRLEArray(self, encoded1: List[List[int]], encoded2: List[List[int]]) -> List[List[int]]:
        res = []
        i, j = 0, 0
        while i < len(encoded1) and j < len(encoded2):
            v1, f1 = encoded1[i]
            v2, f2 = encoded2[j]

            v = v1 * v2
            f = min(f1, f2)

            # 添加
            if not res or res[-1][0] != v:
                res.append([v, f])
            else:
                res[-1] = [v, res[-1][1] + f]

            # 移动 相同逻辑
            if f1 - f > 0:
                encoded1[i] = [v1, f1 - f]
            else:
                i += 1

            if f2 - f > 0:
                encoded2[j] = [v2, f2 - f]
            else:
                j += 1

        return res


print(Solution().findRLEArray(encoded1=[[1, 3], [2, 1], [3, 2]], encoded2=[[2, 3], [3, 3]]))
# 输出: [[2,3],[6,1],[9,2]]
# 解释: encoded1 扩展为 [1,1,1,2,3,3] ，encoded2 扩展为 [2,2,2,3,3,3]。
# prodNums = [2,2,2,6,9,9]，压缩成行程编码数组 [[2,3],[6,1],[9,2]]。
