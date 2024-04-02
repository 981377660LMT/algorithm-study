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

// /**
//  * 模拟退火优化.
//  */
//  public abstract class SimulatedAnnealing<S> {
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

// 模拟退火优化求解最大值(最优解).
// https://www.cnblogs.com/shenben/p/11342308.html
// https://vlight.me/2018/06/08/Simulated-Annealing/
// https://oi-wiki.org//misc/simulated-annealing/
// https://www.luogu.com/article/b2recz8n
// 技巧：可以在时限内重复跑 SA 取最优值，防止脸黑

package main

import (
	"math"
	"math/rand"
	"time"
)

func main() {

}

// 模板题 https://www.luogu.com.cn/problem/P1337
// LC1515 https://leetcode.cn/problems/best-position-for-a-service-centre/
// http://poj.org/problem?id=2420
// UVa 10228 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=14&page=show_problem&problem=1169
// todo 教学题 https://atcoder.jp/contests/intro-heuristics/tasks/intro_heuristics_a
//  https://atcoder.jp/contests/ahc001/tasks/ahc001_a
//  https://atcoder.jp/contests/ahc002/tasks/ahc002_a

// 参数说明:
//
//	initTemperature: 初始温度.一般设置为 1000 到 3000.
//	threshold: 目标温度阈值.一般设置为 1e-8 到 1e-14.
//	k: 玻尔兹曼常数.一般设置为 1.
//	!reduce: 温度衰减率.一般设置为 0.99 到 0.999.对答案影响较大.
//	initSolution: 初始解.
//
// 如何生成新解：
//
//	坐标系内：随机生成一个点，或者生成一个向量。
//	序列问题：random.shuffle()或者随机交换两个元素。
//	网格问题：可以看做二维序列，每次交换两个格子即可。
type SimulatedAnnealing[S any] struct {
	best       S       // 最优解.
	bestWeight float64 // 最优解的分数.

	reduce               float64 // 温度衰减率.默认值为 0.99.
	initTemperature      float64 // 初始温度.默认值为 2000.
	thresholdTemperature float64 // 温度阈值.默认值为 1e-14.
	k                    float64 // 玻尔兹曼常数.默认值为 1.
	timeLimit            float64 // 时间限制.默认为-1，表示不限制时间.

	next    func(old S, temperature float64) S // 生成新解.
	eval    func(status S) float64             // 评估解的质量.
	abandon func(old S)                        // 处理舍弃的解.

	calculated bool
}

func NewSimulatedAnnealing[S any](
	next func(old S, temperature float64) S,
	eval func(status S) float64,
	abandon func(old S),
) *SimulatedAnnealing[S] {
	return &SimulatedAnnealing[S]{
		bestWeight:           -1e100,
		reduce:               0.99,
		initTemperature:      2000,
		thresholdTemperature: 1e-14,
		k:                    1,
		timeLimit:            -1,
		next:                 next,
		eval:                 eval,
		abandon:              abandon,
	}
}

func (sa *SimulatedAnnealing[S]) Optimize(initSolution S) {
	if sa.timeLimit == -1 {
		sa._run(initSolution)
	} else {
		sa._runWithinTimeLimit(initSolution, sa.timeLimit)
	}
}

func (sa *SimulatedAnnealing[S]) GetBest() S            { return sa.best }
func (sa *SimulatedAnnealing[S]) WeightOfBest() float64 { return sa.bestWeight }

func (sa *SimulatedAnnealing[S]) SetReduce(reduce float64) *SimulatedAnnealing[S] {
	sa.reduce = reduce
	return sa
}

func (sa *SimulatedAnnealing[S]) SetInitTemperature(initTemperature float64) *SimulatedAnnealing[S] {
	sa.initTemperature = initTemperature
	return sa
}

func (sa *SimulatedAnnealing[S]) SetThresholdTemperature(thresholdTemperature float64) *SimulatedAnnealing[S] {
	sa.thresholdTemperature = thresholdTemperature
	return sa
}

func (sa *SimulatedAnnealing[S]) SetK(k float64) *SimulatedAnnealing[S] {
	sa.k = k
	return sa
}

func (sa *SimulatedAnnealing[S]) SetTimeLimit(timeLimit float64) *SimulatedAnnealing[S] {
	sa.timeLimit = timeLimit
	return sa
}

func (sa *SimulatedAnnealing[S]) _run(initSolution S) {
	now := initSolution
	weight := sa.eval(now)
	t := sa.initTemperature
	for t > sa.thresholdTemperature {
		next := sa.next(now, t)
		nextWeight := sa.eval(next)
		if nextWeight > weight || rand.Float64() < math.Exp((nextWeight-weight)/(sa.k*t)) {
			sa.abandon(now)
			now = next
			weight = nextWeight
		}
		t *= sa.reduce
	}

	if !sa.calculated || sa.bestWeight < weight {
		sa.best = now
		sa.bestWeight = weight
		sa.calculated = true
	}
}

func (sa *SimulatedAnnealing[S]) _runWithinTimeLimit(initSolution S, timeLimit float64) {
	t0 := time.Now()
	now := initSolution
}
