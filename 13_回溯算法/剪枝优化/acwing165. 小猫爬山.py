# 每租用一辆缆车，翰翰和达达就要付 1 美元，所以他们想知道，最少需要付多少美元才能把这 N 只小猫都运送下山？
# 每辆缆车上的小猫的重量之和不能超过 W。
# N 只小猫下山

# 搜索：
# 把某个猫加到最后一组
# 新开一个组


# 1. 优先安排体重大的猫, 因此, 对猫从大到小排序.
# 2. 如果已有的缆车数量大于等于当前找到的最优答案, 那么可以回溯了.
def dfs(index: int) -> None:
    global res
    if len(groups) >= res:
        return
    if index == n:
        res = min(res, len(groups))
        return

    # 放到一组
    for i in range(len(groups)):
        if groups[i] + nums[index] <= limit:
            groups[i] += nums[index]
            dfs(index + 1)
            groups[i] -= nums[index]

    # 新开一组
    groups.append(nums[index])
    dfs(index + 1)
    groups.pop()


n, limit = map(int, input().split())
nums = []
for _ in range(n):
    nums.append(int(input()))

nums.sort(reverse=True)
groups = []
res = n
dfs(0)
print(res)
