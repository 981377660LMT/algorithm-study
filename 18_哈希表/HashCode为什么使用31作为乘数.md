```JAVA
// 获取hashCode "abc".hashCode();
public int hashCode() {
    int h = hash;
    if (h == 0 && value.length > 0) {
        char val[] = value;
        for (int i = 0; i < value.length; i++) {
            h = 31 * h + val[i];
        }
        hash = h;
    }
    return h;
}

```

`Why does Java's hashCode() in String use 31 as a multiplier?`

1. 31 是一个奇质数。
2. 另外在二进制中，2 个 5 次方是 32，那么也就是 `31 * i == (i << 5) - i`。这主要是说乘积运算可以使用位移提升性能，同时目前的 JVM 虚拟机也会自动支持此类的优化。
3. 用超过 5 千个单词计算 hashCode，这个 hashCode 的运算使用 31、33、37、39 和 41 作为乘积，得到的碰撞结果少于 7
