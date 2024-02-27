邻接表+bitset

https://www.cnblogs.com/alex-wei/p/bitset_yyds.html

1. 处理邻接表保存每种字符的下标集合
2. 匹配前，将答案 bitset 置为全 1
3. 进行子串匹配，匹配到第 i 个字符时，将对应字符下标集合右移 i 位，然后与答案 bitset 进行与操作最后答案 bitset 中 1 的位置就是合法的匹配起点
