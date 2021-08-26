# 给定一组城市和每对城市之间的火车票的价钱，找到每个城市只访问一次并返回起点的最小车费花销。
import functools

n = int(input())
graph = []
for i in range(n):
    graph.append(list(map(int, input().split())))
visited = 1 << (n - 1)


@functools.lru_cache(None)
def dp(fr, to):
    if to == 0:
        return graph[fr][0]
    ans = float('inf')
    for i in range(1, n):
        if ((to >> (i - 1)) & 1) == 1:
            to ^= 1 << (i - 1)  # mark visited
            ans = min(ans, graph[fr][i] + dp(i, to))
            to ^= 1 << (i - 1)  # unmark
    return ans


print(dp(0, visited - 1))

# 输入例子1:
# 4
# 0 2 6 5
# 2 0 4 4
# 6 4 0 2
# 5 4 2 0

# 输出例子1:
# 13

# 例子说明1:
# 共 4 个城市，城市 1 和城市 1 的车费为0，城市 1 和城市 2 之间的车费为 2，城市 1 和城市 3 之间的车费为 6，城市 1 和城市 4 之间的车费为 5，依次类推。假设任意两个城市之间均有单程票可购买，且票价在1000元以内，无需考虑极端情况。