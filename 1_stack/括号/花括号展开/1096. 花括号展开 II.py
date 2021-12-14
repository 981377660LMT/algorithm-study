from typing import List

# https://leetcode-cn.com/problems/brace-expansion-ii/solution/python3-18xing-die-dai-fa-by-yuan-zhi-b/
class Solution:
    def braceExpansionII(self, expression: str) -> List[str]:
        stack = []
        # 第一个list代表已经计算好的组合部分，第二个list代表增长中的未知(单词)部分
        res, cur = [], []

        for char in expression:
            print(res, cur, stack)
            # 看到字母就“乘入”第二个list，例如[a]*b变成“ab”，[a,c]*b就变成['ab','cb']
            if char.isalpha():
                cur = [prefix + char for prefix in cur or ['']]
            # 看到“{”就把两个list推入stack暂存
            elif char == '{':
                stack.append(res)
                stack.append(cur)
                res, cur = [], []
            # 看到“}”就把两个list从stack pop出来并把当前结果“乘入”第二个list。
            elif char == '}':
                pre = stack.pop()
                preRes = stack.pop()
                cur = [preChar + curChar for curChar in res + cur for preChar in pre or ['']]
                res = preRes
            # 第二个list已经不可能再继续增长了，把第二个list并入第一个list并清空
            elif char == ',':
                res += cur
                cur = []

        return sorted(set(res + cur))


# print(Solution().braceExpansionII(expression="{a,b}{c,{d,e}}"))
# 输出：["ac","ad","ae","bc","bd","be"]
print(Solution().braceExpansionII(expression=r"{{a,zk},a{bfg,c},{ab,z}}"))
# 输出：["a","ab","ac","z"]
# 解释：输出中 不应 出现重复的组合结果。
