from collections import defaultdict


# n 个小朋友排成一排，从左到右依次编号为 1∼n。
# 第 i 个小朋友的身高为 hi。
# 虽然队伍已经排好，但是小朋友们对此并不完全满意。
# 对于一个小朋友来说，如果存在其他小朋友身高比他更矮，却站在他右侧的情况，该小朋友就会感到不满。
# 每个小朋友的不满程度都可以量化计算，具体来说，对于第 i 个小朋友：
# 如果存在比他更矮且在他右侧的小朋友，那么他的不满值等于其中最靠右的那个小朋友与他之间的小朋友数量。
# 如果不存在比他更矮且在他右侧的小朋友，那么他的不满值为 −1。
# 请你计算并输出每个小朋友的不满值。
# 注意，第 1 个小朋友和第 2 个小朋友之间的小朋友数量为 0，第 1 个小朋友和第 4 个小朋友之间的小朋友数量为 2。

# !寻找每个数右侧最后一个比他小的数(排序遍历)

if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))

    groups = defaultdict(list)
    for i, num in enumerate(nums):
        groups[num].append(i)

    curMax = -1
    res = [-1] * n
    for key in sorted(groups):
        group = groups[key]
        for i in group:
            if curMax > i:
                res[i] = curMax - i - 1  # !和他之间的小朋友数量
        for i in group:
            curMax = max(curMax, i)

    print(*res)
