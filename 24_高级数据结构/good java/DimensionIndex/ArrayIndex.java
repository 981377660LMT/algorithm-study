package template.utils;

public class ArrayIndex {
  int[] dimensions;

  public ArrayIndex(int... dimensions) {
    this.dimensions = dimensions;
  }

  public int totalSize() {
    int ans = 1;
    for (int x : dimensions) {
      ans *= x;
    }
    return ans;
  }

  public int indexOf(int a, int b) {
    return a * dimensions[1] + b;
  }

  public int indexOf(int a, int b, int c) {
    return (a * dimensions[1] + b) * dimensions[2] + c;
  }

  public int indexOf(int a, int b, int c, int d) {
    return ((a * dimensions[1] + b) * dimensions[2] + c) * dimensions[3] + d;
  }

  public int indexOf(int a, int b, int c, int d, int e) {
    return (((a * dimensions[1] + b) * dimensions[2] + c) * dimensions[3] + d) * dimensions[4] + e;
  }

  public int indexOf(int a, int b, int c, int d, int e, int f) {
    return ((((a * dimensions[1] + b) * dimensions[2] + c) * dimensions[3] + d) * dimensions[4] + e)
        * dimensions[5] + f;
  }

  public int indexOf(int a, int b, int c, int d, int e, int f, int g) {
    return (((((a * dimensions[1] + b) * dimensions[2] + c) * dimensions[3] + d) * dimensions[4]
        + e) * dimensions[5] + f) * dimensions[6] + g;
  }

  public boolean isValid(int a, int d) {
    return dimensions[d] > a && a >= 0;
  }

  public boolean isValidIndex(int a) {
    return isValid(a, 0);
  }

  public boolean isValidIndex(int a, int b) {
    return isValidIndex(a) && isValid(b, 1);
  }

  public boolean isValidIndex(int a, int b, int c) {
    return isValidIndex(a, b) && isValid(c, 2);
  }

  public boolean isValidIndex(int a, int b, int c, int d) {
    return isValidIndex(a, b, c) && isValid(d, 3);
  }

  public int indexOfSpecifiedDimension(int index, int d) {
    return indexOfSpecifiedDimension0(index, d, dimensions.length - 1);
  }

  public int[] inverse(int i) {
    int[] ans = new int[dimensions.length];
    for (int j = dimensions.length - 1; j >= 0; j--) {
      ans[j] = i % dimensions[j];
      i /= dimensions[j];
    }
    return ans;
  }

  private int indexOfSpecifiedDimension0(int index, int t, int now) {
    return now == t ? index % dimensions[now]
        : indexOfSpecifiedDimension0(index / dimensions[now], t, now - 1);
  }
}


