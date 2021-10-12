from collections import defaultdict


# K4(ON(SO3)2)2
class Solution:
    def countOfAtoms(self, s: str) -> str:
        dic = defaultdict(int)
        num_stack = [1]
        lower = ''
        count = ''

        for i in range(len(s) - 1, -1, -1):
            cur = s[i]
            if cur.isdigit():
                count = cur + count
            elif cur.islower():
                lower = cur + lower
            elif cur == ')':
                num_stack.append(num_stack[-1] * int(count or '1'))
                count = ''
            elif cur == '(':
                num_stack.pop()
            elif cur.isupper():
                dic[cur + lower] += num_stack[-1] * int(count or '1')
                lower = ''
                count = ''

        sb = []
        for key, value in sorted(dic.items()):
            sb.append(key + str(value) if value > 1 else key)
        return ''.join(sb)


print(Solution().countOfAtoms('K4(ON(SO3)2)2'))
