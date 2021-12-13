# 将所有x替换成y 求最小值与最大值的最大差值
# 得到的新的整数 不能 有前导 0 ，得到的新整数也 不能 是 0 。

# 1 <= num <= 10^8

# 最大的数： 最左边位不为9的数变为9
# 最小的数：首位不为1，改首位为1；首位为1，把后面出现的不为0/1的数改为0
class Solution:
    def maxDiff(self, num: int) -> int:
        s = str(num)
        maxVal = minVal = str(num)

        for digit in s:
            if digit != '9':
                maxVal = maxVal.replace(digit, '9')
                break

        if minVal[0] != '1':
            minVal = minVal.replace(minVal[0], '1')
        else:
            for digit in minVal[1:]:
                if digit not in '01':
                    minVal = minVal.replace(digit, '0')
                    break

        return int(maxVal) - int(minVal)

    # 暴力
    def maxDiff2(self, num: int) -> int:
        s, res = str(num), []
        for i in range(0, 10):
            for j in range(0, 10):
                tmp = s.replace(str(i), str(j))
                if tmp[0] == '0' or int(tmp) == 0:
                    continue
                res.append(int(tmp))
        return max(res) - min(res)


print(Solution().maxDiff(num=555))
# 输出：888
# 解释：第一次选择 x = 5 且 y = 9 ，并把得到的新数字保存在 a 中。
# 第二次选择 x = 5 且 y = 1 ，并把得到的新数字保存在 b 中。
# 现在，我们有 a = 999 和 b = 111 ，最大差值为 888
