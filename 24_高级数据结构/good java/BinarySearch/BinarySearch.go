// package template.algo;

// import template.math.DigitUtils;

// import java.util.Comparator;
// import java.util.function.DoublePredicate;
// import java.util.function.IntPredicate;
// import java.util.function.LongPredicate;

// public class BinarySearch {
//     private BinarySearch() {
//     }

//     public static long lastTrue(LongPredicate predicate, long l, long r) {
//         assert l <= r;
//         while (l != r) {
//             long mid = DigitUtils.ceilAverage(l, r);
//             if (predicate.test(mid)) {
//                 l = mid;
//             } else {
//                 r = mid - 1;
//             }
//         }
//         if (!predicate.test(l)) {
//             l--;
//         }
//         return l;
//     }

//     public static long firstTrue(LongPredicate predicate, long l, long r) {
//         assert l <= r;
//         while (l != r) {
//             long mid = DigitUtils.floorAverage(l, r);
//             if (predicate.test(mid)) {
//                 r = mid;
//             } else {
//                 l = mid + 1;
//             }
//         }
//         if(!predicate.test(l)){
//             l++;
//         }
//         return l;
//     }

//     public static int lastTrue(IntPredicate predicate, int l, int r) {
//         assert l <= r;
//         while (l != r) {
//             int mid = DigitUtils.ceilAverage(l, r);
//             if (predicate.test(mid)) {
//                 l = mid;
//             } else {
//                 r = mid - 1;
//             }
//         }
//         if (!predicate.test(l)) {
//             l--;
//         }
//         return l;
//     }

//     public static int firstTrue(IntPredicate predicate, int l, int r) {
//         assert l <= r;
//         while (l != r) {
//             int mid = DigitUtils.floorAverage(l, r);
//             if (predicate.test(mid)) {
//                 r = mid;
//             } else {
//                 l = mid + 1;
//             }
//         }
//         if(!predicate.test(l)){
//             l++;
//         }
//         return l;
//     }

//     public static double firstTrue(DoublePredicate predicate, double l, double r, double absError, double relativeError) {
//         if (l > r) {
//             throw new IllegalArgumentException();
//         }
//         while (r - l > absError) {
//             if ((r < 0 && (r - l) < -r * relativeError) || (l > 0 && (r - l) < l * relativeError)) {
//                 break;
//             }
//             double mid = (l + r) / 2;
//             if (predicate.test(mid)) {
//                 r = mid;
//             } else {
//                 l = mid;
//             }
//         }
//         return (l + r) / 2;
//     }

//     public static double lastTrue(DoublePredicate predicate, double l, double r, double absError, double relativeError) {
//         return firstTrue(predicate, l, r, absError, relativeError);
//     }

//     public static int lowerBound(int[] arr, int l, int r, int target) {
//         while (l < r) {
//             int mid = DigitUtils.floorAverage(l, r);
//             if (arr[mid] >= target) {
//                 r = mid;
//             } else {
//                 l = mid + 1;
//             }
//         }
//         if (arr[l] < target) {
//             l++;
//         }
//         return l;
//     }

//     public static int upperBound(int[] arr, int l, int r, int target) {
//         return lowerBound(arr, l, r, target + 1);
//     }

//     public static int lowerBound(long[] arr, int l, int r, long target) {
//         while (l < r) {
//             int mid = DigitUtils.floorAverage(l, r);
//             if (arr[mid] >= target) {
//                 r = mid;
//             } else {
//                 l = mid + 1;
//             }
//         }
//         if (arr[l] < target) {
//             l++;
//         }
//         return l;
//     }

//     public static int upperBound(long[] arr, int l, int r, long target) {
//         return lowerBound(arr, l, r, target + 1);
//     }

//     public static int lowerBound(double[] arr, int l, int r, double target) {
//         while (l < r) {
//             int mid = DigitUtils.floorAverage(l, r);
//             if (arr[mid] >= target) {
//                 r = mid;
//             } else {
//                 l = mid + 1;
//             }
//         }
//         if (arr[l] < target) {
//             l++;
//         }
//         return l;
//     }

//     public static int upperBound(double[] arr, int l, int r, double target) {
//         while (l < r) {
//             int mid = DigitUtils.floorAverage(l, r);
//             if (arr[mid] > target) {
//                 r = mid;
//             } else {
//                 l = mid + 1;
//             }
//         }
//         if (arr[l] <= target) {
//             l++;
//         }
//         return l;
//     }

//     public static <T> int lowerBound(T[] arr, int l, int r, T target, Comparator<T> comp) {
//         while (l < r) {
//             int mid = DigitUtils.floorAverage(l, r);
//             if (comp.compare(arr[mid], target) >= 0) {
//                 r = mid;
//             } else {
//                 l = mid + 1;
//             }
//         }
//         if (comp.compare(arr[l], target) < 0) {
//             l++;
//         }
//         return l;
//     }

//     public static <T> int upperBound(T[] arr, int l, int r, T target, Comparator<T> comp) {
//         while (l < r) {
//             int mid = DigitUtils.floorAverage(l, r);
//             if (comp.compare(arr[mid], target) > 0) {
//                 r = mid;
//             } else {
//                 l = mid + 1;
//             }
//         }
//         if (comp.compare(arr[l], target) <= 0) {
//             l++;
//         }
//         return l;
//     }

// }

// interface IBisectOptions<T> {
//   left?: number
//   right?: number
//   compareFn?: (a: T, b: T) => number
// }

// class BinarySearch {
//   private constructor() {}

//   static firstTrueInt(predicate: (mid: number) => boolean, left: number, right: number): number {}

//   static lastTrueInt(predicate: (mid: number) => boolean, left: number, right: number): number {}

//   static firstTrueFloat64(predicate: (mid: number) => boolean, left: number, right: number, absError: number, relativeError: number): number {}

//   static lastTrueFloat64(predicate: (mid: number) => boolean, left: number, right: number, absError: number, relativeError: number): number {
//     return BinarySearch.firstTrueFloat64(predicate, left, right, absError, relativeError)
//   }

//   static lowerBound<T>(arr: ArrayLike<T>, target: T, options?: IBisectOptions<T>): number {
//     const { left = 0, right = arr.length - 1, compareFn = (a: any, b: any) => (a < b ? -1 : a > b ? 1 : 0) } = options || {}
//   }

//   static upperBound<T>(arr: ArrayLike<T>, target: T, options?: IBisectOptions<T>): number {
//     const { left = 0, right = arr.length - 1, compareFn = (a: any, b: any) => (a < b ? -1 : a > b ? 1 : 0) } = options || {}
//   }
// }

// export { BinarySearch }

// if (require.main === module) {
// }

package main

func main() {

}
