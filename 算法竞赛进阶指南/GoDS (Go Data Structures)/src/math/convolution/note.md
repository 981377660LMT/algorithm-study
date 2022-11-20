https://github.dev/EndlessCheng/codeforces-go/blob/3dd70515200872705893d52dc5dad174f2c3b5f3/copypasta/math_ntt.go#L25

https://www.cnblogs.com/newera/p/10076871.html

- FFT 与 NTT 区别
  fft/ntt 都可以求多项式的卷积
  NTT 解决的是`多项式乘法带模数`的情况，可以说有些受模数的限制，数也比较大。
  NTT 算法流程与 FFT 几乎一样，区别在于 FTT 使用 n 次单位根插值，NTT 使用原根的次方进行插值。 NTT 都是整数运算，速度较快，且不会出现精度不够。
  数论变换(number-theoretic transform, NTT）是离散傅里叶变换（DFT）在数论基础上的实现；快速数论变换(fast number-theoretic transform, FNTT）是 快速傅里叶变换（FFT）在数论基础上的实现，是数论变换（NTT）增加分治操作之后的快速算法。
  一般默认“数论变换”是指“快速数论变换”。
- DFT: Discrete Fourier Transform
  将多项式的系数表达法转化为点值表达法的过程
- IDFT: Inverse Discrete Fourier Transform
  将点值表达法转化为系数表达法的过程
- FFT: DFT -> 复数点值相乘 -> IDFT
  快速傅里叶变换
