from collections import defaultdict, deque


class Solution:
    def canReach(self, s: str, minJump: int, maxJump: int) -> bool:
        """差分数组区间更新，边遍历边还原数组值
        
        直接更新diff字典的写法
        """
        if s[-1] == '1':
            return False

        diff = defaultdict(int)
        for i, char in enumerate(s):
            diff[i] += diff[i - 1]
            if char == '0' and (i == 0 or diff[i] > 0):
                diff[i + minJump] += 1
                diff[i + maxJump + 1] -= 1
        return diff[len(s) - 1] > 0

    def canReach2(self, s: str, minJump: int, maxJump: int) -> bool:
        """差分数组区间更新，边遍历边还原数组值
        
        不修改diff 用 curSum 的写法
        """
        if s[-1] == '1':
            return False

        n = len(s)
        diff, curSum = [0] * (n + 1), 0
        for i, char in enumerate(s):
            curSum += diff[i]
            if char == '0':
                if i == 0 or curSum > 0:
                    diff[min(i + minJump, n)] += 1
                    diff[min(i + maxJump + 1, n)] -= 1
        return curSum > 0
