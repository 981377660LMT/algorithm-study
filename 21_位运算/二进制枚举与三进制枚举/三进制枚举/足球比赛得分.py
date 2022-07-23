# 五支球队
# 胜者得三分 负者得零分
# 踢平得一分
# 求足球比赛可能的得分组合以及组合数

# !三进制枚举


from itertools import combinations


n = int(input())
nums = list(map(int, input().split()))


pairs = []  # 10场比赛
for a, b in combinations(range(5), 2):
    pairs.append((a, b))


res = set()  # 可能的比赛得分组合
for state in range(3 ** len(pairs)):
    points = [0] * 5
    for i in range(10):
        a, b = pairs[i]
        div = state // 3**i
        mod = div % 3
        if mod == 0:
            points[a] += 3
        elif mod == 1:
            points[b] += 3
        else:
            points[a] += 1
            points[b] += 1
    res.add(tuple(sorted(points, reverse=True)))


res1 = "yes" if n == len(res) else "no"
res2 = "yes" if (tuple(nums) in res) else "no"

print(f"{res1} {res2}")
