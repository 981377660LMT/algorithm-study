// https://github.com/EndlessCheng/codeforces-go/blob/ff168d8e767e1d09ce47b4107d5c2e511b8bf41d/copypasta/common.go#L1-L2483
// TODO：solve problems using rust in mac.
package main

import (
	"fmt"
	"math/bits"
	"reflect"
	"time"
	"unsafe"
)

/*
编程入门题单
https://leetcode.cn/studyplan/primers-list/
https://leetcode.cn/studyplan/programming-skills/

我整理的分类题单
https://leetcode.cn/circle/discuss/0viNMK/ 滑动窗口（定长/不定长/多指针）
https://leetcode.cn/circle/discuss/SqopEo/ 二分算法（二分答案/最小化最大值/最大化最小值/第K小）
https://leetcode.cn/circle/discuss/9oZFK9/ 单调栈（矩形系列/字典序最小/贡献法）
https://leetcode.cn/circle/discuss/YiXPXW/ 网格图（DFS/BFS/综合应用）
https://leetcode.cn/circle/discuss/dHn9Vk/ 位运算（基础/性质/拆位/试填/恒等式/贪心/脑筋急转弯）
https://leetcode.cn/circle/discuss/01LUak/ 图论算法（DFS/BFS/拓扑排序/最短路/最小生成树/二分图/基环树/欧拉路径）
https://leetcode.cn/circle/discuss/tXLS3i/ 动态规划（入门/背包/状态机/划分/区间/状压/数位/数据结构优化/树形/博弈/概率期望）
https://leetcode.cn/circle/discuss/mOr1u6/ 常用数据结构（前缀和/差分/栈/队列/堆/字典树/并查集/树状数组/线段树）
https://leetcode.cn/circle/discuss/IYT3ss/ 数学算法（数论/组合/概率期望/博弈/计算几何/随机算法）

力扣题目分类汇总
https://leetcode.cn/circle/article/04PVPY/
https://leetcode.cn/circle/discuss/vEFf96/

## 字符串基础
https://codeforces.com/problemset/problem/1101/B
https://leetcode.cn/problems/apply-operations-to-make-string-empty/

## 暴力枚举
https://codeforces.com/problemset/problem/681/B 1300
- [2207. 字符串中最多数目的子序列](https://leetcode.cn/problems/maximize-number-of-subsequences-in-a-string/) 1550

## 枚举右，维护左
- [1. 两数之和](https://leetcode.cn/problems/two-sum/)
   - https://codeforces.com/problemset/problem/702/B
- [1512. 好数对的数目](https://leetcode.cn/problems/number-of-good-pairs/) 1161 经典题
    - https://leetcode.cn/problems/sum-of-digit-differences-of-all-pairs/
- [2815. 数组中的最大数对和](https://leetcode.cn/problems/max-pair-sum-in-an-array/) 1295
- [2748. 美丽下标对的数目](https://leetcode.cn/problems/number-of-beautiful-pairs/) 1301
- [219. 存在重复元素 II](https://leetcode.cn/problems/contains-duplicate-ii/)
- [121. 买卖股票的最佳时机](https://leetcode.cn/problems/best-time-to-buy-and-sell-stock/)
- [2342. 数位和相等数对的最大和](https://leetcode.cn/problems/max-sum-of-a-pair-with-equal-sum-of-digits/) 1309
- [1679. K 和数对的最大数目](https://leetcode.cn/problems/max-number-of-k-sum-pairs/) 1346
- [1010. 总持续时间可被 60 整除的歌曲](https://leetcode.cn/problems/pairs-of-songs-with-total-durations-divisible-by-60/) 1377
- [2971. 找到最大周长的多边形](https://leetcode.cn/problems/find-polygon-with-the-largest-perimeter/) 1521
- [2874. 有序三元组中的最大值 II](https://leetcode.cn/problems/maximum-value-of-an-ordered-triplet-ii/) 1583
    巧妙安排更新顺序，使得 ans，pre_max 只能使用之前的值，从而符合 i<j<k 的要求
- [1014. 最佳观光组合](https://leetcode.cn/problems/best-sightseeing-pair/) 1730
- [1814. 统计一个数组中好对子的数目](https://leetcode.cn/problems/count-nice-pairs-in-an-array/) 1738
- [454. 四数相加 II](https://leetcode.cn/problems/4sum-ii/)
- [1214. 查找两棵二叉搜索树之和](https://leetcode.cn/problems/two-sum-bsts/)（会员题）
- [2613. 美数对](https://leetcode.cn/problems/beautiful-pairs/)（会员题）
- [2964. 可被整除的三元组数量](https://leetcode.cn/problems/number-of-divisible-triplet-sums/)（会员题）
https://leetcode.com/discuss/interview-question/3685049/25-variations-of-Two-sum-question
异或 https://codeforces.com/problemset/problem/1800/F 1900

## 枚举右，维护左：需要维护两种值（pair）
https://codeforces.com/problemset/problem/1895/C 1400
https://codeforces.com/contest/1931/problem/D
https://leetcode.cn/problems/count-beautiful-substrings-ii/

哈希表
- [2260. 必须拿起的最小连续卡牌数](https://leetcode.cn/problems/minimum-consecutive-cards-to-pick-up/) 1365
- [982. 按位与为零的三元组](https://leetcode.cn/problems/triples-with-bitwise-and-equal-to-zero/) 2085
- [面试题 16.21. 交换和](https://leetcode.cn/problems/sum-swap-lcci/)

前缀和
- [1732. 找到最高海拔](https://leetcode.cn/problems/find-the-highest-altitude/)
- [303. 区域和检索 - 数组不可变](https://leetcode.cn/problems/range-sum-query-immutable/)
- [1310. 子数组异或查询](https://leetcode.cn/problems/xor-queries-of-a-subarray/)
- [2615. 等值距离和](https://leetcode.cn/problems/sum-of-distances/) 1793
- [2602. 使数组元素全部相等的最少操作次数](https://leetcode.cn/problems/minimum-operations-to-make-all-array-elements-equal/) 1903
- [2955. Number of Same-End Substrings](https://leetcode.cn/problems/number-of-same-end-substrings/)（会员题）
https://codeforces.com/problemset/problem/466/C

前缀和+哈希表（双变量思想）
- [930. 和相同的二元子数组](https://leetcode.cn/problems/binary-subarrays-with-sum/) 1592  *同 560，但是数据范围小，存在滑窗做法
- [560. 和为 K 的子数组](https://leetcode.cn/problems/subarray-sum-equals-k/)
- [1524. 和为奇数的子数组数目](https://leetcode.cn/problems/number-of-sub-arrays-with-odd-sum/) 1611
- [974. 和可被 K 整除的子数组](https://leetcode.cn/problems/subarray-sums-divisible-by-k/) 1676
   - 变形：乘积可以被 k 整除
   - a[i] = gcd(a[i], k) 之后窗口乘积是 k 的倍数就行，不会乘爆
- [523. 连续的子数组和](https://leetcode.cn/problems/continuous-subarray-sum/)
- [3026. 最大好子数组和](https://leetcode.cn/problems/maximum-good-subarray-sum/) 1817
- [525. 连续数组](https://leetcode.cn/problems/contiguous-array/) *转换
- [1124. 表现良好的最长时间段](https://leetcode.cn/problems/longest-well-performing-interval/) 1908 *转换
- [2488. 统计中位数为 K 的子数组](https://leetcode.cn/problems/count-subarrays-with-median-k/) 1999 *转换
- [1590. 使数组和能被 P 整除](https://leetcode.cn/problems/make-sum-divisible-by-p/) 2039
- [2949. 统计美丽子字符串 II](https://leetcode.cn/problems/count-beautiful-substrings-ii/) 2445
- [面试题 17.05. 字母与数字](https://leetcode.cn/problems/find-longest-subarray-lcci/)
- [1983. 范围和相等的最宽索引对](https://leetcode.cn/problems/widest-pair-of-indices-with-equal-range-sum/)（会员题）
- [2489. 固定比率的子字符串数](https://leetcode.cn/problems/number-of-substrings-with-fixed-ratio/)（会员题）
https://atcoder.jp/contests/abc233/tasks/abc233_d
交错前缀和 https://codeforces.com/contest/1915/problem/E
https://codeforces.com/problemset/problem/1446/D1 2600 转换
https://www.luogu.com.cn/problem/AT_joisc2014_h 三个字母映射到一些大整数上，从而区分开

前缀和思想 LC1523 https://leetcode.cn/problems/count-odd-numbers-in-an-interval-range/
有点数形结合 https://codeforces.com/problemset/problem/1748/C

前缀和的前缀和（二重前缀和）
LC2281 https://leetcode.cn/problems/sum-of-total-strength-of-wizards/
https://atcoder.jp/contests/abc058/tasks/arc071_b

前缀和+异或
- [1177. 构建回文串检测](https://leetcode.cn/problems/can-make-palindrome-from-substring/) 1848
- [1371. 每个元音包含偶数次的最长子字符串](https://leetcode.cn/problems/find-the-longest-substring-containing-vowels-in-even-counts/) 2041
- [1542. 找出最长的超赞子字符串](https://leetcode.cn/problems/find-longest-awesome-substring/) 2222
- [1915. 最美子字符串的数目](https://leetcode.cn/problems/number-of-wonderful-substrings/) 2235
- [2791. 树中可以形成回文的路径数](https://leetcode.cn/problems/count-paths-that-can-form-a-palindrome-in-a-tree/) 2677
模 3 & 字符集大小为 n https://codeforces.com/problemset/problem/1418/G 2500
https://atcoder.jp/contests/abc295/tasks/abc295_d
https://ac.nowcoder.com/acm/contest/75174/E

https://leetcode.cn/problems/find-longest-subarray-lcci/
https://codeforces.com/problemset/problem/1296/C

## 前后缀分解（右边数字为难度分）
部分题目也可以用状态机 DP 解决
- [42. 接雨水](https://leetcode.cn/problems/trapping-rain-water/)（[视频讲解](https://www.bilibili.com/video/BV1Qg411q7ia/?t=3m05s)）
  注：带修改的接雨水 https://codeforces.com/gym/104821/problem/M
  - https://www.zhihu.com/question/627281278/answer/3280684055
  - 全排列接雨水 https://atcoder.jp/contests/tenka1-2015-final/tasks/tenka1_2015_final_e
- [123. 买卖股票的最佳时机 III](https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-iii/) 拆分成两个 121 题
- [1422. 分割字符串的最大得分](https://leetcode.cn/problems/maximum-score-after-splitting-a-string/) 1238
- [2256. 最小平均差](https://leetcode.cn/problems/minimum-average-difference/) 1395
- [1493. 删掉一个元素以后全为 1 的最长子数组](https://leetcode.cn/problems/longest-subarray-of-1s-after-deleting-one-element/) 1423
- [845. 数组中的最长山脉](https://leetcode.cn/problems/longest-mountain-in-array/) 1437 *也可以分组循环
- [2909. 元素和最小的山形三元组 II](https://leetcode.cn/problems/minimum-sum-of-mountain-triplets-ii/) 1479
- [2483. 商店的最少代价](https://leetcode.cn/problems/minimum-penalty-for-a-shop/) 1495
- [1525. 字符串的好分割数目](https://leetcode.cn/problems/number-of-good-ways-to-split-a-string/) 1500
- [3096. 得到更多分数的最少关卡数目](https://leetcode.cn/problems/minimum-levels-to-gain-more-points/) 1501
- [2874. 有序三元组中的最大值 II](https://leetcode.cn/problems/maximum-value-of-an-ordered-triplet-ii/) 1583
- [1031. 两个非重叠子数组的最大和](https://leetcode.cn/problems/maximum-sum-of-two-non-overlapping-subarrays/) 1680
- [689. 三个无重叠子数组的最大和](https://leetcode.cn/problems/maximum-sum-of-3-non-overlapping-subarrays/)
- [2420. 找到所有好下标](https://leetcode.cn/problems/find-all-good-indices/) 1695
- [2100. 适合野炊的日子](https://leetcode.cn/problems/find-good-days-to-rob-the-bank/) 1702
- [1653. 使字符串平衡的最少删除次数](https://leetcode.cn/problems/minimum-deletions-to-make-string-balanced/) 1794
- [926. 将字符串翻转到单调递增](https://leetcode.cn/problems/flip-string-to-monotone-increasing/)
  - https://codeforces.com/problemset/problem/180/C 1400
  - https://codeforces.com/problemset/problem/846/A 1500
- [1477. 找两个和为目标值且不重叠的子数组](https://leetcode.cn/problems/find-two-non-overlapping-sub-arrays-each-with-target-sum/) 1851
- [1671. 得到山形数组的最少删除次数](https://leetcode.cn/problems/minimum-number-of-removals-to-make-mountain-array/) 1913 *DP
- [238. 除自身以外数组的乘积](https://leetcode.cn/problems/product-of-array-except-self/) ~2000
- [1888. 使二进制字符串字符交替的最少反转次数](https://leetcode.cn/problems/minimum-number-of-flips-to-make-the-binary-string-alternating/) 2006
- [2906. 构造乘积矩阵](https://leetcode.cn/problems/construct-product-matrix/) 2075
- [2167. 移除所有载有违禁货物车厢所需的最少时间](https://leetcode.cn/problems/minimum-time-to-remove-all-cars-containing-illegal-goods/) 2219 *DP
- [2484. 统计回文子序列数目](https://leetcode.cn/problems/count-palindromic-subsequences/) 2223
- [2163. 删除元素后和的最小差值](https://leetcode.cn/problems/minimum-difference-in-sums-after-removal-of-elements/) 2225
- [2565. 最少得分子序列](https://leetcode.cn/problems/subsequence-with-the-minimum-score/) 2432
- [2552. 统计上升四元组](https://leetcode.cn/problems/count-increasing-quadruplets/) 2433
- [3003. 执行操作后的最大分割数量](https://leetcode.cn/problems/maximize-the-number-of-partitions-after-operations/) 3039
- [487. 最大连续 1 的个数 II](https://leetcode.cn/problems/max-consecutive-ones-ii/)（会员题）
- [1746. 经过一次操作后的最大子数组和](https://leetcode.cn/problems/maximum-subarray-sum-after-one-operation/)（会员题）
https://codeforces.com/problemset/problem/1178/B 1300
https://codeforces.com/problemset/problem/1443/B 1300
https://codeforces.com/problemset/problem/1706/C 1400
https://codeforces.com/problemset/problem/1029/C 1600
https://codeforces.com/problemset/problem/1837/F 2400
昆明 2024：至多修改一个子数组 [L,R] ：把元素都加上 k，最大化整个数组的 GCD
- 预处理前后缀 GCD，由于前缀 GCD 只有 O(logU) 个不同的值，可以只枚举 O(logU) 个 L 和 O(n) 个 R，
- 枚举 R 的同时计算修改后的子数组 GCD，然后和前后缀 GCD 求 GCD

#### 定长滑动窗口（右边数字为难度分）
- [1456. 定长子串中元音的最大数目](https://leetcode.cn/problems/maximum-number-of-vowels-in-a-substring-of-given-length/) 1263
- [2269. 找到一个数字的 K 美丽值](https://leetcode.cn/problems/find-the-k-beauty-of-a-number/) 1280
- [1984. 学生分数的最小差值](https://leetcode.cn/problems/minimum-difference-between-highest-and-lowest-of-k-scores/) 1306
- [643. 子数组最大平均数 I](https://leetcode.cn/problems/maximum-average-subarray-i/)
- [1343. 大小为 K 且平均值大于等于阈值的子数组数目](https://leetcode.cn/problems/number-of-sub-arrays-of-size-k-and-average-greater-than-or-equal-to-threshold/) 1317
- [2090. 半径为 k 的子数组平均值](https://leetcode.cn/problems/k-radius-subarray-averages/) 1358
- [2379. 得到 K 个黑块的最少涂色次数](https://leetcode.cn/problems/minimum-recolors-to-get-k-consecutive-black-blocks/) 1360
- [1652. 拆炸弹](https://leetcode.cn/problems/defuse-the-bomb/) 1417
- [1052. 爱生气的书店老板](https://leetcode.cn/problems/grumpy-bookstore-owner/) 1418
- [2841. 几乎唯一子数组的最大和](https://leetcode.cn/problems/maximum-sum-of-almost-unique-subarray/) 1546
- [2461. 长度为 K 子数组中的最大和](https://leetcode.cn/problems/maximum-sum-of-distinct-subarrays-with-length-k/) 1553
- [1423. 可获得的最大点数](https://leetcode.cn/problems/maximum-points-you-can-obtain-from-cards/) 1574
- [2134. 最少交换次数来组合所有的 1 II](https://leetcode.cn/problems/minimum-swaps-to-group-all-1s-together-ii/) 1748
- [2653. 滑动子数组的美丽值](https://leetcode.cn/problems/sliding-subarray-beauty/) 1786
- [567. 字符串的排列](https://leetcode.cn/problems/permutation-in-string/)
- [438. 找到字符串中所有字母异位词](https://leetcode.cn/problems/find-all-anagrams-in-a-string/)
- [2156. 查找给定哈希值的子串](https://leetcode.cn/problems/find-substring-with-given-hash-value/) 2063
- [2953. 统计完全子字符串](https://leetcode.cn/problems/count-complete-substrings/) 2449 *分组循环
- [346. 数据流中的移动平均值](https://leetcode.cn/problems/moving-average-from-data-stream/)（会员题）
- [1100. 长度为 K 的无重复字符子串](https://leetcode.cn/problems/find-k-length-substrings-with-no-repeated-characters/)（会员题）
- [1852. 每个子数组的数字种类数](https://leetcode.cn/problems/distinct-numbers-in-each-subarray/)（会员题）
- [2067. 等计数子串的数量](https://leetcode.cn/problems/number-of-equal-count-substrings/)（会员题）
- [2107. 分享 K 个糖果后独特口味的数量](https://leetcode.cn/problems/number-of-unique-flavors-after-sharing-k-candies/)（会员题）
https://codeforces.com/problemset/problem/608/B 1500
https://codeforces.com/problemset/problem/69/E 1800
https://codeforces.com/problemset/problem/371/E 2000

#### 不定长滑动窗口（求最长/最大）
- [3. 无重复字符的最长子串](https://leetcode.cn/problems/longest-substring-without-repeating-characters/)
   - 翻转至多一个任意子串后的无重复字符的最长子串 https://codeforces.com/contest/1234/problem/F
- [3090. 每个字符最多出现两次的最长子字符串](https://leetcode.cn/problems/maximum-length-substring-with-two-occurrences/) 1329
- [1493. 删掉一个元素以后全为 1 的最长子数组](https://leetcode.cn/problems/longest-subarray-of-1s-after-deleting-one-element/) 1423
- [1208. 尽可能使字符串相等](https://leetcode.cn/problems/get-equal-substrings-within-budget/) 1497
- [2730. 找到最长的半重复子字符串](https://leetcode.cn/problems/find-the-longest-semi-repetitive-substring/) 1502
- [904. 水果成篮](https://leetcode.cn/problems/fruit-into-baskets/) 1516
- [1695. 删除子数组的最大得分](https://leetcode.cn/problems/maximum-erasure-value/) 1529
- [2958. 最多 K 个重复元素的最长子数组](https://leetcode.cn/problems/length-of-longest-subarray-with-at-most-k-frequency/) 1535
- [2024. 考试的最大困扰度](https://leetcode.cn/problems/maximize-the-confusion-of-an-exam/) 1643
- [1004. 最大连续1的个数 III](https://leetcode.cn/problems/max-consecutive-ones-iii/) 1656
- [1438. 绝对差不超过限制的最长连续子数组](https://leetcode.cn/problems/longest-continuous-subarray-with-absolute-diff-less-than-or-equal-to-limit/) 1672  *需要 SortedList
- [2401. 最长优雅子数组](https://leetcode.cn/problems/longest-nice-subarray/) 1750 *位运算
- [1658. 将 x 减到 0 的最小操作数](https://leetcode.cn/problems/minimum-operations-to-reduce-x-to-zero/) 1817
    - https://codeforces.com/problemset/problem/1692/E 1200
- [1838. 最高频元素的频数](https://leetcode.cn/problems/frequency-of-the-most-frequent-element/) 1876
- [2516. 每种字符至少取 K 个](https://leetcode.cn/problems/take-k-of-each-character-from-left-and-right/) 1948
- [2831. 找出最长等值子数组](https://leetcode.cn/problems/find-the-longest-equal-subarray/) 1976
- [2106. 摘水果](https://leetcode.cn/problems/maximum-fruits-harvested-after-at-most-k-steps/) 2062
- [1610. 可见点的最大数目](https://leetcode.cn/problems/maximum-number-of-visible-points/) 2147
- [2781. 最长合法子字符串的长度](https://leetcode.cn/problems/length-of-the-longest-valid-substring/) 2204
- [2968. 执行操作使频率分数最大](https://leetcode.cn/problems/apply-operations-to-maximize-frequency-score/) 2444
- [395. 至少有 K 个重复字符的最长子串](https://leetcode.cn/problems/longest-substring-with-at-least-k-repeating-characters/)
- [1763. 最长的美好子字符串](https://leetcode.cn/problems/longest-nice-substring/)
- [424. 替换后的最长重复字符](https://leetcode.cn/problems/longest-repeating-character-replacement/) *有些特殊
- [159. 至多包含两个不同字符的最长子串](https://leetcode.cn/problems/longest-substring-with-at-most-two-distinct-characters/)（会员题）
- [340. 至多包含 K 个不同字符的最长子串](https://leetcode.cn/problems/longest-substring-with-at-most-k-distinct-characters/)（会员题）
与单调队列结合 https://www.luogu.com.cn/problem/P3594
https://codeforces.com/problemset/problem/1873/F 1300

#### 不定长滑动窗口（求最短/最小）
- [209. 长度最小的子数组](https://leetcode.cn/problems/minimum-size-subarray-sum/)
- [1234. 替换子串得到平衡字符串](https://leetcode.cn/problems/replace-the-substring-for-balanced-string/) 1878
- [2875. 无限数组的最短子数组](https://leetcode.cn/problems/minimum-size-subarray-in-infinite-array/) 1914
- [1574. 删除最短的子数组使剩余数组有序](https://leetcode.cn/problems/shortest-subarray-to-be-removed-to-make-array-sorted/) 1932
- [76. 最小覆盖子串](https://leetcode.cn/problems/minimum-window-substring/)
- [面试题 17.18. 最短超串](https://leetcode.cn/problems/shortest-supersequence-lcci/)
https://codeforces.com/problemset/problem/1354/B 1200
https://codeforces.com/problemset/problem/224/B 1500 和最小
https://codeforces.com/problemset/problem/701/C 1500
https://codeforces.com/problemset/problem/1777/C 1700

#### 不定长滑动窗口（求子数组个数）
- [2799. 统计完全子数组的数目](https://leetcode.cn/problems/count-complete-subarrays-in-an-array/) 1398
- [713. 乘积小于 K 的子数组](https://leetcode.cn/problems/subarray-product-less-than-k/)
- [1358. 包含所有三种字符的子字符串数目](https://leetcode.cn/problems/number-of-substrings-containing-all-three-characters/) 1646
- [2962. 统计最大元素出现至少 K 次的子数组](https://leetcode.cn/problems/count-subarrays-where-max-element-appears-at-least-k-times/) 1701
- [LCP 68. 美观的花束](https://leetcode.cn/problems/1GxJYY/)
- [2302. 统计得分小于 K 的子数组数目](https://leetcode.cn/problems/count-subarrays-with-score-less-than-k/) 1808
- [2537. 统计好子数组的数目](https://leetcode.cn/problems/count-the-number-of-good-subarrays/) 1892
- [2762. 不间断子数组](https://leetcode.cn/problems/continuous-subarrays/) 1940
- [2972. 统计移除递增子数组的数目 II](https://leetcode.cn/problems/count-the-number-of-incremovable-subarrays-ii/) 2153
- [1918. 第 K 小的子数组和](https://leetcode.cn/problems/kth-smallest-subarray-sum/)（会员题）*二分答案
- [2743. 计算没有重复字符的子字符串数量](https://leetcode.cn/problems/count-substrings-without-repeating-character/)（会员题）
和至少为 k 的子数组个数 https://atcoder.jp/contests/abc130/tasks/abc130_d
变形：改成子数组 https://codeforces.com/problemset/problem/550/B
其它题目见【前缀和】

#### 滑窗的同时维护数据
https://codeforces.com/problemset/problem/898/D 1600

#### 进阶 多指针滑动窗口
- [930. 和相同的二元子数组](https://leetcode.cn/problems/binary-subarrays-with-sum/) 1592 恰好等于
- [1248. 统计「优美子数组」](https://leetcode.cn/problems/count-number-of-nice-subarrays/) 1624 类似 930
- [2563. 统计公平数对的数目](https://leetcode.cn/problems/count-the-number-of-fair-pairs/) 1721
- [1477. 找两个和为目标值且不重叠的子数组](https://leetcode.cn/problems/find-two-non-overlapping-sub-arrays-each-with-target-sum/) 1851
- [1712. 将数组分成三个子数组的方案数](https://leetcode.cn/problems/ways-to-split-array-into-three-subarrays/) 2079
- [2444. 统计定界子数组的数目](https://leetcode.cn/problems/count-subarrays-with-fixed-bounds/) 2093
- [1638. 统计只差一个字符的子串数目](https://leetcode.cn/problems/count-substrings-that-differ-by-one-character/) *非暴力做法
- [992. K 个不同整数的子数组](https://leetcode.cn/problems/subarrays-with-k-different-integers/) 2210
- [1989. 捉迷藏中可捕获的最大人数](https://leetcode.cn/problems/maximum-number-of-people-that-can-be-caught-in-tag/)（会员题）

### 多指针
- [1213. 三个有序数组的交集](https://leetcode.cn/problems/intersection-of-three-sorted-arrays/)（会员题）
https://codeforces.com/problemset/problem/1971/F 1600

LC2234 https://leetcode.cn/problems/maximum-total-beauty-of-the-gardens/ 2562
类似 [795. 区间子数组个数](https://leetcode.cn/problems/number-of-subarrays-with-bounded-maximum/) 1817
入门题 https://codeforces.com/problemset/problem/602/B
入门题 https://codeforces.com/problemset/problem/279/B
https://atcoder.jp/contests/abc229/tasks/abc229_d
LC2271 毯子覆盖的最多白色砖块数 需要多思考一点点 https://leetcode.cn/problems/maximum-white-tiles-covered-by-a-carpet/
- https://atcoder.jp/contests/abc098/tasks/arc098_b
较为复杂 https://atcoder.jp/contests/abc294/tasks/abc294_e
      - https://ac.nowcoder.com/acm/contest/62033/D
https://codeforces.com/problemset/problem/1208/B
https://codeforces.com/problemset/problem/1765/D
多指针 https://codeforces.com/problemset/problem/895/B
https://codeforces.com/contest/1833/problem/F
计算有多少子数组，其中有至少 k 个相同的数 https://codeforces.com/problemset/problem/190/D
mex https://atcoder.jp/contests/abc194/tasks/abc194_e
https://codeforces.com/problemset/problem/165/C
- [1099. 小于 K 的两数之和](https://leetcode.cn/problems/two-sum-less-than-k/)（会员题）

双序列双指针
LC88 https://leetcode.cn/problems/merge-sorted-array/
LC360（背向双指针）https://leetcode.cn/problems/sort-transformed-array/
- [986. 区间列表的交集](https://leetcode.cn/problems/interval-list-intersections/) 1542
- [1537. 最大得分](https://leetcode.cn/problems/get-the-maximum-score/) 1961
https://codeforces.com/contest/489/problem/B 1200

相向双指针
题单 https://leetcode.cn/leetbook/read/sliding-window-and-two-pointers/odt2yh/
LC2824 https://leetcode.cn/problems/count-pairs-whose-sum-is-less-than-target/
LC923 https://leetcode.cn/problems/3sum-with-multiplicity/
https://www.facebook.com/codingcompetitions/hacker-cup/2023/practice-round/problems/C

同时用到同向双指针和相向双指针的题
https://atcoder.jp/contests/abc155/tasks/abc155_d
- 相似题目 https://leetcode.cn/problems/kth-smallest-product-of-two-sorted-arrays/

a[i] + b[j] = target 的方案数
a[i] + b[j] < target 的方案数    相向双指针 https://leetcode.cn/problems/count-pairs-whose-sum-is-less-than-target/
                                         https://codeforces.com/problemset/problem/1538/C 1300
                               - [259. 较小的三数之和](https://leetcode.cn/problems/3sum-smaller/)（会员题）
a[i] + b[j] > target 的方案数    同上
a[i] - b[j] = target 的方案数
a[i] - b[j] < target 的方案数    滑窗
a[i] - b[j] > target 的方案数    同上 >= https://atcoder.jp/contests/abc353/tasks/abc353_c
子数组元素和 = < > target 的方案数：用前缀和，转换成上面 a[i] - b[j] 的形式
子序列元素和 = < > target 的方案数：0-1 背包恰好/至多/至少，见 https://www.bilibili.com/video/BV16Y411v7Y6/ 末尾的总结

## 分组循环
https://leetcode.cn/problems/longest-even-odd-subarray-with-threshold/solution/jiao-ni-yi-ci-xing-ba-dai-ma-xie-dui-on-zuspx/
**适用场景**：按照题目要求，数组会被分割成若干组，每一组的判断/处理逻辑是相同的。
**核心思想**：
- 外层循环负责遍历组之前的准备工作（记录开始位置），和遍历组之后的统计工作（更新答案最大值）。
- 内层循环负责遍历组，找出这一组最远在哪结束。
这个写法的好处是，各个逻辑块分工明确，也不需要特判最后一组（易错点）。以我的经验，这个写法是所有写法中最不容易出 bug 的，推荐大家记住。
- [1446. 连续字符](https://leetcode.cn/problems/consecutive-characters/) 1165
- [1869. 哪种连续子字符串更长](https://leetcode.cn/problems/longer-contiguous-segments-of-ones-than-zeros/) 1205
- [1957. 删除字符使字符串变好](https://leetcode.cn/problems/delete-characters-to-make-fancy-string/) 1358
- [674. 最长连续递增序列](https://leetcode.cn/problems/longest-continuous-increasing-subsequence/)
- [978. 最长湍流子数组](https://leetcode.cn/problems/longest-turbulent-subarray/) 1393
- [2110. 股票平滑下跌阶段的数目](https://leetcode.cn/problems/number-of-smooth-descent-periods-of-a-stock/) 1408
- [228. 汇总区间](https://leetcode.cn/problems/summary-ranges/)
- [2760. 最长奇偶子数组](https://leetcode.cn/problems/longest-even-odd-subarray-with-threshold/) 1420
- [1887. 使数组元素相等的减少操作次数](https://leetcode.cn/problems/reduction-operations-to-make-the-array-elements-equal/) 1428
- [845. 数组中的最长山脉](https://leetcode.cn/problems/longest-mountain-in-array/) 1437
- [2038. 如果相邻两个颜色均相同则删除当前颜色](https://leetcode.cn/problems/remove-colored-pieces-if-both-neighbors-are-the-same-color/) 1468
- [1759. 统计同质子字符串的数目](https://leetcode.cn/problems/count-number-of-homogenous-substrings/) 1491
- [3011. 判断一个数组是否可以变为有序](https://leetcode.cn/problems/find-if-array-can-be-sorted/) 1497
- [1578. 使绳子变成彩色的最短时间](https://leetcode.cn/problems/minimum-time-to-make-rope-colorful/) 1574
- [1839. 所有元音按顺序排布的最长子字符串](https://leetcode.cn/problems/longest-substring-of-all-vowels-in-order/) 1580
- [2765. 最长交替子序列](https://leetcode.cn/problems/longest-alternating-subarray/) 1581
- [3105. 最长的严格递增或递减子数组](https://leetcode.cn/problems/longest-strictly-increasing-or-strictly-decreasing-subarray/)
- [467. 环绕字符串中唯一的子字符串](https://leetcode.cn/problems/unique-substrings-in-wraparound-string/) ~1700
- [2948. 交换得到字典序最小的数组](https://leetcode.cn/problems/make-lexicographically-smallest-array-by-swapping-elements/) 2047
- [2393. 严格递增的子数组个数](https://leetcode.cn/problems/count-strictly-increasing-subarrays/)（会员题）
- [2436. 使子数组最大公约数大于一的最小分割数](https://leetcode.cn/problems/minimum-split-into-subarrays-with-gcd-greater-than-one/)（会员题）
- [2495. 乘积为偶数的子数组数](https://leetcode.cn/problems/number-of-subarrays-having-even-product/)（会员题）
- [3063. 链表频率](https://leetcode.cn/problems/linked-list-frequency/)（会员题）
LC1180（会员）https://leetcode.cn/problems/count-substrings-with-only-one-distinct-letter/
LC2257 https://leetcode.cn/problems/count-unguarded-cells-in-the-grid/
- https://atcoder.jp/contests/abc317/tasks/abc317_e
LC2495（会员）逆向思维 https://leetcode.cn/problems/number-of-subarrays-having-even-product/
https://codeforces.com/problemset/problem/1272/C 1200
https://codeforces.com/problemset/problem/1343/C 1200
https://codeforces.com/problemset/problem/1821/C 1300 枚举答案
https://codeforces.com/problemset/problem/1873/F 1300
https://codeforces.com/problemset/problem/1380/C 1400
https://codeforces.com/problemset/problem/620/C 1500
https://codeforces.com/problemset/problem/525/C 1600
https://codeforces.com/problemset/problem/1748/C 1600
https://codeforces.com/problemset/problem/1849/D 1700

### 哨兵
- [1465. 切割后面积最大的蛋糕](https://leetcode.cn/problems/maximum-area-of-a-piece-of-cake-after-horizontal-and-vertical-cuts/) 1445
- [2975. 移除栅栏得到的正方形田地的最大面积](https://leetcode.cn/problems/maximum-square-area-by-removing-fences-from-a-field/) 1873
不是哨兵，但图像类似 [2943. 最大化网格图中正方形空洞的面积](https://leetcode.cn/problems/maximize-area-of-square-hole-in-grid/) 1677

### 巧妙枚举
LC939 https://leetcode.cn/problems/minimum-area-rectangle/
- [1577. 数的平方等于两数乘积的方法数](https://leetcode.cn/problems/number-of-ways-where-square-of-number-is-equal-to-product-of-two-numbers/) 1594
https://codeforces.com/problemset/problem/846/C 1800
https://codeforces.com/problemset/problem/1181/C 1900
https://codeforces.com/problemset/problem/1626/D 2100
https://codeforces.com/problemset/problem/339/E 2700

### 贪心及其证明
- [2587. 重排数组以得到最大前缀分数](https://leetcode.cn/problems/rearrange-array-to-maximize-prefix-score/) 1337
- [455. 分发饼干](https://leetcode.cn/problems/assign-cookies/)
- [1029. 两地调度](https://leetcode.cn/problems/two-city-scheduling/) 1348
- [2165. 重排数字的最小值](https://leetcode.cn/problems/smallest-value-of-the-rearranged-number/) 1362
- [2410. 运动员和训练师的最大匹配数](https://leetcode.cn/problems/maximum-matching-of-players-with-trainers/) 1381
- [3111. 覆盖所有点的最少矩形数目](https://leetcode.cn/problems/minimum-rectangles-to-cover-points/) 1401
- [2139. 得到目标值的最少行动次数](https://leetcode.cn/problems/minimum-moves-to-reach-target-score/) 1417
- [2645. 构造有效字符串的最少插入数](https://leetcode.cn/problems/minimum-additions-to-make-valid-string/) 1478
- [3091. 执行操作使数据元素之和大于等于 K](https://leetcode.cn/problems/apply-operations-to-make-sum-of-array-greater-than-or-equal-to-k/) 1522
- [881. 救生艇](https://leetcode.cn/problems/boats-to-save-people/) 1530
    - https://codeforces.com/problemset/problem/1690/E
    - https://www.lanqiao.cn/problems/4174/learning/?contest_id=135
    - https://codeforces.com/problemset/problem/1765/D
- [2522. 将字符串分割成值不超过 K 的子字符串](https://leetcode.cn/problems/partition-string-into-substrings-with-values-at-most-k/) 1605
- [2086. 喂食仓鼠的最小食物桶数](https://leetcode.cn/problems/minimum-number-of-food-buckets-to-feed-the-hamsters/) 1623 注：原标题是「从房屋收集雨水需要的最少水桶数」
- [2375. 根据模式串构造最小数字](https://leetcode.cn/problems/construct-smallest-number-from-di-string/) 1642
- [2611. 老鼠和奶酪](https://leetcode.cn/problems/mice-and-cheese/) 1663
- [1567. 乘积为正数的最长子数组长度](https://leetcode.cn/problems/maximum-length-of-subarray-with-positive-product/) 1710
- [3085. 成为 K 特殊字符串需要删除的最少字符数](https://leetcode.cn/problems/minimum-deletions-to-make-string-k-special/) 1765
- [2952. 需要添加的硬币的最小数量](https://leetcode.cn/problems/minimum-number-of-coins-to-be-added/) 1784
- [553. 最优除法](https://leetcode.cn/problems/optimal-division/)
- [330. 按要求补齐数组](https://leetcode.cn/problems/patching-array/)
- [2931. 购买物品的最大开销](https://leetcode.cn/problems/maximum-spending-after-buying-items/) 1822
- [2311. 小于等于 K 的最长二进制子序列](https://leetcode.cn/problems/longest-binary-subsequence-less-than-or-equal-to-k/) 1840
- [3035. 回文字符串的最大数量](https://leetcode.cn/problems/maximum-palindromes-after-operations/) 1857
- [3081. 替换字符串中的问号使分数最小](https://leetcode.cn/problems/replace-question-marks-in-string-to-minimize-its-value/) 1905
- [1147. 段式回文](https://leetcode.cn/problems/longest-chunked-palindrome-decomposition/) 1912
- [1686. 石子游戏 VI](https://leetcode.cn/problems/stone-game-vi/) 2001
    - https://codeforces.com/contest/1914/problem/E2 1400
- [2333. 最小差值平方和](https://leetcode.cn/problems/minimum-sum-of-squared-difference/) 2011
    - 有 k%(i+1) 个元素可以多减少 1
- [2136. 全部开花的最早一天](https://leetcode.cn/problems/earliest-possible-day-of-full-bloom/) 2033
- [1648. 销售价值减少的颜色球](https://leetcode.cn/problems/sell-diminishing-valued-colored-balls/) 2050
- todo 复习 [2193. 得到回文串的最少操作次数](https://leetcode.cn/problems/minimum-number-of-moves-to-make-palindrome/) 2091
- todo 复习 [659. 分割数组为连续子序列](https://leetcode.cn/problems/split-array-into-consecutive-subsequences/)
- [1889. 装包裹的最小浪费空间](https://leetcode.cn/problems/minimum-space-wasted-from-packaging/) 2214
- [1505. 最多 K 次交换相邻数位后得到的最小整数](https://leetcode.cn/problems/minimum-possible-integer-after-at-most-k-adjacent-swaps-on-digits/) 2337
- [420. 强密码检验器](https://leetcode.cn/problems/strong-password-checker/)
- [LCP 26. 导航装置](https://leetcode.cn/problems/hSRGyL/)
- [3088. 使字符串反回文](https://leetcode.cn/problems/make-string-anti-palindrome/)（会员题）
https://codeforces.com/problemset/problem/1920/B 1100
https://codeforces.com/problemset/problem/545/D 1300
https://codeforces.com/problemset/problem/1443/B 1300
https://codeforces.com/problemset/problem/388/A 1400
https://codeforces.com/problemset/problem/492/C 1400
https://codeforces.com/problemset/problem/1369/C 1400
	提示 1：前 k 大的数一定可以作为最大值。且尽量把大的数放在 w[i] = 1 的组中，这样可以计入答案两次。
	如果某个前 k 大的数 x 没有作为最大值（其中一个组的最大值是不在前 k 大中的 y），那么把 x 和 y 交换，
	如果 x 是某个组的最小值，那么交换后 y 必然也是最小值，此时答案不变。
	如果 x 不是某个组的最小值（这个组的最小值是 z）：
		   如果 y 交换后变成了最小值，那么答案变大了 x-z。
		   如果 y 交换后也不是最小值，那么答案变大了 x-y。
	无论如何，这样交换都不会使答案变小，因此前 k 大的数一定可以作为最大值。
	提示 2：然后来说最小值。a 的最小值必然要分到某个组中，为了「跳过」尽量多的较小的数，优先把 a 中较小的数分到 w 较大的组中。所以 a 从小到大遍历，w 从大到小遍历。
https://codeforces.com/problemset/problem/1443/C 1400
https://codeforces.com/problemset/problem/1691/C 1400
https://codeforces.com/problemset/problem/864/D 1500
https://codeforces.com/problemset/problem/985/C 1500
https://codeforces.com/problemset/problem/1659/C 1500
https://codeforces.com/problemset/problem/1759/E 1500
https://codeforces.com/problemset/problem/1873/G 1500
https://codeforces.com/problemset/problem/913/C 1600
https://codeforces.com/problemset/problem/1707/A 1600 倒序思维
https://codeforces.com/problemset/problem/1157/C2 1700
https://codeforces.com/problemset/problem/1661/C 1700 奇数天+1 偶数天 +2
https://codeforces.com/problemset/problem/3/B 1900
https://codeforces.com/problemset/problem/1479/B1 1900
https://codeforces.com/problemset/problem/1804/D 2000
https://codeforces.com/problemset/problem/1479/B2 2100
    https://www.luogu.com.cn/blog/wsyhb/post-ti-xie-cf1479b1-painting-the-array-i
https://codeforces.com/problemset/problem/442/C 2500
    如果 x>=y<=z，那么删除 y 最优
    结束后剩下一个长为 m 的 /\ 形状的序列，由于无法取到最大值和次大值，那么加上剩下最小的 m-2 个数
https://atcoder.jp/contests/arc147/tasks/arc147_e 难
https://www.luogu.com.cn/problem/P1016
https://www.luogu.com.cn/problem/UVA11384 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=25&page=show_problem&problem=2379

数学思维
https://codeforces.com/problemset/problem/23/C 2500
- https://codeforces.com/problemset/problem/798/D 2400

### 乘法贪心
https://codeforces.com/problemset/problem/45/I 1400
https://codeforces.com/problemset/problem/934/A 1400
最大 3 个数的乘积
最大 k 个数的乘积
删除一个数后，最小化最大 k 个数的乘积

### 区间贪心
- [435. 无重叠区间](https://leetcode.cn/problems/non-overlapping-intervals/)
- [452. 用最少数量的箭引爆气球](https://leetcode.cn/problems/minimum-number-of-arrows-to-burst-balloons/)
- [646. 最长数对链](https://leetcode.cn/problems/maximum-length-of-pair-chain/)
- [1288. 删除被覆盖区间](https://leetcode.cn/problems/remove-covered-intervals/) 1375
- [757. 设置交集大小至少为2](https://leetcode.cn/problems/set-intersection-size-at-least-two/) 2379
- [2589. 完成所有任务的最少时间](https://leetcode.cn/problems/minimum-time-to-complete-all-tasks/) 2381
另见 misc.go 中的 mergeIntervals 和 minJumpNumbers

### 中位数贪心（右边数字为难度分） // 注：算长度用左闭右开区间思考，算中间值用闭区间思考    两个中位数分别是 a[(n-1)/2] 和 a[n/2]
有两种证明方法，见 https://leetcode.cn/problems/5TxKeK/solution/zhuan-huan-zhong-wei-shu-tan-xin-dui-din-7r9b/
题单（右边数字为难度分）
- [462. 最小操作次数使数组元素相等 II](https://leetcode.cn/problems/minimum-moves-to-equal-array-elements-ii/)
- [2033. 获取单值网格的最小操作数](https://leetcode.cn/problems/minimum-operations-to-make-a-uni-value-grid/) 1672
- [2448. 使数组相等的最小开销](https://leetcode.cn/problems/minimum-cost-to-make-array-equal/) 2005
- [2607. 使子数组元素和相等](https://leetcode.cn/problems/make-k-subarray-sums-equal/) 2071
- [2967. 使数组成为等数数组的最小代价](https://leetcode.cn/problems/minimum-cost-to-make-array-equalindromic/) 2116
- [1478. 安排邮筒](https://leetcode.cn/problems/allocate-mailboxes/) 2190
- [2968. 执行操作使频率分数最大](https://leetcode.cn/problems/apply-operations-to-maximize-frequency-score/) 2444
- [1703. 得到连续 K 个 1 的最少相邻交换次数](https://leetcode.cn/problems/minimum-adjacent-swaps-for-k-consecutive-ones/) 2467
- [3086. 拾起 K 个 1 需要的最少行动次数](https://leetcode.cn/problems/minimum-moves-to-pick-k-ones/) 2673
- [LCP 24. 数字游戏](https://leetcode.cn/problems/5TxKeK/)
- [296. 最佳的碰头地点](https://leetcode.cn/problems/best-meeting-point/) 二维的情况（会员题）
https://codeforces.com/problemset/problem/710/B 1400
中位数相关 https://codeforces.com/problemset/problem/166/C 1500 *可以做到对不同的 x 用 O(log n) 时间回答

### 排序不等式
- [2285. 道路的最大总重要性](https://leetcode.cn/problems/maximum-total-importance-of-roads/) 1496
- [3016. 输入单词需要的最少按键次数 II](https://leetcode.cn/problems/minimum-number-of-pushes-to-type-word-ii/) 1534
- [1402. 做菜顺序](https://leetcode.cn/problems/reducing-dishes/) 1679
- [2931. 购买物品的最大开销](https://leetcode.cn/problems/maximum-spending-after-buying-items/) 1822
- [2809. 使数组和小于等于 x 的最少时间](https://leetcode.cn/problems/minimum-time-to-make-array-sum-at-most-x/) 2979
https://codeforces.com/problemset/problem/1165/E 1600

## 相邻不同
每次取两个数减一，最后剩下的数最小 / 操作次数最多 https://cs.stackexchange.com/a/145450
- [1753. 移除石子的最大得分](https://leetcode.cn/problems/maximum-score-from-removing-stones/) 1488
- [767. 重构字符串](https://leetcode.cn/problems/reorganize-string/) 1681
- [1054. 距离相等的条形码](https://leetcode.cn/problems/distant-barcodes/) 1702
- [1953. 你可以工作的最大周数](https://leetcode.cn/problems/maximum-number-of-weeks-for-which-you-can-work/) 1804
   - https://codeforces.com/problemset/problem/1579/D 1400
- [3139. 使数组中所有元素相等的最小开销](https://leetcode.cn/problems/minimum-cost-to-equalize-array/) 2666
- [621. 任务调度器](https://leetcode.cn/problems/task-scheduler/) 相同元素至少间隔 n
- [358. K 距离间隔重排字符串](https://leetcode.cn/problems/rearrange-string-k-distance-apart/)（会员题）
https://codeforces.com/problemset/problem/1521/E 2700 二维+对角不同

每次取数组中大于 0 的连续一段同时减 1，求使数组全为 0 的最少操作次数
https://leetcode.cn/problems/minimum-number-of-increments-on-subarrays-to-form-a-target-array/solutions/371326/xing-cheng-mu-biao-shu-zu-de-zi-shu-zu-zui-shao-ze/
https://codeforces.com/problemset/problem/448/C

邻项交换（最小代价排序/字典序最小） Exchange Arguments
https://codeforces.com/blog/entry/63533
某些题目和逆序对有关
- [1665. 完成所有任务的最少初始能量](https://leetcode.cn/problems/minimum-initial-energy-to-finish-tasks/) 1901
https://codeforces.com/problemset/problem/1638/B 1100
https://codeforces.com/problemset/problem/920/C 1400
https://codeforces.com/problemset/problem/435/B 1400
https://codeforces.com/contest/246/problem/A 900
https://atcoder.jp/contests/arc147/tasks/arc147_b
https://atcoder.jp/contests/abc268/tasks/abc268_f
相邻两数之差的绝对值为 1 https://ac.nowcoder.com/acm/contest/65259/C

非邻项交换（最小代价排序/字典序最小）
某些题目可以在 i 到 a[i] 之间连边建图
LC1202 https://leetcode.cn/problems/smallest-string-with-swaps/ 1855
LC2948 https://leetcode.cn/problems/make-lexicographically-smallest-array-by-swapping-elements/ 2047
https://codeforces.com/contest/252/problem/B
https://codeforces.com/problemset/problem/1768/D 1800
https://codeforces.com/contest/109/problem/D 2000
shift+reverse https://codeforces.com/contest/1907/problem/F

区间与点的最大匹配/覆盖问题
https://www.luogu.com.cn/problem/P2887
https://codeforces.com/problemset/problem/555/B
https://codeforces.com/problemset/problem/863/E

倒序
LC2718 https://leetcode.cn/problems/sum-of-matrix-after-queries/
- 加强版 https://www.luogu.com.cn/problem/P9715        ?contestId=126251

思维：观察、结论
- [2498. 青蛙过河 II](https://leetcode.cn/problems/frog-jump-ii/) 1759
- [782. 变为棋盘](https://leetcode.cn/problems/transform-to-chessboard/) 2430
https://codeforces.com/problemset/problem/1811/C 1100
https://codeforces.com/problemset/problem/1822/D 1200
https://codeforces.com/problemset/problem/1077/C 1300
https://codeforces.com/problemset/problem/1364/B 1300
https://codeforces.com/problemset/problem/1608/C 1700
https://codeforces.com/problemset/problem/1442/A 1800
https://codeforces.com/problemset/problem/558/C  1900
https://codeforces.com/problemset/problem/1744/F 2000
https://codeforces.com/problemset/problem/1610/E 2300

思维：脑筋急转弯
LC1503 https://leetcode.cn/problems/last-moment-before-all-ants-fall-out-of-a-plank/
LC2731 https://leetcode.cn/problems/movement-of-robots/
LC280 https://leetcode.cn/problems/wiggle-sort/
LC3012 https://leetcode.cn/problems/minimize-length-of-array-using-operations/
https://codeforces.com/problemset/problem/1009/B 1400
https://codeforces.com/problemset/problem/1883/F 1400
https://codeforces.com/problemset/problem/1169/B 1500
https://codeforces.com/problemset/problem/500/C 1600
https://codeforces.com/problemset/problem/601/A 1600
https://codeforces.com/problemset/problem/1763/C 2000
https://atcoder.jp/contests/abc194/tasks/abc194_e
https://atcoder.jp/contests/abc196/tasks/abc196_e
https://www.luogu.com.cn/problem/UVA10881 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=20&page=show_problem&problem=1822
- [LCS 01. 下载插件](https://leetcode.cn/problems/Ju9Xwi/)

注意值域
LC2653 https://leetcode.cn/problems/sliding-subarray-beauty/ 1786
LC2250 https://leetcode.cn/problems/count-number-of-rectangles-containing-each-point/ 1998
LC2857 https://leetcode.cn/problems/count-pairs-of-points-with-distance-k/ 2082
LC1906 https://leetcode.cn/problems/minimum-absolute-difference-queries/ 2147
LC1766 https://leetcode.cn/problems/tree-of-coprimes/ 2232
LC2198 https://leetcode.cn/problems/number-of-single-divisor-triplets/（会员题）

注意指数/对数
LC2188 https://leetcode.cn/problems/minimum-time-to-finish-the-race/ 2315
LC2920 https://leetcode.cn/problems/maximum-points-after-collecting-coins-from-all-nodes/ 2351

枚举答案
https://codeforces.com/contest/1977/problem/C

构造
题单 https://www.luogu.com.cn/training/14#problems
LC767 https://leetcode.cn/problems/reorganize-string/
LC667 https://leetcode.cn/problems/beautiful-arrangement-ii/
LC2745 https://leetcode.cn/problems/construct-the-longest-new-string/ 1607
LC2573 https://leetcode.cn/problems/find-the-string-with-lcp/ 2682
构造反例 https://leetcode.cn/problems/parallel-courses-iii/solution/tuo-bu-pai-xu-dong-tai-gui-hua-by-endles-dph6/2310439
构造 TLE 数据 https://leetcode.cn/problems/maximum-total-reward-using-operations-ii/solutions/2805413/bitset-you-hua-0-1-bei-bao-by-endlessche-m1xn/comments/2320111
https://codeforces.com/problemset/problem/1028/B  1200
https://codeforces.com/problemset/problem/1713/C  1200
https://codeforces.com/problemset/problem/1717/C  1300
https://codeforces.com/problemset/problem/1788/C  1300
https://codeforces.com/problemset/problem/1815/A  1300
https://codeforces.com/problemset/problem/1978/C  1300
https://codeforces.com/problemset/problem/803/A   1400
https://codeforces.com/problemset/problem/1838/C  1400
https://codeforces.com/problemset/problem/1863/D  1400
https://codeforces.com/problemset/problem/1896/C  1400
https://codeforces.com/problemset/problem/1630/A  1500
https://codeforces.com/problemset/problem/1710/A  1500
https://codeforces.com/problemset/problem/1722/G  1500
https://codeforces.com/problemset/problem/1809/C  1500
https://codeforces.com/problemset/problem/1968/E  1600
https://codeforces.com/problemset/problem/584/C   1700 分类讨论
https://codeforces.com/problemset/problem/1332/D  1700 给你一个错误代码，构造 hack 数据
https://codeforces.com/problemset/problem/142/B   1800 棋盘放最多的马
https://codeforces.com/problemset/problem/847/C   1800
https://codeforces.com/problemset/problem/1156/B  1800 相邻字母在字母表中不相邻
https://codeforces.com/problemset/problem/1267/L  1800
https://codeforces.com/problemset/problem/1304/D  1800 最短/最长 LIS
https://codeforces.com/problemset/problem/1554/D  1800
https://codeforces.com/problemset/problem/118/C   1900 贪心
https://codeforces.com/problemset/problem/327/D   1900
https://codeforces.com/problemset/problem/388/B   1900 两点间恰好 k 条最短路径
https://codeforces.com/problemset/problem/550/D   1900 度数均为 k 且至少（恰好）有一条割边
https://codeforces.com/problemset/problem/708/B   1900 分类讨论
https://codeforces.com/problemset/problem/1823/D  1900
https://codeforces.com/problemset/problem/1854/A2 1900 分类讨论
https://codeforces.com/problemset/problem/515/D   2000
https://codeforces.com/problemset/problem/1558/C  2000
https://codeforces.com/problemset/problem/1789/D  2200 推荐 位运算 把 X 变成 Y 不断靠近答案
https://codeforces.com/problemset/problem/1761/E  2400
https://codeforces.com/problemset/problem/1521/E  2700 二维相邻不同
https://codeforces.com/problemset/problem/1838/F  3000 交互 二分
https://atcoder.jp/contests/arc145/tasks/arc145_a
https://atcoder.jp/contests/agc015/tasks/agc015_d bit OR

不好想到的构造
https://codeforces.com/contest/1659/problem/D
https://atcoder.jp/contests/abc178/tasks/abc178_f
https://codeforces.com/problemset/problem/1689/E 脑筋急转弯
https://codeforces.com/problemset/problem/1787/E

不变量（想一想，操作不会改变什么）
https://codeforces.com/contest/1775/problem/E 有点差分的味道，想想前缀和
https://atcoder.jp/contests/arc119/tasks/arc119_c 操作不影响交错和
https://codeforces.com/problemset/problem/1365/F 仍然对称

不变量 2（总和）
把一个环形数组切两刀，分成两段，要求相等，求方案数 => 和为 sum(a)/2 的子数组个数
LC494 https://leetcode.cn/problems/target-sum/

行列独立 LC3189 https://leetcode.cn/problems/minimum-moves-to-get-a-peaceful-board/

分类讨论（部分题是易错题）
https://codeforces.com/problemset/problem/870/C 1300
https://codeforces.com/problemset/problem/1698/C 1300
https://codeforces.com/problemset/problem/30/A 1400
https://codeforces.com/problemset/problem/45/I 1400
https://codeforces.com/problemset/problem/489/C 1400
https://codeforces.com/problemset/problem/934/A 1400
https://codeforces.com/problemset/problem/1009/B 1400 脑筋急转弯
https://codeforces.com/problemset/problem/1251/B 1400
https://codeforces.com/problemset/problem/1292/A 1400 也有简单写法
https://codeforces.com/problemset/problem/1605/C 1400
https://codeforces.com/problemset/problem/960/B 1500
https://codeforces.com/problemset/problem/1051/C 1500
https://codeforces.com/problemset/problem/1180/B 1500
https://codeforces.com/problemset/problem/750/C 1600 *也有偏数学的做法
https://codeforces.com/problemset/problem/898/E 1600
https://codeforces.com/problemset/problem/1822/E 1600 样例给的挺良心的
https://codeforces.com/problemset/problem/1861/C 1600 好题！
https://codeforces.com/problemset/problem/1978/D 1600
https://codeforces.com/problemset/problem/193/A 1700
https://codeforces.com/problemset/problem/382/C 1700
https://codeforces.com/problemset/problem/411/C 1700
https://codeforces.com/problemset/problem/1516/C 1700
https://codeforces.com/problemset/problem/1799/C 1700
https://codeforces.com/problemset/problem/1833/G 1800 样例给的挺良心的
https://codeforces.com/problemset/problem/796/C 1900
https://codeforces.com/problemset/problem/1095/E 1900
https://codeforces.com/problemset/problem/1714/F 1900 锻炼代码实现技巧的好题
https://codeforces.com/problemset/problem/1914/F 1900
https://codeforces.com/problemset/problem/1763/C 2000
https://codeforces.com/problemset/problem/1978/E 2000
https://codeforces.com/problemset/problem/1811/F 2100
https://codeforces.com/problemset/problem/1798/E 2300
https://codeforces.com/problemset/problem/209/C 2400
https://codeforces.com/problemset/problem/1594/F 2400
https://codeforces.com/problemset/problem/1761/E 2400
https://codeforces.com/problemset/problem/1832/D2 2400
https://codeforces.com/problemset/problem/1730/E 2700
https://codeforces.com/gym/105139/problem/L
https://atcoder.jp/contests/diverta2019/tasks/diverta2019_c
https://atcoder.jp/contests/abc155/tasks/abc155_d
https://atcoder.jp/contests/abc125/tasks/abc125_d
https://atcoder.jp/contests/arc134/tasks/arc134_d 1998
- [335. 路径交叉](https://leetcode.cn/problems/self-crossing/)
- [2162. 设置时间的最少代价](https://leetcode.cn/problems/minimum-cost-to-set-cooking-time/) 1852
https://leetcode.cn/problems/maximize-the-number-of-partitions-after-operations/
https://leetcode.cn/problems/count-the-number-of-houses-at-a-certain-distance-ii/

大量分类讨论
- [420. 强密码检验器](https://leetcode.cn/problems/strong-password-checker/)
https://codeforces.com/problemset/problem/796/C 1900
https://codeforces.com/problemset/problem/1647/D 1900
https://codeforces.com/problemset/problem/356/C 2100
https://codeforces.com/problemset/problem/460/D 2300
https://codeforces.com/problemset/problem/1527/D 2400
https://codeforces.com/problemset/problem/1374/E2 2500
https://atcoder.jp/contests/arc153/tasks/arc153_c +构造
https://atcoder.jp/contests/agc015/tasks/agc015_d

贡献法
- [2063. 所有子字符串中的元音](https://leetcode.cn/problems/vowels-of-all-substrings/) 1663
LC979 https://leetcode.cn/problems/distribute-coins-in-binary-tree/ 1709
LC2477 https://leetcode.cn/problems/minimum-fuel-cost-to-report-to-the-capital/ 2012
LC891 https://leetcode.cn/problems/sum-of-subsequence-widths/
LC1588 https://leetcode.cn/problems/sum-of-all-odd-length-subarrays/
LC2681 https://leetcode.cn/problems/power-of-heroes/
- https://atcoder.jp/contests/arc116/tasks/arc116_b
LC2763 https://leetcode.cn/problems/sum-of-imbalance-numbers-of-all-subarrays/
更多贡献法题目，见 monotone_stack.go
https://codeforces.com/problemset/problem/1648/A 1400 维度独立
https://codeforces.com/problemset/problem/1691/C 1400
https://codeforces.com/problemset/problem/1789/C 1500 好题！
https://codeforces.com/problemset/problem/383/A 1600 好题
https://codeforces.com/problemset/problem/1165/E 1600
https://codeforces.com/problemset/problem/1715/C 1700 也可以用增量法
https://codeforces.com/problemset/problem/1777/D 1900 树
https://codeforces.com/problemset/problem/1788/D 2000 好题！
https://codeforces.com/problemset/problem/912/D 2100
https://codeforces.com/problemset/problem/1808/D 2100
https://codeforces.com/problemset/problem/1208/E 2200
https://codeforces.com/problemset/problem/749/E 2400
https://codeforces.com/problemset/problem/915/F 2400
https://atcoder.jp/contests/abc290/tasks/abc290_e 好题！
https://atcoder.jp/contests/abc159/tasks/abc159_f 与 0-1 背包结合
^+ https://atcoder.jp/contests/abc201/tasks/abc201_e
https://www.lanqiao.cn/problems/12467/learning/?contest_id=167

增量法
- [2262. 字符串的总引力](https://leetcode.cn/problems/total-appeal-of-a-string/) 2033
      结合线段树优化 DP https://codeforces.com/contest/833/problem/B 2200
- [828. 统计子串中的唯一字符](https://leetcode.cn/problems/count-unique-characters-of-all-substrings-of-a-given-string/) 2034
- [2916. 子数组不同元素数目的平方和 II](https://leetcode.cn/problems/subarrays-distinct-element-sum-of-squares-ii/) 2816
https://codeforces.com/problemset/problem/1715/C 1700 也可以用贡献法
https://codeforces.com/problemset/problem/1428/F 2400

小模拟
LC2534 https://leetcode.cn/problems/time-taken-to-cross-the-door/
https://atcoder.jp/contests/abc279/tasks/abc279_f

中模拟
https://atcoder.jp/contests/abc319/tasks/abc319_f

其他
删除一个字符 + 删除最长连续前缀 https://codeforces.com/problemset/problem/1430/D
https://codeforces.com/problemset/problem/521/D

先撤销，再恢复
LC3187 https://leetcode.cn/problems/peaks-in-array/

合法括号字符串 (Regular Bracket Sequence, RBS)
https://codeforces.com/problemset/problem/1097/C 1400
https://codeforces.com/problemset/problem/1837/D 1400
https://codeforces.com/problemset/problem/990/C 1500
https://codeforces.com/problemset/problem/847/C 1800 构造
https://codeforces.com/problemset/problem/1821/E 2100
https://codeforces.com/problemset/problem/1830/C 2400
https://codeforces.com/problemset/problem/3/D 2600 反悔贪心（反悔堆）

= 变成 <= 或者 >=
求前缀和/后缀和
https://leetcode.cn/problems/maximum-product-of-the-length-of-two-palindromic-substrings/

连续性 + 上下界
https://atcoder.jp/contests/arc137/tasks/arc137_b
https://codeforces.com/contest/1695/problem/C
*/

// 异类双变量：固定某变量统计另一变量的 [0,n)
//     EXTRA: 值域上的双变量，见 https://codeforces.com/contest/486/problem/D
// 同类双变量①：固定 i 统计 [0,n)
// 同类双变量②：固定 i 统计 [0,i-1]
// 套路：预处理数据（按照某种顺序排序/优先队列/BST/...），或者边遍历边维护，
//      然后固定变量 i，用均摊 O(1)~O(logn) 的复杂度统计范围内的另一变量 j
// 这样可以将复杂度从 O(n^2) 降低到 O(n) 或 O(nlogn)
// https://codeforces.com/problemset/problem/1194/E
// 进阶：https://codeforces.com/problemset/problem/1483/D
// 删除一段的最长连续递增 CERC10D https://codeforces.com/gym/101487
// 统计量是二元组的情形 https://codeforces.com/problemset/problem/301/D
// 好题 空间优化 https://codeforces.com/contest/1830/problem/B

// 双变量+下取整：枚举分母，然后枚举分子的范围，使得在该范围内的分子/分母是一个定值
// LC1862 https://leetcode.cn/problems/sum-of-floored-pairs/
// https://codeforces.com/problemset/problem/1706/D2

// 利用前缀和实现巧妙的构造 https://www.luogu.com.cn/blog/duyi/qian-zhui-he
// 邻项修改->前缀和->单项修改 https://codeforces.com/problemset/problem/1254/B2 https://ac.nowcoder.com/acm/contest/7612/C

/* 二进制枚举
https://www.luogu.com.cn/problem/UVA11464 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=26&page=show_problem&problem=2459
*/

/* 横看成岭侧成峰
转换为距离的众数 https://codeforces.com/problemset/problem/1365/C
转换为差分数组 https://codeforces.com/problemset/problem/1110/E
             https://codeforces.com/problemset/problem/1442/A
             https://codeforces.com/problemset/problem/1700/C
             https://codeforces.com/problemset/problem/1779/D 改成修改长为 x 的数组？
             https://www.luogu.com.cn/problem/P4552
转换为差 http://www.51nod.com/Challenge/Problem.html#problemId=1217
考虑每个点产生的贡献 https://codeforces.com/problemset/problem/1009/E
考虑每条边产生的负贡献 https://atcoder.jp/contests/abc173/tasks/abc173_f
考虑符合范围要求的贡献 https://codeforces.com/problemset/problem/1151/E
和式的另一视角。若每一项的值都在一个范围，不妨考虑另一个问题：值为 x 的项有多少个？https://atcoder.jp/contests/abc162/tasks/abc162_e
对所有排列考察所有子区间的性质，可以转换成对所有子区间考察所有排列。将子区间内部的排列和区间外部的排列进行区分，内部的性质单独研究，外部的当作 (n-(r-l))! 个排列 https://codeforces.com/problemset/problem/1284/C
从最大值入手 https://codeforces.com/problemset/problem/1381/B
等效性 LC1183 https://leetcode.cn/problems/maximum-number-of-ones/
LC1526 https://leetcode.cn/problems/minimum-number-of-increments-on-subarrays-to-form-a-target-array/
置换 https://atcoder.jp/contests/abc250/tasks/abc250_e
排序+最小操作次数 https://codeforces.com/contest/1367/problem/F2
https://codeforces.com/contest/1830/problem/A
从绝对值最大的开始思考 https://codeforces.com/contest/351/problem/E
https://codeforces.com/problemset/problem/777/C 1600

棋盘染色 LC2577 https://leetcode.cn/problems/minimum-time-to-visit-a-cell-in-a-grid/
        https://codeforces.com/contest/1848/problem/A

others https://codeforces.com/blog/entry/118706
*/

/*
## 练习：离线（按难度分排序）

> 由于所有的询问数据都给出了，我们可以通过修改询问的顺序，达到降低时间复杂度的效果。相应的，在线算法就是按照输入的顺序处理，来一个处理一个。

- [2343. 裁剪数字后查询第 K 小的数字](https://leetcode.cn/problems/query-kth-smallest-trimmed-number/) 1652
- [2070. 每一个查询的最大美丽值](https://leetcode.cn/problems/most-beautiful-item-for-each-query/) 1724
- [1847. 最近的房间](https://leetcode.cn/problems/closest-room/) 2082
- [2503. 矩阵查询可获得的最大分数](https://leetcode.cn/problems/maximum-number-of-points-from-grid-queries/) 2196
- [1851. 包含每个查询的最小区间](https://leetcode.cn/problems/minimum-interval-to-include-each-query/) 2286
- [1697. 检查边长度限制的路径是否存在](https://leetcode.cn/problems/checking-existence-of-edge-length-limited-paths/) 2300
- [2747. 统计没有收到请求的服务器数目](https://leetcode.cn/problems/count-zero-request-servers/) 2405
- [1938. 查询最大基因差](https://leetcode.cn/problems/maximum-genetic-difference-query/) 2503
- [2736. 最大和查询](https://leetcode.cn/problems/maximum-sum-queries/) 2533
*/

/* 逆向思维 / 正难则反
不可行方案通常比可行方案好求
- [2171. 拿出最少数目的魔法豆](https://leetcode.cn/problems/removing-minimum-number-of-magic-beans/) 1748
- [1354. 多次求和构造目标数组](https://leetcode.cn/problems/construct-target-array-with-multiple-sums/) 2015
LC803 https://leetcode.cn/problems/bricks-falling-when-hit/
LC936 https://leetcode.cn/problems/stamping-the-sequence/
LC1199 https://leetcode.cn/problems/minimum-time-to-build-blocks/
LC2382 https://leetcode.cn/problems/maximum-segment-sum-after-removals/
LCP52 https://leetcode.cn/problems/QO5KpG/
https://codeforces.com/problemset/problem/1792/C 1500
- 相似题目 https://codeforces.com/problemset/problem/1367/F1 2100
https://codeforces.com/problemset/problem/1882/B
https://codeforces.com/problemset/problem/712/C 1600
https://codeforces.com/problemset/problem/621/C 1700
https://codeforces.com/problemset/problem/1301/C 1700
https://codeforces.com/problemset/problem/1644/D 1700
https://codeforces.com/problemset/problem/1672/D 1700
https://codeforces.com/problemset/problem/1759/G 1900 求字典序最小，通常可以从大往小思考
https://codeforces.com/problemset/problem/1638/D 2000
https://codeforces.com/problemset/problem/571/A 2100
https://codeforces.com/problemset/problem/369/E 2200

删除变添加
https://codeforces.com/problemset/problem/295/B
https://leetcode.cn/problems/maximum-segment-sum-after-removals/
*/

/* 奇偶性
https://codeforces.com/problemset/problem/763/B
https://codeforces.com/problemset/problem/1270/E
https://codeforces.com/problemset/problem/1332/E 配对法：将合法局面与非法局面配对
LC932 https://leetcode.cn/problems/beautiful-array/ 分治
*/

/* 相邻 传递性
https://codeforces.com/problemset/problem/1582/E
*/

/* 归纳：solve(n)->solve(n-1) 或者 solve(n-1)->solve(n)
https://codeforces.com/problemset/problem/1517/C
https://codeforces.com/problemset/problem/412/D
https://codeforces.com/problemset/problem/266/C
*/

/* 见微知著：考察单个点的规律，从而推出全局规律
https://codeforces.com/problemset/problem/1510/K
LC1806 https://leetcode.cn/problems/minimum-number-of-operations-to-reinitialize-a-permutation/ 1491
*/

// 「恰好」转换成「至少/至多」https://codeforces.com/problemset/problem/1188/C

/* 反悔贪心
另见 heap.go 中的「反悔堆」
https://djy-juruo.blog.luogu.org/qian-tan-fan-hui-tan-xin
https://www.jvruo.com/archives/1844/
https://www.cnblogs.com/nth-element/p/11768155.html
题单 https://www.luogu.com.cn/training/8793
LC1388 双向链表反悔贪心 https://leetcode.cn/problems/pizza-with-3n-slices/
LC2813 https://leetcode.cn/problems/maximum-elegance-of-a-k-length-subsequence/
*/

/* 集合哈希
https://codeforces.com/problemset/problem/1394/B
https://www.luogu.com.cn/problem/P6688
*/

/* 操作树
和莫队类似，通过改变查询的顺序来优化复杂度
https://codeforces.com/problemset/problem/707/D
*/

/* Golang 卡常技巧（注：关于 IO 的加速见 io.go）
对于存在海量小对象的情况（如 trie, treap 等），使用 debug.SetGCPercent(-1) 来禁用 GC，能明显减少耗时
对于可以回收的情况（如 append 在超过 cap 时），使用 debug.SetGCPercent(-1) 虽然会减少些许耗时，但若有大量内存没被回收，会有 MLE 的风险
其他情况下使用 debug.SetGCPercent(-1) 对耗时和内存使用无明显影响
对于多组数据的情况，若禁用 GC 会 MLE，可在每组数据的开头或末尾调用 runtime.GC() 或 debug.FreeOSMemory() 手动 GC
参考 https://draveness.me/golang/docs/part3-runtime/ch07-memory/golang-garbage-collector/
    https://zhuanlan.zhihu.com/p/77943973

128MB ~1e7 个 int64
256MB ~3e7 个 int64
512MB ~6e7 个 int64
1GB   ~1e8 个 int64

如果没有禁用 GC 但 MLE，可以尝试 1.19 新增的 debug.SetMemoryLimit
例如 debug.SetMemoryLimit(200<<20)，其中 200 可以根据题目的约束来修改
具体见如下测试：
180<<20 1996ms 255100KB https://codeforces.com/contest/1800/submission/203769679
195<<20  779ms 257800KB https://codeforces.com/contest/1800/submission/203768086
200<<20  654ms 259300KB https://codeforces.com/contest/1800/submission/203768768
205<<20  764ms 220100KB https://codeforces.com/contest/1800/submission/203771041
210<<20        MLE
参考 https://go.dev/doc/gc-guide#Memory_limit

对于二维矩阵，以 make([][mx]int, n) 的方式使用，比 make([][]int, n) 嵌套 make([]int, m) 更高效（100MB 以上时可以快 ~150ms）
但需要注意这种方式可能会向 OS 额外申请一倍的内存
对比 https://codeforces.com/problemset/submission/375/118043978
    https://codeforces.com/problemset/submission/375/118044262

函数内的递归 lambda 会额外消耗非常多的内存（100~200MB / 1e6 递归深度）
写在 main 里面 + slice MLE      https://codeforces.com/contest/767/submission/174193385
写在 main 外面 + slice 188364KB https://codeforces.com/contest/767/submission/174194380
附：
写在 main 里面 + array 257424KB https://codeforces.com/contest/767/submission/174194515
写在 main 外面 + array 154500KB https://codeforces.com/contest/767/submission/174193693

在特殊情况下，改为手动模拟栈可以减少 > 100MB 的内存
见这题的 Go 提交记录 https://codeforces.com/problemset/problem/163/E

测试：哈希表用时是数组的 13 倍（本题瓶颈）
slice    249ms https://codeforces.com/problemset/submission/570/209063267
hashmap 3259ms https://codeforces.com/problemset/submission/570/209063603
*/

func main() {
	time1 := time.Now()
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	_ = arr
	for i := 0; i < 1e8; i++ {
		// cast[[]int64](arr)
		// bool2int(true)
		cast[int](true)
	}
	fmt.Println(time.Since(time1))
	intSliceAsMapKeyExample(map[string]int{}, []int{1, 2, 3})
}

func cast[To, From any](v From) To {
	return *(*To)(unsafe.Pointer(&v))
}

// bool2int returns 0 if x is false or 1 if x is true.
func bool2int(x bool) int {
	return int(*(*uint8)(unsafe.Pointer(&x)))
}

// slice 作为 map 的 key
// 长度为 0 的 slice 对应空字符串
func intSliceAsMapKeyExample(cnt map[string]int, a []int) {
	// 如果后面还会修改 a，可以先 copy 一份
	//a = append(a[:0:0], a...)
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	sh.Len *= bits.UintSize / 8 // 装作 byte slice
	s := *(*string)(unsafe.Pointer(sh))
	fmt.Println(s)
	cnt[s]++
}

// 力扣上的 int 和 int64 是一样的，但是有些题目要求返回 []int64
// !此时可以用指针强转
func intsToInt64s(a []int) []int64 {
	int64s := *(*[]int64)(unsafe.Pointer(&a))
	return int64s
}
