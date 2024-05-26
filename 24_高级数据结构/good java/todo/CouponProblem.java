package template.problem;

import java.util.Arrays;
import java.util.Comparator;
import java.util.PriorityQueue;

// 股票问题
public class CouponProblem {
    /**
     * res[i] = 0 means not buy the good, res[1] = 1 means buy the good without coupon, res[i] = 2 means buy the good
     * with coupon.
     * The target is buy as many goods as possible.
     *
     *  O(n\log_2n)
     * @param price
     * @param discount
     * @param numCoupon
     * @return
     */
    public int[] solve(long[] price, long[] discount, int numCoupon, long money) {
        int n = price.length;
        Item[] items = new Item[n];
        for (int i = 0; i < n; i++) {
            items[i] = new Item();
            items[i].index = i;
            items[i].price = price[i];
            items[i].discount = discount[i];
            items[i].profit = items[i].price - items[i].discount;
        }
        int[] ans = new int[n];
        PriorityQueue<Item> discountGroup = new PriorityQueue<>(n, Comparator.comparingLong(x -> x.profit));
        PriorityQueue<Item> withoutDiscount = new PriorityQueue<>(n, Comparator.comparingLong(x -> x.price));
        PriorityQueue<Item> withDiscount = new PriorityQueue<>(n, Comparator.comparingLong(x -> x.discount));
        withoutDiscount.addAll(Arrays.asList(items));
        withDiscount.addAll(Arrays.asList(items));
        while (!withDiscount.isEmpty() && !withoutDiscount.isEmpty()) {
            if (ans[withDiscount.peek().index] != 0) {
                withDiscount.remove();
                continue;
            }
            if (ans[withoutDiscount.peek().index] != 0) {
                withoutDiscount.remove();
                continue;
            }
            long discountTop = 0;
            if (numCoupon == 0) {
                discountTop = money + 1;
            } else if (discountGroup.size() == numCoupon) {
                discountTop = discountGroup.peek().profit;
            }
            if (discountTop + withDiscount.peek().discount <= withoutDiscount.peek().price) {
                Item top = withDiscount.remove();
                long cost = discountTop + top.discount;
                if (cost > money) {
                    break;
                }
                money -= cost;
                ans[top.index] = 2;
                if (discountGroup.size() == numCoupon) {
                    ans[discountGroup.remove().index] = 1;
                }
                discountGroup.add(top);
            } else {
                Item top = withoutDiscount.remove();
                long cost = top.price;
                if (cost > money) {
                    break;
                }
                money -= cost;
                ans[top.index] = 1;
            }
        }
        return ans;
    }

    static class Item {
        long price;
        long discount;
        long profit;
        int index;
    }
}
