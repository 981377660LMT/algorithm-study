# 工作分配,任务安排,最大利润
# n,m<=1e5
class Solution:
    def solve(self, people, jobs, profits):
        # 1e5说明肯定要对工人和工作都排序
        pairs = [(p, j) for p, j in zip(profits, jobs)]
        pairs.sort(reverse=True)
        people.sort(reverse=True)
        res = 0

        j = 0
        for i in range(len(people)):
            while j < len(pairs):
                if people[i] >= pairs[j][1]:
                    res += pairs[j][0]
                    break
                j += 1
        return res


print(Solution().solve(people=[5, 7, 8], jobs=[6, 5, 8], profits=[1, 2, 3]))
