package template.binary;

/**
 * Log operations
 */
public class Log2 {
    public static int ceilLog(int x) {
        if (x <= 0) {
            return 0;
        }
        return 32 - Integer.numberOfLeadingZeros(x - 1);
    }

    public static int floorLog(int x) {
        if (x <= 0) {
            throw new IllegalArgumentException();
        }
        return 31 - Integer.numberOfLeadingZeros(x);
    }

    public static int ceilLog(long x) {
        if (x <= 0) {
            return 0;
        }
        return 64 - Long.numberOfLeadingZeros(x - 1);
    }

    public static int floorLog(long x) {
        if (x <= 0) {
            throw new IllegalArgumentException();
        }
        return 63 - Long.numberOfLeadingZeros(x);
    }
}
