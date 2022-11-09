# !模拟费用流
# https://www.cnblogs.com/Nero-Claudius/p/Simulation_of_CostFlow3.html#_
# https://www.cnblogs.com/Forever-666/p/14623282.html
# https://www.mina.moe/archives/11762
# https://blog.csdn.net/wyy603/article/details/105038432
# https://www.xht37.com/tag/%E6%A8%A1%E6%8B%9F%E8%B4%B9%E7%94%A8%E6%B5%81/


from heapq import heappop, heappush
from typing import List

INF = int(1e18)


class Event:
    __slots__ = ("pos", "kind", "count", "wage")

    def __init__(self, pos: int, kind: int, count: int, wage: int):
        self.pos = pos
        """位置"""
        self.kind = kind
        """类型 0:鼠 1:洞"""
        self.count = count
        """每个位置上的数量"""
        self.wage = wage
        """进入洞的花费"""


class Pair:
    __slots__ = ("cost", "count")

    def __init__(self, cost: int, count: int):
        self.cost = cost
        self.count = count

    def __lt__(self, other: "Pair"):
        return self.cost < other.cost


class Solution:
    def minimumTotalDistance(self, mouses: List[int], holes: List[List[int]]) -> int:
        # https://m-sea-blog.com/archives/uoj455s
        n, m = len(mouses), len(holes)
        sum_ = 0
        mouses.sort()
        holes.sort()
        mousePos = [INF] * int(1e5)
        holePos = [INF] * int(1e5)
        holeCost = [INF] * int(1e5)
        holeCount = [INF] * int(1e5)
        for i in range(1, n + 1):
            mousePos[i] = mouses[i - 1]
        for i in range(1, m + 1):
            holePos[i] = holes[i - 1][0]
            holeCost[i] -= INF
            holeCount[i] = holes[i - 1][1]
            sum_ += holeCost[i]
        if sum_ < 0:
            return -1

        return solve(n, m, mousePos, holePos, holeCost, holeCount)


# #include <bits/stdc++.h>
# #define pb push_back
# #define fi first
# #define se second
# using std::vector; using std::pair; using std::min;
# using ll = long long;
# int n, m;
# struct node { int op, x, w, c; };
# vector<node> vec;
# std::priority_queue<pair<ll, int>, vector<pair<ll, int>>, std::greater<pair<ll, int>>> Q0, Q1;
# int main() {
# 	scanf("%d%d", &n, &m);
# 	for (int i = 1; i <= n; i++) {
# 		int x; scanf("%d", &x);
# 		vec.pb({1, x, 0, 1});
# 	}
# 	ll sum = 0;
# 	for (int i = 1; i <= m; i++) {
# 		int x, w, c; scanf("%d%d%d", &x, &w, &c);
# 		vec.pb({2, x, w, c}); sum += c;
# 	}
# 	if (sum < n) { puts("-1"); return 0; }
# 	std::sort(vec.begin(), vec.end(), [](node a, node b) { return a.x < b.x; });
# 	Q0.push({1e12, 1e5}), Q1.push({1e12, 1e5});
# 	ll ans = 0;
# 	for (node u : vec)
# 		if (u.op == 1) {
# 			pair<ll, int> tmp = Q1.top(); Q1.pop();
# 			ans += tmp.fi + u.x; tmp.se--;
# 			Q0.push({-u.x * 2 - tmp.fi, 1});
# 			if (tmp.se) Q1.push(tmp);
# 		} else {
# 			int flow = 0;
# 			while (flow < u.c && Q0.top().fi + u.x + u.w < 0) {
# 				pair<ll, int> tmp = Q0.top(); Q0.pop();
# 				int f = min(tmp.se, u.c - flow);
# 				ans += (ll)f * (tmp.fi + u.x + u.w);
# 				flow += f;
# 				tmp.se -= f;
# 				Q1.push({-u.x * 2 - tmp.fi, f});
# 				if (tmp.se) Q0.push(tmp);
# 			}
# 			if (flow) Q0.push({-u.x - u.w, flow});
# 			if (flow < u.c) Q1.push({-u.x + u.w, u.c - flow});
# 		}
# 	printf("%lld\n", ans);
# 	return 0;
# }


def solve(
    n: int,
    m: int,
    mousePos: List[int],
    holePos: List[int],
    holeWage: List[int],
    holeCount: List[int],
) -> int:
    events: List["Event"] = []
    Q0: List["Pair"] = []
    Q1: List["Pair"] = []
    for i in range(1, n + 1):
        events.append(Event(mousePos[i], 0, 1, 0))
    for i in range(1, m + 1):
        events.append(Event(holePos[i], 1, holeCount[i], holeWage[i]))

    events.sort(key=lambda x: x.pos)
    heappush(Q0, Pair(int(1e12), int(1e5)))
    heappush(Q1, Pair(int(1e12), int(1e5)))
    res = 0
    for event in events:
        if event.kind == 0:
            tmp = heappop(Q1)
            res += tmp.cost + event.pos
            tmp.count -= 1
            heappush(Q0, Pair(-event.pos * 2 - tmp.cost, 1))
            if tmp.count:
                heappush(Q1, tmp)
        else:
            flow = 0
            while flow < event.count and Q0[0].cost + event.pos + event.wage < 0:
                tmp = heappop(Q0)
                f = min(tmp.count, event.count - flow)
                res += f * (tmp.cost + event.pos + event.wage)
                flow += f
                tmp.count -= f
                heappush(Q1, Pair(-event.pos * 2 - tmp.cost, f))
                if tmp.count:
                    heappush(Q0, tmp)
            if flow:
                heappush(Q0, Pair(-event.pos - event.wage, flow))
            if flow < event.count:
                heappush(Q1, Pair(-event.pos + event.wage, event.count - flow))
    return res


print(Solution().minimumTotalDistance(mouses=[0, 4, 6], holes=[[2, 2], [6, 2]]))
# https://uoj.ac/submission/575876
# https://uoj.ac/submissions?problem_id=455&min_score=100&max_score=100&language=c%2B%2B

print(solve(4, 4, [0] + [8, 8, 9, 10], [0] + [1, 3, 5, 10], [0] + [4, 0, 0, 2], [0] + [1, 1, 1, 1]))
