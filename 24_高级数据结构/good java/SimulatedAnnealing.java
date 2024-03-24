package template.rand;

public abstract class SimulatedAnnealing<S> {
    public SimulatedAnnealing(double threshold, double k, double reduce) {
        this.threshold = threshold;
        this.k = k;
        this.reduce = reduce;
    }

    public abstract S next(S old, double temperature);

    public abstract double eval(S status);

    public void abandon(S old) {
    }

    public void optimize(double temperature, S init) {
        S now = init;
        double weight = eval(now);
        double t = temperature;
        while (t > threshold) {
            S next = next(now, t);
            double nextWeight = eval(next);
            if (nextWeight > weight || RandomWrapper.INSTANCE.nextDouble() < Math.exp((nextWeight - weight) / (k * t))) {
                abandon(now);
                now = next;
                weight = nextWeight;
            }
            t *= reduce;
        }

        if (best == null || bestWeight < weight) {
            best = now;
            bestWeight = weight;
        }
    }

    public S getBest() {
        return best;
    }

    public double weightOfBest(){
        return bestWeight;
    }

    private S best;
    private double bestWeight = -1e100;
    private double threshold;
    /**
     * The larger k is, the more possible to challenge
     */
    private double k;
    /**
     * The smaller reduce is, the fast to reduce temperature
     */
    private double reduce;
}