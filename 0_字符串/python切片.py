import time


size = 1000
for _ in range(24):
    s = "*" * size
    start = time.time()
    a = s[:-1]
    # b = s[:-1]
    delta = time.time() - start
    print(f"Size {size:9d}, time={delta:.3f}")
    size *= 2

# !假设O(n)可以1s跑1e7的数据
# !字符串切片大约是 O(n/400) 的复杂度，而字符串比较大约是O(n/200)的复杂度
# !小于1000的长度可以将比较/切片视为常数

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
