# m,n<=15
# 所有任务是否能在cores上执行
class Solution:
    def solve(self, cores, tasks):
        def bt(index: int) -> bool:
            if index == len(tasks):
                return True

            for i in range(len(cores)):
                if i and cores[i] == cores[i - 1]:
                    continue
                if cores[i] >= tasks[index]:
                    cores[i] -= tasks[index]
                    if bt(index + 1):
                        return True
                    cores[i] += tasks[index]

            return False

        tasks = sorted(tasks, reverse=True)
        return bt(0)


# tasks是否能在cores上运行
print(Solution().solve(cores=[8, 10], tasks=[2, 3, 3, 3, 7]))
# We can put tasks[0], tasks[1], and tasks[2] into cores[0] and the rest of the tasks into cores[1].
