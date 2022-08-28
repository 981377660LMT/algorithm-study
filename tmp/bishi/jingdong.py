# # TODO 力扣奇偶数组原题

# from collections import Counter


# n = int(input())
# nums = list(map(int, input().split()))
# if n == 1:
#     print(0)
#     exit(0)

# if n == 2:
#     print(1 if nums[0] == nums[1] else 0)
#     exit(0)

# odd = []
# even = []
# for i in range(n):
#     if i & 1:
#         odd.append(nums[i])
#     else:
#         even.append(nums[i])

# counter1, counter2 = Counter(odd), Counter(even)
# c1 = counter1.most_common()
# c2 = counter2.most_common()

# if len(c1) <= 1:
#     c1.append((-1, 0))
# if len(c2) <= 1:
#     c2.append((-2, 0))

# res = n
# for num1, count1 in c1[:2]:
#     for num2, count2 in c2[:2]:
#         if num1 == num2:
#             res = min(res, n - max(count1, count2))
#         else:
#             res = min(res, n - count1 - count2)
# print(res)
