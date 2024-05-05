# 加工一个模块需要拧 n 个螺丝。
# 但是，每隔一段固定的时间，
# 小 E 的老板就会到小 E 的工位上收走一个未完工的模块。
# 在这期间，只够小 E 拧 k 个螺丝。

# 小 E 的老板刚刚离开。
# 问老板接下来第几次来的时候小 E 才可能有一个完工的模块？
# 假设小 E 的老板极力不想让小 E 达成这件事。
# 如果小 E 无论如何也不能有一个完工的模块，那么输出一行 Poor E.S.!。
FAIL = "Poor E.S.!"

# n, k = map(int, input().split())

# if k == 1:
#     if n == 1:
#         print(1)
#     else:
#         print(FAIL)
#     exit(0)

# if n == k:
#     print(1)
#     exit(0)


# WA
n, k = map(int, input().split())
if k == 1:
    if n == 1:
        print(1)
    else:
        print("Poor E.S.!")
else:
    if k == n:
        print(1)
        exit()
    v = 0
    cnt = 0
    while True:
        if v + k >= n:
            print(cnt + 1)
            exit()
        new_v = (v + k) // 2  # 为啥是//2
        if new_v == v:
            print("Poor E.S.!")
            exit()
        v = new_v
        cnt += 1

# 把鸡蛋放在所有不同的篮子里 => 1 1 1 1 1 1 1
# 放一个等差数列 => 1 2 3 4 5
