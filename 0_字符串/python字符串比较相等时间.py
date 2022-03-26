import time


size = 1000
for _ in range(24):

    start = time.time()
    for _ in range(int(1e7)):
        size += 1
    delta = time.time() - start
    print(f'time={delta:.3f}')

    start = time.time()
    assert 'a' * int(1e7) == 'a' * int(1e7)
    delta = time.time() - start

    print(f'time={delta:.3f}')


# 字符串比较相等切片似乎是 O(n/200) 的复杂度

# Size      1000, time=0.000
# Size      2000, time=0.000
# Size      4000, time=0.000
# Size      8000, time=0.000
# Size     16000, time=0.000
# Size     32000, time=0.000
# Size     64000, time=0.000
# Size    128000, time=0.000
# Size    256000, time=0.000
# Size    512000, time=0.000
# Size   1024000, time=0.000
# Size   2048000, time=0.000
# Size   4096000, time=0.001
# Size   8192000, time=0.001
# Size  16384000, time=0.004
# Size  32768000, time=0.007
# Size  65536000, time=0.014
# Size 131072000, time=0.027
# Size 262144000, time=0.058
# Size 524288000, time=0.111
# Size 1048576000, time=0.225
# Size 2097152000, time=0.457
# Size 4194304000, time=1.677
# Size 8388608000, time=10.228
