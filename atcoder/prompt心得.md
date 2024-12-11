prompt 心得

1. A-D(简单题) 可以直接给描述、数据范围、样例
2. 要求写暴力，然后给定数据范围让 gpt 优化
3. 输入 gpt 生成的代码，要 gpt 修复
4. o1-preview 不行就换 o1-mini
5. 最好都让 gpt 写，不要加自己的代码
6. 确保代码通过所有测试用例
7. 没过样例时，回复哪个样例不对，持续回复"不对"，需要一直提示错误的样例
8. 样例通过了，但是始终 WA 几个点，考虑更换模型。
9. 从问题文开始复制到样例结束
10. 切换语言？

    - 翻译成 rust，使用 proconio 读取
    - 翻译成 C++ 20 (gcc 12.2)，并增加头

    ```CPP
    #pragma GCC target("arch=skylake-avx512")
    #pragma GCC target("avx2")
    #pragma GCC optimize("O3")
    #pragma GCC target("sse4")
    #pragma GCC optimize("unroll-loops")
    #pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")

    #include <bits/stdc++.h>
    using namespace std;

    int main() {
      int N, Q;
      cin >> N >> Q;
      string S;
      cin >> S;
      vector< int > one(N + 1), two(N + 1);
      vector< int > three;
      for(int i = 0; i < N; i++) {
        one[i + 1] = one[i] + (S[i] == '1');
        two[i + 1] = two[i] + (S[i] == '2');
        if(S[i] == '/') three.emplace_back(i);
      }
      while(Q--) {
        int L, R;
        cin >> L >> R;
        --L;
        int l = lower_bound(three.begin(), three.end(), L) - three.begin();
        int r = lower_bound(three.begin(), three.end(), R) - three.begin();
        int ret = 0;
        for(int i = l; i < r; i++) {
          int j = three[i];
          ret = max(ret, min(one[j] - one[L], two[R] - two[j]) * 2 + 1);
        }
        cout << ret << endl;
      }
    }
    ```

11. 网页版**o1-mini 有时更加厉害**，copilot-chat 的 preview 更加厉害。不需要角色。
    模型问题，要多试。
12. 力扣**反馈错误格式**
    返回错误:

    ```
    提交结果：解答错误

    输入：
    [[0,1,14],[0,3,44],[3,2,37]]
    1

    输出：
    44

    预期：
    51
    ```

13. **重复 prompt 陷入循环后：**
    `回答错误。需要更换解法`

14. 重新打开对话，重新抽签

---

你是一个算法专家，精通解算法题，智慧超群。
在当前聊天中，你需要做到以下要求：
你负责生成代码，代码要求：

- 生成的代码清晰可读
- 代码尽量简洁
- 代码不包含注释
- 使用驼峰命名法
- 合理使用空行分隔

请输出代码。

---

bonus:

1. 限定时间复杂度：你应该使用一个 `O(nlogn)` 的算法来解决这个问题。
2. 限定方法：你应该使用动态规划来解决这个问题。

---

anti-gpt:

atcoder 一般 E 题难用 gpt 求解，需要自己动手实现。
leetcode 有的 t3/t4 anti-gpt 也不一定能解决，需要自己动手实现。
