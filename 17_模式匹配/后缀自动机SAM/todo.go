package copypasta

import (
	. "fmt"
)

// TODO:
// 一起来学学后缀自动机SAM吧~(*￣▽￣*)~ - Ximena的文章 - 知乎
// https://zhuanlan.zhihu.com/p/548853957
// SAM与倍增
// C-葫芦的考验之定位子串_牛客竞赛字符串专题班SAM（后缀自动机简单应用）习题 (nowcoder.com)
/* 后缀自动机 Suffix automaton (SAM)
https://ac.nowcoder.com/acm/contest/37092


区间本质不同子串个数（与 LCT 结合）https://www.luogu.com.cn/problem/P6292
动态子串出现次数（与 LCT 结合）SPOJ NSUBSTR2 https://www.luogu.com.cn/problem/SP8747
*/
