如何模拟实现 java 的 bitset(位图)

Bitset 是 Java 中的一种数据结构。Bitset 中主要存储的是二进制位，做的也都是位运算，每一位只用来存储 0，1 值，主要用于对数据的标记。
BitSet 的默认初始大小是一个 long 数组，一个 long 数组就是 64 个 bit

```JAVA

使用：
private static final int DEFAULT_SIZE = 2 << 24;
/**
* 位数组。数组中的元素只能是 0 或者 1
*/
private BitSet bits = new BitSet(DEFAULT_SIZE);

/////////////////////////////////////////////////////////
实现：

public BitSet(int nbits) {
    // nbits can't be negative; size 0 is OK
    if (nbits < 0)
        throw new NegativeArraySizeException("nbits < 0: " + nbits);

    initWords(nbits);
    sizeIsSticky = true;
}


private void initWords(int nbits) {
    words = new long[wordIndex(nbits-1) + 1];
}
private static final int ADDRESS_BITS_PER_WORD = 6;
private static int wordIndex(int bitIndex) {
    return bitIndex >> ADDRESS_BITS_PER_WORD;
}
```
