class Solution {
public:
    int lastStoneWeightII(vector<int>& stones) {
        int sum = accumulate(stones.begin(), stones.end(), 0);
        int volumn = sum >> 1;
        vector<bool> dp(volumn + 1, false);
        dp[0] = true;

        for (int i = 0; i < stones.size(); i++) {
            for (int j = volumn; j >= 0; j--) {
                j >= stones[i] && (dp[j] = dp[j] || dp[j - stones[i]]);
            }
        }

        int maxWeight = find(dp.rbegin(), dp.rend(), true) - dp.rbegin();
        return sum - 2 * maxWeight;
    }
};