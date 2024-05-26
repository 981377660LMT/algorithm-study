package template.problem;

import template.math.IntMath;
import template.math.LongPollardRho;

import java.util.Arrays;
import java.util.Set;

public class SumOfSquares {
    private static int[] cantBeRepresent = new int[]{2, 3, 6, 7, 8, 11, 12, 15, 18, 19, 22, 23, 24, 27, 28, 31, 32, 33, 43, 44, 47, 48, 60, 67, 72, 76, 92, 96, 108, 112, 128};

    public static boolean isSumOfDistinctSquares(long n) {
        return !(n < 0 || n <= 128 && Arrays.binarySearch(cantBeRepresent, (int) n) >= 0);
    }

    public static boolean isSumOfOneSquare(long n) {
        if (n == 0) {
            return true;
        }
        long sqrt = IntMath.floorSqrt(n);
        return sqrt * sqrt == n;
    }

    public static boolean isSumOfTwoSquare(long n) {
        if (n == 0) {
            return true;
        }
        Set<Long> set = LongPollardRho.findAllFactors(n);
        for (long x : set) {
            int pow = 0;
            long y = n;
            while (y % x == 0) {
                y /= x;
                pow++;
            }
            if (x % 4 == 3 && pow % 2 == 1) {
                return false;
            }
        }
        return true;
    }

    public static boolean isSumOfThreeSquare(long n) {
        if (n == 0) {
            return true;
        }
        while (true) {
            if (n % 8 == 7) {
                return false;
            }
            if (n % 4 != 0) {
                break;
            }
            n /= 4;
        }
        return true;
    }

    public static boolean isSumOfFourSquare(long n) {
        return true;
    }
}
