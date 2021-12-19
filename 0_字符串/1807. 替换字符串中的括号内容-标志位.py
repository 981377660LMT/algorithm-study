from typing import List


class Solution:
    def evaluate(self, s: str, knowledge: List[List[str]]) -> str:
        d = {k: v for k, v in knowledge}
        res = []
        cur = ''
        isReplaceMode = False
        
        for char in s:
            if char == '(':
                isReplaceMode = True
            elif char == ')':
                isReplaceMode = False
                res.append(d.get(cur, '?'))
                cur = ''
            elif isReplaceMode:
                cur += char
            else:
                res.append(char)

        return ''.join(res)


print(Solution().evaluate(s="(name)is(age)yearsold", knowledge=[["name", "bob"], ["age", "two"]]))
# 输出："bobistwoyearsold"
# 解释：
# 键 "name" 对应的值为 "bob" ，所以将 "(name)" 替换为 "bob" 。
# 键 "age" 对应的值为 "two" ，所以将 "(age)" 替换为 "two" 。
