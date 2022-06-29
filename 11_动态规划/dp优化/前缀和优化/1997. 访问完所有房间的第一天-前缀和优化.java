class Solution {
  private static final int MOD = (int) 1e9 + 7;

  public int firstDayBeenInAllRooms(int[] nextVisit) {
    int[] dp = new int[nextVisit.length], sum = new int[nextVisit.length];
    for (int i = 1; i < nextVisit.length; i++) {
      dp[i] = (sum[i - 1] - sum[nextVisit[i - 1]] + MOD) % MOD;
      sum[i] = (sum[i - 1] + dp[i]) % MOD;
    }

    return sum[nextVisit.length - 1];

  }

}