class Solution {

	public int firstDayBeenInAllRooms(int[] nextVisit) {
		int[] dp = new int[nextVisit.length], sum = new int[nextVisit.length];
		for (int i = 1; i < nextVisit.length; i++) {
			dp[i] = (sum[i - 1] - sum[nextVisit[i - 1]] + 1000000009) % 1000000007;
			sum[i] = (sum[i - 1] + dp[i]) % 1000000007;
		}

    
		return sum[nextVisit.length - 1];
      
	}
}