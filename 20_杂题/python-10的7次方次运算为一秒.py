import time

tot_time = 0

for t in range(0, 5):
    st = time.time()
    sum = 0
    for i in range(0, int(1e7)):
        sum += i

    ed = time.time()
    inv = ed - st
    tot_time += inv

print(tot_time / 5)

# 1.0903070449829102
