fibo = [1, 1]
for i in range(100):
    fibo.append(fibo[-1] + fibo[-2])


class Solution:
    def stick(self, a: int):
        # write code here
        res = 0
        curSum = 0
        for num in fibo:
            curSum += num
            if curSum > a:
                break
            res += 1
        return res
