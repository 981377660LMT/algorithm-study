class Solution {
  private static final int MOD = (int) 1e9 + 7;

  public int firstDayBeenInAllRooms(int[] nextVisit) {
    int[] dp = new int[nextVisit.length], dpSum = new int[nextVisit.length];
    for (int i = 1; i < nextVisit.length; i++) {
      dp[i] = (dpSum[i - 1] - dpSum[nextVisit[i - 1]] + MOD) % MOD;
      dpSum[i] = (dpSum[i - 1] + dp[i]) % MOD;
    }

    return dpSum[nextVisit.length - 1];

  }

}