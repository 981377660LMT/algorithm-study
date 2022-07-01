from operator import add, mul, sub, truediv
import re
from typing import Dict, List
from functools import lru_cache

# 测试数据保证正确表达式结果在范围 [0, 1000] 以内。
# 3 <= s.length <= 31

# 枚举分割点

# 如果一位学生的答案 等于 表达式的正确结果，这位学生将得到 5 分。
# 否则，如果答案由 一处或多处错误的运算顺序 计算得到，那么这位学生能得到 2 分。
# 否则，这位学生将得到 0 分。

OPT = {"+": add, "-": sub, "*": mul, "/": truediv}


class Solution:
    def scoreOfStudents(self, s: str, answers: List[int]) -> int:
        @lru_cache(None)
        def dfs(left: int, right: int) -> Dict[int, int]:
            if left > right:
                return {}
            if left == right:
                return {int(arr[left]): 0}

            res: Dict[int, int] = {}
            for i in range(left + 1, right, 2):
                opt = OPT[arr[i]]
                for leftRes in dfs(left, i - 1):
                    for rightRes in dfs(i + 1, right):
                        cur = opt(leftRes, rightRes)
                        # !注意这里剪枝
                        if cur <= 1000:
                            res[cur] = 2
            return res

        arr = re.split(r"(\D)", s)
        res = {**dfs(0, len(arr) - 1), eval(s): 5}  # 字典解构
        return sum(res.get(answer, 0) for answer in answers)


print(Solution().scoreOfStudents(s="7+3*1*2", answers=[20, 13, 42]))
# 输出：7
# 解释：如上图所示，正确答案为 13 ，因此有一位学生得分为 5 分：[20,13,42] 。
# 一位学生可能通过错误的运算顺序得到结果 20 ：7+3=10，10*1=10，10*2=20 。所以这位学生得分为 2 分：[20,13,42] 。
# 所有学生得分分别为：[2,5,0] 。所有得分之和为 2+5+0=7 。
