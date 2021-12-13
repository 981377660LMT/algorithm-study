from collections import defaultdict


# K4(ON(SO3)2)2
class Solution:
    def countOfAtoms(self, s: str) -> str:
        dic = defaultdict(int)
        cur_rate = [1]  # 记录当前倍率
        lower = ''
        count = ''

        for i in range(len(s) - 1, -1, -1):
            cur = s[i]
            if cur.isdigit():
                count = cur + count
            elif cur.islower():
                lower = cur + lower
            elif cur == ')':
                cur_rate.append(cur_rate[-1] * int(count or '1'))
                count = ''
            elif cur == '(':
                cur_rate.pop()
            elif cur.isupper():
                dic[cur + lower] += cur_rate[-1] * int(count or '1')
                lower = ''
                count = ''

        sb = []
        for key, value in sorted(dic.items()):
            sb.append(key + str(value) if value > 1 else key)
        return ''.join(sb)


print(Solution().countOfAtoms('K4(ON(SO3)2)2'))
# 输出："K4N2O14S4"
# 解释：原子的数量是 {'K': 4, 'N': 2, 'O': 14, 'S': 4}。
