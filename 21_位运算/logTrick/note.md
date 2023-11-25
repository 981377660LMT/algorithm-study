# LogTrick 解决子数组按位与、按位或、gcd、lcm 问题

**固定子数组右端点，最多对应 log 段 op 值不同的子数组**

原理：固定右端点时，向左扩展，GCD 要么不变，要么至少减半，所以固定右端点时，只有 O(log U) 个 GCD

https://github.com/EndlessCheng/codeforces-go/blob/497f0b7a3f7853cca693027f60b20409729cce79/copypasta/bits.go#L644
