不要删!!!

- 52
- 68
- 88
  https://leetcode.cn/problems/merge-sorted-array/solutions/2385610/dao-xu-shuang-zhi-zhen-wei-shi-yao-dao-x-xxkp/

- 90
  https://leetcode.cn/problems/subsets-ii/solutions/3036436/liang-chong-fang-fa-xuan-huo-bu-xuan-mei-v0js/

- 99
  Morris

- 172
  https://leetcode.cn/problems/factorial-trailing-zeroes/solutions/2972637/yan-jin-shu-xue-zheng-ming-pythonjavaccg-fe5t/

- 186
  https://leetcode.cn/problems/reverse-words-in-a-string-ii/solutions/2416792/fan-zhuan-zi-fu-chuan-zhong-de-dan-ci-ii-wzli/

- 229
- 超级简短的分治加归并的写法
  https://leetcode.cn/problems/the-skyline-problem/solutions/2010141/by-mei-mei-16-6jtg/

- 258

https://leetcode.cn/problems/add-digits/solutions/3008465/o1-zuo-fa-jie-shi-qing-chu-9-shi-zen-yao-wmoo/

- 288
- 291
- 296. 最佳的碰头地点
       https://leetcode.cn/problems/best-meeting-point/solutions/2381941/zui-jia-de-peng-tou-di-dian-by-leetcode-folxq/
- 302
  https://leetcode.cn/problems/smallest-rectangle-enclosing-black-pixels/description/
- 313. 超级丑数
       https://leetcode.cn/problems/super-ugly-number/description/

- 317
  法三
  https://leetcode.cn/problems/shortest-distance-from-all-buildings/description/

- 318 优化二
  https://leetcode.cn/problems/maximum-product-of-word-lengths/solutions/1104441/zui-da-dan-ci-chang-du-cheng-ji-by-leetc-lym9/
- 326. 3 的幂
       https://leetcode.cn/problems/power-of-three/solutions/2974674/o1-shu-xue-zuo-fa-yi-xing-gao-ding-pytho-w0uh/
- 330
  https://leetcode.cn/problems/patching-array/solutions/2551840/yong-gui-na-fa-zheng-ming-pythonjavacgo-mvyu1/
- 334. 递增的三元子序列
       https://leetcode.cn/problems/increasing-triplet-subsequence/description/
- https://leetcode.cn/problems/counting-bits/solutions/627418/bi-te-wei-ji-shu-by-leetcode-solution-0t1i/

- 350
  https://leetcode.cn/problems/intersection-of-two-arrays-ii/solutions/3056041/yi-ci-bian-li-jian-ji-xie-fa-fu-jin-jie-szdb1/

- 358
  https://leetcode.cn/problems/rearrange-string-k-distance-apart/solutions/1882750/by-424479543-1bwj/

- 366
  https://leetcode.cn/problems/find-leaves-of-binary-tree/

- 386
  https://leetcode.cn/problems/lexicographical-numbers/solutions/1428281/zi-dian-xu-pai-shu-by-leetcode-solution-98mz/

- 391
  https://leetcode.cn/problems/perfect-rectangle/description/

```C++
class Solution {
public:
    bool isRectangleCover(vector<vector<int>>& rectangles) {
        map<pair<int, int>, int> mp;
        for (auto &item : rectangles) {
            int x = item[0], y = item[1], a = item[2], b = item[3];
            if (++mp[{x, y}] == 0) mp.erase({x, y});
            if (--mp[{x, b}] == 0) mp.erase({x, b});
            if (++mp[{a, b}] == 0) mp.erase({a, b});
            if (--mp[{a, y}] == 0) mp.erase({a, y});
        }
        return mp.size() == 4 && mp.begin()->second == 1;
    }
};
```

- 395
  https://leetcode.cn/problems/longest-substring-with-at-least-k-repeating-characters/description/
- 408. 有效单词缩写
       https://leetcode.cn/problems/valid-word-abbreviation/description/

- 418. 屏幕可显示句子的数量

https://leetcode.cn/problems/sentence-screen-fitting/description/

- 424. 替换后的最长重复字符

https://leetcode.cn/problems/longest-repeating-character-replacement/description/

- 425. 单词方块
       https://leetcode.cn/problems/word-squares/

- 467
  https://leetcode.cn/problems/unique-substrings-in-wraparound-string/description/
- 471. 编码最短长度的字符串
       gpt
       https://leetcode.cn/problems/encode-string-with-shortest-length/description/
- 472. 连接词
       https://leetcode.cn/problems/concatenated-words/submissions/252968643/

- 482. 密钥格式化
       https://leetcode.cn/problems/license-key-formatting/description/

- 492. 构造矩形
       https://leetcode.cn/problems/construct-the-rectangle/description/

- 527. 单词缩写

https://leetcode.cn/problems/word-abbreviation/solutions/2329140/an-dan-ci-chang-du-yu-mo-wei-zi-mu-fen-z-d93s/

- 540. 有序数组中的单一元素
       https://leetcode.cn/problems/single-element-in-a-sorted-array/solutions/2983333/er-fen-xing-zhi-fen-xi-jian-ji-xie-fa-py-0rng/

- 611
  https://leetcode.cn/problems/valid-triangle-number/description/
- 616. 给字符串添加加粗标签
       https://leetcode.cn/problems/add-bold-tag-in-string/description/
- 632. 最小区间
       滑窗 https://leetcode.cn/problems/smallest-range-covering-elements-from-k-lists/solutions/2982588/liang-chong-fang-fa-dui-pai-xu-hua-dong-luih5/
- 644. 子数组最大平均数 II
       https://leetcode.cn/problems/maximum-average-subarray-ii/solutions/860135/fu-za-du-wei-onde-dan-diao-zhan-fa-by-li-trzz/
- 650
  https://leetcode.cn/problems/2-keys-keyboard/description/
- 654. 最大二叉树
       https://leetcode.cn/problems/maximum-binary-tree/description/
- 665. 非递减数列
       https://leetcode.cn/problems/non-decreasing-array/solutions/594758/fei-di-jian-shu-lie-by-leetcode-solution-zdsm/

- 672. 灯泡开关 Ⅱ

线性基
https://leetcode.cn/problems/bulb-switcher-ii/description/

- 673. 最长递增子序列的个数
       https://leetcode.cn/problems/number-of-longest-increasing-subsequence/solutions/1007075/zui-chang-di-zeng-zi-xu-lie-de-ge-shu-by-w12f/
- 679. 24 点游戏
       https://leetcode.cn/problems/24-game/description/
- 680. 验证回文串 II
       https://leetcode.cn/problems/valid-palindrome-ii/solutions/3053249/tan-xin-wei-shi-yao-xiang-deng-de-shi-ho-wtll/
- 683. K 个关闭的灯泡
       https://leetcode.cn/problems/k-empty-slots/description/
       滑窗、单调队列解法
- 685. 冗余连接 II

有向基环树删边成树
https://leetcode.cn/problems/redundant-connection-ii/description/

- 686. 重复叠加字符串匹配
- 687. 最长同值路径
       https://leetcode.cn/problems/longest-univalue-path/solutions/2227160/shi-pin-che-di-zhang-wo-zhi-jing-dpcong-524j4/
- 696. 计数二进制子串
- 702. 搜索长度未知的有序数组
       https://leetcode.cn/problems/search-in-a-sorted-array-of-unknown-size/description/
- 710. 黑名单中的随机数

https://leetcode.cn/problems/random-pick-with-blacklist/solutions/1626454/by-xie-dai-ma-de-huo-che-a31p/

- 711. 不同岛屿的数量 II
       https://leetcode.cn/problems/number-of-distinct-islands/description/
       https://leetcode.cn/problems/number-of-distinct-islands-ii/description/
- 712. 两个字符串的最小ASCII删除和
       https://leetcode.cn/problems/minimum-ascii-delete-sum-for-two-strings/
- https://leetcode.cn/problems/binary-number-with-alternating-bits/description/
- 722. 删除注释
       https://leetcode.cn/problems/remove-comments/description/
       正则
- 754. 到达终点数字
       https://leetcode.cn/problems/reach-a-number/solutions/1947254/fen-lei-tao-lun-xiang-xi-zheng-ming-jian-sqj2/
- 756. 金字塔转换矩阵
       https://leetcode.cn/problems/pyramid-transition-matrix/description/
- 766.  托普利茨矩阵
        起点+方向
        https://leetcode.cn/problems/toeplitz-matrix/solutions/613732/tuo-pu-li-ci-ju-zhen-by-leetcode-solutio-57bb/

              进阶：

              如果矩阵存储在磁盘上，并且内存有限，以至于一次最多只能将矩阵的一行加载到内存中，该怎么办？
              如果矩阵太大，以至于一次只能将不完整的一行加载到内存中，该怎么办？

  对于进阶问题一，一次最多只能将矩阵的一行加载到内存中，我们将每一行复制到一个连续数组中，随后在读取下一行时，就与内存中此前保存的数组进行比较。

        对于进阶问题二，一次只能将不完整的一行加载到内存中，我们将整个矩阵竖直切分成若干子矩阵，并保证两个相邻的矩阵至少有一列或一行是重合的，然后判断每个子矩阵是否符合要求。

- 777. 在 LR 字符串中交换相邻字符
       https://leetcode.cn/problems/swap-adjacent-in-lr-string/description/
- 795. 区间子数组个数
       https://leetcode.cn/problems/number-of-subarrays-with-bounded-maximum/solutions/1988198/tu-jie-yi-ci-bian-li-jian-ji-xie-fa-pyth-n75l/

       1. 统计定界子数组的数目
          https://leetcode.cn/problems/count-subarrays-with-fixed-bounds/

- 827. 最大人工岛
       https://leetcode.cn/problems/making-a-large-island/solutions/2808887/jian-ji-gao-xiao-ji-suan-dao-yu-de-mian-ab4h7/
- 833. 字符串中的查找与替换
       https://leetcode.cn/problems/find-and-replace-in-string/description/
- 835. 图像重叠

https://leetcode.cn/problems/image-overlap/solutions/527350/ni-ke-neng-wu-fa-xiang-xiang-de-on2lognd-gc5j/

- 862. 和至少为 K 的最短子数组
       https://leetcode.cn/problems/shortest-subarray-with-sum-at-least-k/solutions/1925036/liang-zhang-tu-miao-dong-dan-diao-dui-li-9fvh/

- 917. 仅仅反转字母
       https://leetcode.cn/problems/reverse-only-letters/description/
- 939. 最小面积矩形
       https://leetcode.cn/problems/minimum-area-rectangle/description/
- 963. 最小面积矩形 II
       https://leetcode.cn/problems/minimum-area-rectangle-ii/solutions/707666/c-on2-0ms-100-by-hqztrue-9ij7/
       最坏情况下n个点可以组成Θ(n2logn)个矩形

- 1004. 最大连续1的个数 III
        https://leetcode.cn/problems/max-consecutive-ones-iii/solutions/2126631/hua-dong-chuang-kou-yi-ge-shi-pin-jiang-yowmi/

TODO: 优化fix

- 1092. 最短公共超序列
        https://leetcode.cn/problems/shortest-common-supersequence/solutions/2194615/cong-di-gui-dao-di-tui-jiao-ni-yi-bu-bu-auy8z/
- 1105
  https://leetcode.cn/problems/filling-bookcase-shelves/
  排版

- https://leetcode.cn/problems/divide-array-into-increasing-sequences/description/
- 1131. 绝对值表达式的最大值

- 1203. 项目管理
        https://leetcode.cn/problems/sort-items-by-groups-respecting-dependencies/description/
- 1213
  https://leetcode.cn/problems/intersection-of-three-sorted-arrays/description/

- https://leetcode.cn/problems/missing-number-in-arithmetic-progression/description/
- https://leetcode.cn/problems/find-positive-integer-solution-for-a-given-equation/solutions/2117698/xiang-xiang-shuang-zhi-zhen-yi-ge-shi-pi-nr4y/
  `inspect.getsource(customfunction.__class__)`
  python inspect
- 1265. 逆序打印不可变链表
        https://leetcode.cn/problems/print-immutable-linked-list-in-reverse/solutions/1075083/fen-er-zhi-zhi-zhen-zheng-de-onshi-jian-5ljbe/
- 1326. 灌溉花园的最少水龙头数目
        https://leetcode.cn/problems/minimum-number-of-taps-to-open-to-water-a-garden/solutions/2123855/yi-zhang-tu-miao-dong-pythonjavacgo-by-e-wqry/
        maxJump模板
- 1329. 将矩阵按对角线排序
        原地排序
        https://leetcode.cn/problems/sort-the-matrix-diagonally/solutions/2760094/dui-jiao-xian-pai-xu-fu-yuan-di-pai-xu-p-uts8/
- 1401
  https://leetcode.cn/problems/circle-and-rectangle-overlapping/solutions/2319756/xi-jie-tai-duo-san-fen-tao-san-fen-bing-o6f3l/
- 1409. 查询带键的排列
        树状数组排队问题
        https://leetcode.cn/problems/queries-on-a-permutation-with-key/solutions/228032/cha-xun-dai-jian-de-pai-lie-by-leetcode-solution/
- https://leetcode.cn/problems/html-entity-parser/description/
- https://leetcode.cn/problems/making-file-names-unique/description/
- https://leetcode.cn/problems/clone-n-ary-tree/description/
- 1534. 统计好三元组
        https://leetcode.cn/problems/count-good-triplets/solutions/3622921/liang-chong-fang-fa-bao-li-mei-ju-qian-z-apcv/
- https://leetcode.cn/problems/dot-product-of-two-sparse-vectors/description/
- 1585. 检查字符串是否可以通过排序子字符串得到另一个字符串
        https://leetcode.cn/problems/check-if-string-is-transformable-with-substring-sort-operations/solutions/412180/20-xing-jian-dan-zuo-fa-by-zerotrac2/
- 1612. 检查两棵二叉表达式树是否等价
        进阶：当你的答案需同时支持 '-' 运算符（减法）时，你该如何修改你的答案
        https://leetcode.cn/problems/check-if-two-expression-trees-are-equivalent/solutions/453074/ji-lu-mei-yi-bian-liang-de-ge-shu-zhi-chi-jia-jian/
- 1618. 找出适应屏幕的最大字号
        api设计
        https://leetcode.cn/problems/maximum-font-to-fit-a-sentence-in-a-screen/description/
- 1638. 统计只差一个字符的子串数目
        https://leetcode.cn/problems/count-substrings-that-differ-by-one-character/solutions/2192600/tu-jie-fei-bao-li-onm-suan-fa-pythonjava-k5og/
        相似题目 795、2444
- 1778. 未知网格中的最短路径
        https://leetcode.cn/problems/shortest-path-in-a-hidden-grid/description/
- 1810. 隐藏网格下的最小消耗路径
        https://leetcode.cn/problems/minimum-path-cost-in-a-hidden-grid/description/
- 1868. 两个行程编码数组的积
        https://leetcode.cn/problems/product-of-two-run-length-encoded-arrays/description/
- 1989. 捉迷藏中可捕获的最大人数
        https://leetcode.cn/problems/maximum-number-of-people-that-can-be-caught-in-tag/solutions/1298212/go-shuang-zhi-zhen-fu-za-du-on-by-while-50reg/
- 2052. 将句子分隔成行的最低成本
        **cost modal** 的应用
        https://leetcode.cn/problems/minimum-cost-to-separate-sentence-into-rows/description/

- 2248
  https://leetcode.cn/problems/intersection-of-multiple-arrays/description/
- 非均摊的O(1)做法
  https://leetcode.cn/problems/longest-uploaded-prefix/solutions/1869943/fei-jun-tan-de-o-by-mei-mei-16-bvij/
- 2532. 过桥的时间
        https://leetcode.cn/problems/time-to-cross-a-bridge/solutions/2140586/shi-yong-sheng-cheng-qi-mo-ni-fu-za-guo-yxbnp/

---

多线程(go)

- https://leetcode.cn/problems/print-foobar-alternately/?envType=problem-list-v2&envId=concurrency
- https://leetcode.cn/problems/web-crawler-multithreaded/?envType=problem-list-v2&envId=concurrency

---

蛙佬

- 1240. 铺瓷砖
        这个问题的研究状况总结
        https://leetcode.cn/problems/tiling-a-rectangle-with-the-fewest-squares/solutions/2316994/zhe-ge-wen-ti-de-yan-jiu-zhuang-kuang-zo-b7jv/

- 875. 爱吃香蕉的珂珂
       https://leetcode.cn/problems/koko-eating-bananas/solutions/1539347/by-hqztrue-3gs3/

- Alias method O(1) pick
  https://leetcode.cn/problems/random-point-in-non-overlapping-rectangles/solutions/1543639/o1-pickde-suan-fa-by-hqztrue-vupd/

- 2952. 需要添加的硬币的最小数量
        分桶线性
        https://leetcode.cn/problems/minimum-number-of-coins-to-be-added/solutions/3052162/xian-xing-suan-fa-by-hqztrue-rf5n/
        https://leetcode.cn/problems/destroying-asteroids/solutions/1187851/xian-xing-shi-jian-suan-fa-by-hqztrue-vr3u/
- 3399. 字符相同的最短子字符串 II
        https://leetcode.cn/problems/smallest-substring-with-identical-characters-ii/solutions/3027426/xian-xing-suan-fa-by-hqztrue-kraj/

- 3097. 或值至少为 K 的最短子数组 II
        slidingWindowAggregation
        https://leetcode.cn/problems/shortest-subarray-with-or-at-least-k-ii/description/

- 力扣评测流程
  https://leetcode.cn/problems/serialize-and-deserialize-binary-tree/solutions/2125871/tou-gong-jian-liao-de-zuo-fa-li-yong-pin-gd7e/

- **线性分治**
  https://leetcode.cn/problems/house-robber-iv/solutions/2094194/geng-kuai-de-onzuo-fa-by-hqztrue-f85k/
  https://leetcode.cn/problems/house-robber-iv/solutions/2097454/wan-quan-onde-zuo-fa-fu-xian-wei-chang-s-96h7/
- 中位数分治
  https://leetcode.cn/problems/minimum-operations-to-halve-array-sum/solutions/1352188/onsuan-fa-by-hqztrue-jalf/
  https://leetcode.cn/problems/maximal-score-after-applying-k-operations/solutions/2052490/liang-chong-xian-xing-zuo-fa-by-hqztrue-q2wb/
- 线性选择
  https://leetcode.cn/problems/put-marbles-in-bags/solutions/2081183/onxian-xing-xuan-ze-by-hqztrue-ykur/

  np.partition

- 跑汇编
  https://leetcode.cn/problems/count-increasing-quadruplets/solutions/2164712/chao-yue-cde-su-du-by-hqztrue-dx7t/
  https://leetcode.cn/discuss/post/3447530/

- 524. 通过删除字母匹配到字典里最长单词
       https://leetcode.cn/problems/longest-word-in-dictionary-through-deleting/solutions/1010188/xian-xing-suan-fa-16ms-100-by-hqztrue-71co/
- 更快的 C(n,k) 复杂度的 dfs
  https://leetcode.cn/problems/maximum-rows-covered-by-columns/solutions/1815891/by-hqztrue-ny27/
- 有漏洞的哈希函数
  https://leetcode.cn/problems/delete-duplicate-folders-in-system/solutions/904603/he-ji-yi-xie-you-lou-dong-de-ha-xi-han-s-nylg/
- 868. 二进制间距
       https://leetcode.cn/problems/binary-gap/solutions/837523/wei-yun-suan-olog-log-n-by-hqztrue-aflh/
- median finding
  https://leetcode.cn/problems/maximum-ice-cream-bars/solutions/732895/c-on-80ms-100-by-hqztrue-x357/

---

灵剑

- 使用A\*优化运行速度
  https://leetcode.cn/problems/cherry-pickup/solutions/453318/shi-yong-ayou-hua-yun-xing-su-du-by-ling-jian-2012/

  - 值域有限时，使用多个stack代替优先队列(TODO: 封装一下)

- 1203. 项目管理
        利用扩展的拓扑排序一次性求出结果

  https://leetcode.cn/problems/sort-items-by-groups-respecting-dependencies/solutions/467233/li-yong-kuo-zhan-de-tuo-bu-pai-xu-yi-ci-xing-qiu-c/

- 1235. 规划兼职工作
        https://leetcode.cn/problems/maximum-profit-in-job-scheduling/solutions/467842/fei-chang-hao-dong-de-dong-tai-gui-hua-shi-xian-by/

- https://leetcode.cn/problems/maximum-number-of-non-overlapping-substrings/solutions/474713/jian-dan-ming-kuai-de-onsuan-fa-by-ling-jian-2012/
  合法的有效子字符串之间要么互不重叠，要么一个包含另一个，这样的结构很适合用栈来表示

- https://leetcode.cn/problems/image-overlap/solutions/527350/ni-ke-neng-wu-fa-xiang-xiang-de-on2lognd-gc5j/
- https://leetcode.cn/problems/count-good-meals/solutions/863540/yao-shi-yao-logcshuo-onjiu-shi-zhen-on-b-pwtn/

- 超级复杂的O(N + P)算法——后缀树和O(1) LCA
  先暴力
  https://leetcode.cn/problems/pattern-matching-lcci/solutions/875256/chao-ji-fu-za-de-on-psuan-fa-hou-zhui-sh-5q5h/

- https://leetcode.cn/problems/amount-of-new-area-painted-each-day/solutions/1542333/li-yong-dui-de-by-ling-jian-2012-fooy/

- 永远查询的区间都是[i - k, i)形式的
  https://leetcode.cn/problems/longest-increasing-subsequence-ii/solutions/2022560/shu-zhuang-shu-zu-onlogk-by-ling-jian-20-iqzq/

- 2589. 完成所有任务的最少时间

- Main–Lorentz
  https://leetcode.cn/problems/count-beautiful-splits-in-an-array/solutions/3023632/main-lorentzcha-zhao-zhong-fu-zi-chuan-f-nwsc/
  https://leetcode.cn/problems/count-beautiful-splits-in-an-array/solutions/3020684/onlogn-hou-zhui-shu-zu-jie-fa-by-vclip-suiq/
  https://leetcode.cn/problems/distinct-echo-substrings/

---

白

- O(n) 桶解法
  https://leetcode.cn/problems/smallest-substring-with-identical-characters-ii/solutions/3027075/on-tong-jie-fa-by-vclip-oneq/
  https://leetcode.cn/problems/smallest-substring-with-identical-characters-ii/solutions/3027426/xian-xing-suan-fa-by-hqztrue-kraj/

- O(√qnlogn) 的二分查找+缓存查询结果解法
  https://leetcode.cn/problems/shortest-word-distance-ii/solutions/1941801/zheng-que-de-o-by-vclip-ypg6/

- O(nlog^2n)解法 二维树状数组 cdq分治
  https://leetcode.cn/problems/maximum-height-by-stacking-cuboids/solutions/2014599/by-vclip-gmet/

- 字符串更长时的O(nm)解法 后缀自动机优化dp
  https://leetcode.cn/problems/longest-string-chain/solutions/2247668/zi-fu-chuan-geng-chang-shi-de-onmjie-fa-bjk1b/

- 1397. 找到所有好字符串
        更快的 O(nmlogm) 解法
        https://leetcode.cn/problems/find-all-good-strings/solutions/2429620/geng-kuai-de-onmlogm-jie-fa-by-vclip-m7in/

- 3076. 数组中的最短非公共子字符串
        线性的后缀数组解法
        https://leetcode.cn/problems/shortest-uncommon-substring-in-an-array/solutions/2678131/xian-xing-de-hou-zhui-shu-zu-jie-fa-by-v-g7bc/
- O(logn)解决任意鸡蛋个数的问题
  https://leetcode.cn/problems/egg-drop-with-2-eggs-and-n-floors/solutions/2948790/olognjie-jue-ren-yi-ji-dan-ge-shu-de-wen-yjlq/

- O(n) 桶解法
  https://leetcode.cn/problems/smallest-substring-with-identical-characters-ii/solutions/3027075/on-tong-jie-fa-by-vclip-oneq/
