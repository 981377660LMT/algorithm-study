// !java 预处理+多维dp

import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Deque;
import java.util.List;

// 11_动态规划/经典题/LCP 69. Hello LeetCode!.py
// 1. 考虑 dp, $dp[index][C][D][E][T][O][L][H]表示 前 index 个单词中,每种字母还要取多少个
// 2. 对每个单词，枚举子集预处理合理的转移状态
// 3. 当删除的字符固定时，从两端向中间删除字符，双指针或 deque 计算最小删除成本 (实际上，这里的成本也可以预处理)
// dp 的状态数为 $n*960$，转移的复杂度为 $O(2^{len(w)}*len(w))$
class Solution {
  private static final String CDETOLH = "cdetolh";

  public int Leetcode(String[] words) {
    int n = words.length;

    // !1.枚举子集预处理合理的转移状态
    List<Integer>[] states = new List[n];
    for (int i = 0; i < n; i++) {
      states[i] = new ArrayList<>();
      int m = words[i].length();
      for (int state = 0; state < 1 << m; state++) {
        boolean isValid = true;
        for (int j = 0; j < m; j++) {
          if (((state >> j) & 1) != 0 && CDETOLH.indexOf(words[i].charAt(j)) == -1) {
            isValid = false;
            break;
          }
        }

        if (isValid) {
          states[i].add(state);
        }
      }
    }

    int[][][][][][][] dp = new int[2][2][5][2][3][4][2], ndp = new int[2][2][5][2][3][4][2]; // CDETOLH
    for (int c = 0; c < 2; ++c) {
      for (int d = 0; d < 2; ++d) {
        for (int e = 0; e < 5; ++e) {
          for (int t = 0; t < 2; ++t) {
            for (int o = 0; o < 3; ++o) {
              for (int l = 0; l < 4; ++l) {
                for (int h = 0; h < 2; ++h) {
                  dp[c][d][e][t][o][l][h] = ndp[c][d][e][t][o][l][h] = -1;
                }
              }
            }
          }
        }
      }
    }
    dp[1][1][4][1][2][3][1] = 0;

    for (int i = 0; i < n; i++) {
      for (int state : states[i]) {
        int[] C = new int[7];
        int m = words[i].length();

        // !2.双指针或deque计算最小删除成本
        Deque<Integer> queue = new ArrayDeque<>();
        for (int j = 0; j < m; j++) {
          if (((state >> j) & 1) != 0) {
            C[CDETOLH.indexOf(words[i].charAt(j))]++;
            queue.offerLast(j);
          }
        }

        int cost = 0;
        int leftMoved = 0, rightMoved = 0;
        while (!queue.isEmpty()) {
          int left = queue.peekFirst(), right = queue.peekLast();
          int leftCost = (left - leftMoved) * (m - 1 - left - rightMoved);
          int rightCost = (right - leftMoved) * (m - 1 - right - rightMoved);
          if (leftCost <= rightCost) {
            queue.pollFirst();
            cost += leftCost;
            leftMoved++;
          } else {
            queue.pollLast();
            cost += rightCost;
            rightMoved++;
          }
        }

        // !3.状态转移
        for (int c = C[0]; c < 2; c++) {
          for (int d = C[1]; d < 2; d++) {
            for (int e = C[2]; e < 5; e++) {
              for (int t = C[3]; t < 2; t++) {
                for (int o = C[4]; o < 3; o++) {
                  for (int l = C[5]; l < 4; l++) {
                    for (int h = C[6]; h < 2; h++) {
                      if (dp[c][d][e][t][o][l][h] != -1) {
                        if (ndp[c - C[0]][d - C[1]][e - C[2]][t - C[3]][o - C[4]][l - C[5]][h - C[6]] == -1) {
                          ndp[c - C[0]][d - C[1]][e - C[2]][t - C[3]][o - C[4]][l - C[5]][h - C[6]] = cost
                              + dp[c][d][e][t][o][l][h];
                        } else {
                          ndp[c - C[0]][d - C[1]][e - C[2]][t - C[3]][o - C[4]][l - C[5]][h - C[6]] = Math.min(
                              ndp[c - C[0]][d - C[1]][e - C[2]][t - C[3]][o - C[4]][l - C[5]][h - C[6]],
                              dp[c][d][e][t][o][l][h] + cost);
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }

      for (int d = 0; d < 2; d++) {
        for (int c = 0; c < 2; c++) {
          for (int e = 0; e < 5; e++) {
            for (int t = 0; t < 2; t++) {
              for (int o = 0; o < 3; o++) {
                for (int l = 0; l < 4; l++) {
                  for (int h = 0; h < 2; h++) {
                    dp[c][d][e][t][o][l][h] = ndp[c][d][e][t][o][l][h];
                  }
                }
              }
            }
          }
        }
      }
    }

    return dp[0][0][0][0][0][0][0];
  }

  public static void main(String[] args) {
    Solution solution = new Solution();
    String[] words = new String[] { "hold", "engineer", "cost", "level" };
    int res = solution.Leetcode(words);
    System.out.println(res);
  }
}