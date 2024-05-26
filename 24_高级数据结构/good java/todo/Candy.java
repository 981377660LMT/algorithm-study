package template.problem;

import template.algo.BinarySearch;
import template.math.DigitUtils;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;

/**
 * This problem is about on a circle, there are c points enumerated as 0, 1, ...,c - 1.
 * <br>
 * There are some candis on some locations, xi means there are xi candis on point x.
 * <br>
 * You can deliver candy between adjacent points, (0 and c - 1 is adjacent), deliver each candy one unit cost 1.
 * <br>
 * You have to ensure location i has yi candies, while x_0 + x_1 + ... + x_{c_1} = y_0 + y_1 + ... + y_{c_1}.
 */
public class CandyAssignProblem {
    private long c;
    List<Candy> addedCandies;
    Candy[] candies;
    int candieCnt;
    long minimumCost;


    public CandyAssignProblem(long c, int exp) {
        this.c = c;
        addedCandies = new ArrayList<>(exp);
    }

    public void requestOn(long i, long x, long y) {
        Candy candy = new Candy();
        candy.location = i;
        candy.x = x;
        candy.y = y;
        addedCandies.add(candy);
    }

    public void assignCandyOn(long i, long x) {
        requestOn(i, x, 0);
    }

    public void requireCandyOn(long i, long y) {
        requestOn(i, 0, y);
    }

    public long solve() {
        if (addedCandies.isEmpty()) {
            addedCandies.add(new Candy());
        }
        candies = addedCandies.toArray(new Candy[0]);
        Arrays.sort(candies, Candy.sortByLocation);
        candieCnt = 0;
        for (int i = 1; i < candies.length; i++) {
            if (candies[i].location == candies[candieCnt].location) {
                candies[candieCnt].x += candies[i].x;
                candies[candieCnt].y += candies[i].y;
            } else {
                candieCnt++;
                candies[candieCnt] = candies[i];
            }
        }
        candieCnt++;

        for (int i = 0; i < candieCnt - 1; i++) {
            candies[i].w = candies[i + 1].location - candies[i].location;
        }
        candies[candieCnt - 1].w = c + candies[0].location - candies[candieCnt - 1].location;

        for (int i = 1; i < candieCnt; i++) {
            candies[i].a = candies[i - 1].a + candies[i].x - candies[i].y;
        }

        Candy[] sortedByA = Arrays.copyOf(candies, candieCnt);
        Arrays.sort(sortedByA, (a, b) -> Long.compare(a.a, b.a));
        long prefix = 0;
        long half = DigitUtils.ceilDiv(c, 2);
        for (int i = 0; i < candieCnt; i++) {
            prefix += sortedByA[i].w;
            if (prefix >= half) {
                candies[0].a = -sortedByA[i].a;
                break;
            }
        }

        for (int i = 1; i < candieCnt; i++) {
            candies[i].a += candies[0].a;
        }

        for (int i = 0; i < candieCnt; i++) {
            minimumCost += Math.abs(candies[i].a) * candies[i].w;
        }

        return minimumCost;
    }

    public long minimumCost() {
        return minimumCost;
    }


    Candy tmp = new Candy();

    /**
     * How many candies are delivered from i to i + 1, the returned value could be negative.
     */
    public long deliverBetween(long i) {
        tmp.location = i;
        int index = BinarySearch.lowerBound(candies, 0, candieCnt - 1, tmp, Candy.sortByLocation);
        if (index < 0) {
            index = candieCnt - 1;
        }
        return candies[index].a;
    }

    static class Candy {
        long location;
        long x;
        long y;

        long a;
        long w;

        static Comparator<Candy> sortByLocation = (a, b) -> Long.compare(a.location, b.location);
    }
}