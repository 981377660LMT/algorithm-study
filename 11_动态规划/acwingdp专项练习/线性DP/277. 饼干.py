# 圣诞老人共有 M 个饼干，准备全部分给 N 个孩子。
# 每个孩子有一个贪婪度，第 i 个孩子的贪婪度为 g[i]。
# 如果有 k 个孩子拿到的饼干数比第 i 个孩子多，那么第 i 个孩子会产生 nums[i]×k 的怨气。
# 1≤N≤30 ,
# N≤M≤5000,

# 输出格式
# 第一行一个整数表示最小怨气总和。
# 第二行 N 个空格隔开的整数表示每个孩子分到的饼干数，若有多种方案，输出任意一种均可。
n, m = map(int, input().split())
angry = list(map(int, input().split()))


# 将所有小朋友按g[i]从大到小排序
