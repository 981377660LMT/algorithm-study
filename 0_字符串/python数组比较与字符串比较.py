import time

# !数组比较需要除以50-200的常数
res = 0.0
for _ in range(100000):
    arr1 = [0] * 1000 + [1]
    arr2 = [0] * 1000 + [2]
    time1 = time.time()
    tmp = arr1 < arr2
    time2 = time.time()
    res += time2 - time1
print(res)  # 0.13074040412902832


# !字符串比较需要除以250-1000的常数
res = 0.0
for _ in range(100000):
    arr1 = "0" * 1000 + "1"
    arr2 = "0" * 1000 + "2"
    time1 = time.time()
    tmp = arr1 < arr2
    time2 = time.time()
    res += time2 - time1
print(res)  # 0.020961999893188477
