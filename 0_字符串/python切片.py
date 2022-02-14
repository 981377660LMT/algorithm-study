import time


size = 1000
for _ in range(24):
    # create string of size "size"
    s = '*' * size

    # now time reverse slice
    start = time.time()
    r = s[2:size]
    delta = time.time() - start

    print(f'Size {size:9d}, time={delta:.3f}')

    # double size of the string
    size *= 2

# 字符串切片似乎是 O(n/1000) 的复杂度，而字符串比较是O(n)的复杂度

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
