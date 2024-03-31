package template.algo;

import template.primitve.generated.datastructure.DoubleFunction;

public class NewtonMethod {
    private double prec;

    public NewtonMethod(double prec) {
        this.prec = prec;
    }

    public double search(DoubleFunction func, DoubleFunction derivative, double x0) {
        while (Math.abs(func.apply(x0)) > prec) {
            x0 = x0 - func.apply(x0) / derivative.apply(x0);
        }
        return x0;
    }
}
