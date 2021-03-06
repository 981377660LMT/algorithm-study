# 给定 n 个变量和 m 个不等式。其中 n 小于等于 26，变量分别用前 n 的大写英文字母表示。
# 不等式之间具有传递性，即若 A>B 且 B>C，则 A>C。

# 请从前往后遍历每对关系，每次遍历时判断：
# 如果能够确定全部关系且无矛盾，则结束循环，输出确定的次序；
# 如果发生矛盾，则结束循环，输出有矛盾；
# 如果循环结束时没有发生上述两种情况，则输出无定解。
