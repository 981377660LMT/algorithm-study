"""しりとり词语接龙 (bfs+博弈dp)
https://blog.csdn.net/weixin_51216553/article/details/124272620

Takahashi先手(叫做小T),Aoki后手(叫做小A )
每个人轮流说一个单词,不能说重复的单词
# !这个单词的前三个字母必须等于上一个单词的后三个字母 大小写有区别
现在问当小T 说出的第一个单词为si时,能否获得胜利??
n<=2e5
len(word)<=8 且只含有英文大小写


Solution:
# !一般博弈都是dfs 这里存在环(平局) 需要改成bfs倒序自底向上dp (后退解析)
# !因此需要从叶子节点开始建立`反图`
建图,每个单词前三个字符作为节点u,后三个字符作为节点v,连边
# !一共有52^3个顶点,n条边

1. 如果一个点没有任何的出边，那这个点就是先手必输，即不存在串使的该点能转移到别的点
2. 如果一个点`可以`一步走到一个先手必输点，则他先手必赢 (dfs找到对手必输的分支)
3. 如果一个点一步`只能`走到先手必胜点，则他先手必输  (dfs找到的分支全是对手必胜)
4. 如果这个点既不是先手必胜点和不是先手必输点，则是平局点

对于入度为0的点,他们的状态就是先手必输,状态为0
如果当前点u是先手必输,那对于u的子节点v,他的状态就是先手必胜,即1
"""

from collections import defaultdict, deque
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def shiritori(words: List[str]) -> List[int]:
    """各个单词作为先手开局单词时,先手必胜/必败/平局"""
    rAdjMap = defaultdict(list)
    deg = defaultdict(int)
    allVertex = set()
    for word in words:
        u, v = word[:3], word[-3:]
        rAdjMap[v].append(u)
        deg[u] += 1
        allVertex.add(u)
        allVertex.add(v)

    state = defaultdict(lambda: -1)  # 每个顶点的状态:先手必输:0,先手必胜:1,平局:-1
    queue = deque()
    for v in allVertex:
        if deg[v] == 0:
            queue.append(v)
            state[v] = 0

    while queue:
        cur = queue.popleft()
        for next in rAdjMap[cur]:
            if state[next] != -1:
                continue
            deg[next] -= 1
            # !1.正图上`可以`一步走到一个先手必输点，则他先手必赢 (对应dfs找到一个对手输的分支)
            if state[cur] == 0:
                state[next] = 1
                queue.append(next)
            # !2.正图上没有边可以走,则为必败态 (对应dfs找到的分支全是对手必胜)
            elif deg[next] == 0:
                state[next] = 0
                queue.append(next)
            # !3.其余情况为平局

    return [state[w[-3:]] for w in words]


if __name__ == "__main__":
    n = int(input())
    words = [input() for _ in range(n)]
    res = shiritori(words)
    for num in res:
        if num == -1:
            print("Draw")
        elif num == 0:
            print("Takahashi")
        else:
            print("Aoki")
