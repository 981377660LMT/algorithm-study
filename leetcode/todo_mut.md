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
