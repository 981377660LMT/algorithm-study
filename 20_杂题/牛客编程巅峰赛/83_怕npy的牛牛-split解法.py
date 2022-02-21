# 字符串中不同时含有n,p,y三个字母的最长字串的长度是多少？
# 思路1：滑窗
# 思路二：split
class Solution:
    def Maximumlength(self, x: str) -> int:
        """不同时包含npy三个字母的最长字串的长度是多少？"""
        # # write code here
        # n = x.split('n')
        # p = x.split('p')
        # y = x.split('y')
        # res = 0
        # for i in n:
        #     res = max(len(i), res)
        # for i in y:
        #     res = max(len(i), res)
        # for i in p:
        #     res = max(len(i), res)
        # return res
        left, right = 0, 0
        res = 0
        counter = dict(n=0, p=0, y=0)
        while right < len(x):
            if x[right] in "npy":
                counter[x[right]] += 1
            while counter['n'] and counter['p'] and counter['y']:
                if x[left] in "npy" and counter[x[left]]:
                    counter[x[left]] -= 1
                left += 1
            res = max(res, right - left + 1)
            right += 1
        return res
