from typing import List


class Solution:
    def addOperators(self, num: str, target: int) -> List[str]:
        #         (初始化字符串， 目标值， 表达式， 上一次的数字， 目前表达式的运算结果)
        def helper(num, target, exp, pre, sum):
            if not num:
                # 遍历完了num，并且找到了结果为target的表达式
                if sum == target:
                    ans.append(exp)
                return

            for i in range(1, len(num) + 1):
                # 含前导0的数字无效
                if i > 1 and num[0] == '0':
                    continue

                cur = int(num[:i])
                # 第一次递归的时候，表达式exp为空
                if exp == '':
                    helper(num[i:], target, num[:i], cur, cur)
                else:
                    helper(num[i:], target, exp + '+' + num[:i], cur, sum + cur)
                    helper(num[i:], target, exp + '-' + num[:i], -cur, sum - cur)
                    # 像 3 + 4 * 5 这种，pre=4, SUM=3+4, 所以碰到*的时候，需要先减掉pre的值4，即：(3+4)-4+4*5
                    helper(num[i:], target, exp + '*' + num[:i], pre * cur, sum - pre + pre * cur)

        ans = []
        helper(num, target, '', 0, 0)
        return ans


#                 if i > 1 and num[0] == '0':
#                     continue
# 可以剪枝成

#                 if i > 1 and num[0] == '0':
#                     break
