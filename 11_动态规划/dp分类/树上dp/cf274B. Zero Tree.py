# https://blog.csdn.net/guhaiteng/article/details/52738661?spm=1001.2101.3001.6650.1&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1-52738661-blog-48175347.pc_relevant_multi_platform_featuressortv2dupreplace&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1-52738661-blog-48175347.pc_relevant_multi_platform_featuressortv2dupreplace&utm_relevant_index=2
# 【codeforces 274B】【树形DP】
# B. Zero Tree
# 【一棵树，每个点有权值，每次操作可以对一个联通子集中的点全部加或者减1，
# 且每次操作必须包含根节点，问最少多少次操作权值全为0】

# !所有非叶节点需要维护一下本结点最小需要执行多少次+和-操作，
# 父节点根据子节点的这些信息计算出自己的操作次数，一直返回到根节点
# 所以对于每个点u,它被加修改和减修改的次数，
# 就是它的儿子中进行该操作的最大次数，
# !因为如果有两个儿子都需要进行该操作，
# !那么完全可以两步并一步，所以只需要取最大值就可以了。


n = int(input())
parents = list(map(int, input().split()))
values = list(map(int, input().split()))

adjList = [[] for _ in range(n)]
for cur, pre in enumerate(parents):
    if pre == -1:
        continue
    adjList[pre].append(cur)
    adjList[cur].append(pre)

up, down = [0] * n, [0] * n  # 自己需要加的次数，需要减的次数


def dfs(cur: int, pre: int) -> None:
    for next in adjList[cur]:
        if next == pre:
            continue
        dfs(next, cur)
        up[cur] = max(up[cur], up[next])
        down[cur] = max(down[cur], down[next])

    values[cur] += up[cur] - down[cur]
    if values[cur] > 0:
        down[cur] += values[cur]
    else:
        up[cur] -= values[cur]


dfs(0, -1)
print(up[0] + down[0])
