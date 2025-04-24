- [P ≠ NP 予想について 〜 NP、NP 完全、NP 困難を整理 〜](https://qiita.com/drken/items/5187e49082f7437349c2)
- [ソートを極める！ 〜 なぜソートを学ぶのか 〜](https://qiita.com/drken/items/44c60118ab3703f7727f)

- [百花繚乱！なないろ言語で競技プログラミングをする資料まとめ](https://qiita.com/drken/items/6edb1c0542d4c3b7179c#rust)

- [計算量オーダーの求め方を総整理！ 〜 どこから log が出て来るか 〜](https://qiita.com/drken/items/872ebc3a2b5caaa4a0d0)

- [AtCoder に登録したら次にやること ～ これだけ解けば十分闘える！過去問精選 10 問 ～](https://qiita.com/drken/items/fd4e5e3630d0f5859067)

- [二分探索アルゴリズムを一般化 〜 めぐる式二分探索法のススメ 〜](https://qiita.com/drken/items/97e37dd6143e33a64c8c)
  https://x.com/meguru_comp/status/697008509376835584/photo/4
- [k 番目の値を高速に取り出せるデータ構造のまとめ - BIT上二分探索や平衡二分探索木など](https://qiita.com/drken/items/1b7e6e459c24a83bb7fd)

- [ビット演算 (bit 演算) の使い方を総特集！ 〜 マスクビットから bit DP まで 〜](https://qiita.com/drken/items/7c6ff2aa4d8fce1c9361)

  - XorShift：使用位运算的简单随机数生成方法。随机数质量高且超高速

  ```cpp
  unsigned int randInt() {
      static unsigned int tx = 123456789, ty=362436069, tz=521288629, tw=88675123;
      unsigned int tt = (tx^(tx<<11));
      tx = ty; ty = tz; tz = tw;
      return ( tw=(tw^(tw>>19))^(tt^(tt>>8)) );
  }
  ```

  - 提取最右边的1：`nbit = bit & -bit`
  - 将最右边的1变成0：`nbit = bit & (bit - 1)`
  - O(3^n) frameWork：子集的子集dp

    ```go
    // n: 集合大小
    n := 4
    N := 1 << n // 总状态数

    // dp[i] 表示状态i的最优解
    dp := make([]int, N)
    // 初始化 dp...

    for i := 0; i < N; i++ {
        // 枚举 i 的所有非空子集 sub
        for sub := i; sub > 0; sub = (sub - 1) & i {
            // 在这里写转移逻辑
            // 例如：dp[i] = min(dp[i], dp[sub] + dp[i^sub])
        }
        // 如需处理 sub=0 的情况，可单独处理
    }
    ```
