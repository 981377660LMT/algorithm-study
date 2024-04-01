
// // https://www.cnblogs.com/matt-su/p/16101002.html

// import java.util.Random;

// /**
//  * InnerSimulatedAnnealing
//  */
// public class InnerSimulatedAnnealing {
//   public static void main(String[] args) {
//     double initialTemperature = 1000.0; // 初始温度
//     double threshold = 1e-8; // 温度阈值
//     double k = 1.0; // k值
//     double reduce = 0.99; // 温度衰减率
//     Integer initialSolution = 0; // 初始解

//     IntegerSolution sa = new IntegerSolution(threshold, k, reduce);
//     sa.optimize(initialTemperature, initialSolution); // 执行优化

//     Integer bestSolution = sa.getBest(); // 获取最优解
//     double bestWeight = sa.weightOfBest(); // 获取最优解的质量

//     System.out.println("Best solution: " + bestSolution);
//     System.out.println("Best solution weight: " + bestWeight);
//   }
// }

// public class RandomWrapper {
//   private Random random;

//   public RandomWrapper() {
//     this(new Random());
//   }

//   public RandomWrapper(Random random) {
//     this.random = random;
//   }

//   public RandomWrapper(long seed) {
//     this(new Random(seed));
//   }

//   public int nextInt(int l, int r) {
//     return random.nextInt(r - l + 1) + l;
//   }

//   public int nextInt(int n) {
//     return random.nextInt(n);
//   }

//   public double nextDouble(double l, double r) {
//     return random.nextDouble() * (r - l) + l;
//   }

//   public double nextDouble() {
//     return random.nextDouble();
//   }

//   public long nextLong(long l, long r) {
//     return nextLong(r - l + 1) + l;
//   }

//   public long nextLong(long n) {
//     return Math.round(random.nextDouble() * (n - 1));
//   }

//   public long nextLong() {
//     return random.nextLong();
//   }

//   public String nextString(char l, char r, int len) {
//     StringBuilder builder = new StringBuilder(len);
//     for (int i = 0; i < len; i++) {
//       builder.append((char) nextInt(l, r));
//     }
//     return builder.toString();
//   }

//   public String nextString(char[] s, int len) {
//     StringBuilder builder = new StringBuilder(len);
//     for (int i = 0; i < len; i++) {
//       builder.append(s[nextInt(0, s.length - 1)]);
//     }
//     return builder.toString();
//   }

//   public Random getRandom() {
//     return random;
//   }

//   public int range(int... x) {
//     return x[nextInt(0, x.length - 1)];
//   }

//   public char range(char... x) {
//     return x[nextInt(0, x.length - 1)];
//   }

//   public long range(long... x) {
//     return x[nextInt(0, x.length - 1)];
//   }

//   public <T> T rangeT(T... x) {
//     return x[nextInt(0, x.length - 1)];
//   }

//   public static final RandomWrapper INSTANCE = new RandomWrapper();

//   public static void main(String[] args) {
//     RandomWrapper random = new RandomWrapper();
//     System.out.println(random.nextInt(1, 10));
//     System.out.println(random.nextDouble(1.0, 10.0));
//     System.out.println(random.nextLong(1, 10));
//     System.out.println(random.nextString('a', 'z', 10));
//     System.out.println(random.nextString(new char[] {'a', 'b', 'c'}, 10));
//   }
// }

// // 假设我们的解是一个整数
// class IntegerSolution extends SimulatedAnnealing<Integer> {
//   public IntegerSolution(double threshold, double k, double reduce) {
//     super(threshold, k, reduce);
//   }

//   @Override
//   public Integer next(Integer old, double temperature) {
//     // 根据当前温度生成新的解，这里只是示例
//     return old + (RandomWrapper.INSTANCE.nextInt(3) - 1);
//   }

//   @Override
//   public double eval(Integer status) {
//     // 评估解的质量，这里只是示例
//     System.out.println("Eval: " + status);
//     return -Math.abs(status - 100); // 假设我们的目标是接近100
//   }

//   @Override
//   public void abandon(Integer old) {
//     // 处理舍弃的解，这里不需要实现
//   }

//   // 使用模拟退火算法
//   public static void main(String[] args) {
//     double initialTemperature = 1000.0; // 初始温度
//     double threshold = 1e-8; // 温度阈值
//     double k = 1.0; // k值
//     double reduce = 0.99; // 温度衰减率
//     Integer initialSolution = 0; // 初始解

//     IntegerSolution sa = new IntegerSolution(threshold, k, reduce);
//     sa.optimize(initialTemperature, initialSolution); // 执行优化

//     Integer bestSolution = sa.getBest(); // 获取最优解
//     double bestWeight = sa.weightOfBest(); // 获取最优解的质量

//     System.out.println("Best solution: " + bestSolution);
//     System.out.println("Best solution weight: " + bestWeight);
//   }
// }

// /**
//  * 模拟退火优化.
//  */
// public abstract class SimulatedAnnealing<S> {
//   private S best;
//   private double bestWeight = -1e100;

//   private double threshold;

//   /**
//    * 玻尔兹曼常数.
//    *
//    * The larger k is, the more possible to challenge .
//    */
//   private double k;

//   /**
//    * 学习率.
//    *
//    * The smaller reduce is, the fast to reduce temperature
//    */
//   private double reduce;

//   public SimulatedAnnealing(double threshold, double k, double reduce) {
//     this.threshold = threshold;
//     this.k = k;
//     this.reduce = reduce;
//   }

//   public abstract S next(S old, double temperature);

//   public abstract double eval(S status);

//   public void abandon(S old) {}

//   public void optimize(double temperature, S init) {
//     S now = init;
//     double weight = eval(now);
//     double t = temperature;
//     while (t > threshold) {
//       S next = next(now, t);
//       double nextWeight = eval(next);
//       if (nextWeight > weight
//           || RandomWrapper.INSTANCE.nextDouble() < Math.exp((nextWeight - weight) / (k * t))) {
//         abandon(now);
//         now = next;
//         weight = nextWeight;
//       }
//       t *= reduce;
//     }

//     if (best == null || bestWeight < weight) {
//       best = now;
//       bestWeight = weight;
//     }
//   }

//   public S getBest() {
//     return best;
//   }

//   public double weightOfBest() {
//     return bestWeight;
//   }
// }
