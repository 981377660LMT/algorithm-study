import re
from typing import List
from functools import lru_cache

# 测试数据保证正确表达式结果在范围 [0, 1000] 以内。
# 3 <= s.length <= 31

# 枚举分割点


class Solution:
    def scoreOfStudents(self, s: str, answers: List[int]) -> int:
        @lru_cache(None)
        def dfs(left: int, right: int) -> dict:
            if left + 1 >= right:
                return {int(arr[left]): 0}

            res = {}
            for mid in range(left + 1, right, 2):
                for leftRes in dfs(left, mid):
                    for rightRes in dfs(mid + 1, right):
                        cur = leftRes * rightRes if arr[mid] == '*' else leftRes + rightRes
                        # 注意这里剪枝
                        if cur <= 1000:
                            res[cur] = 2
            return res

        arr = list(re.split(r'(\D)', s))
        res = {**dfs(0, len(arr) - 1), eval(s): 5}
        return sum(res.get(answer, 0) for answer in answers)


print(Solution().scoreOfStudents(s="7+3*1*2", answers=[20, 13, 42]))
# 输出：7
# 解释：如上图所示，正确答案为 13 ，因此有一位学生得分为 5 分：[20,13,42] 。
# 一位学生可能通过错误的运算顺序得到结果 20 ：7+3=10，10*1=10，10*2=20 。所以这位学生得分为 2 分：[20,13,42] 。
# 所有学生得分分别为：[2,5,0] 。所有得分之和为 2+5+0=7 。

