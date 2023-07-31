
import java.util.function.Supplier;

/**
 * 惰性求值的流.类似IO函子.
 */
public class MyStream<T> {
  private final Supplier<T> value;

  private MyStream(Supplier<T> effect) {
    this.value = effect;
  }

  @SuppressWarnings("unchecked")
  public static <V> MyStream<V> of(V value) {
    if (value instanceof MyStream) {
      return new MyStream<>(((MyStream<V>) value)::run);
    }
    return new MyStream<>(() -> value);
  }

  public <U> MyStream<U> map(java.util.function.Function<T, U> f) {
    return new MyStream<>(() -> f.apply(value.get()));
  }

  public <U> MyStream<U> flatMap(java.util.function.Function<T, MyStream<U>> f) {
    return new MyStream<>(() -> f.apply(value.get()).run());
  }

  public T run() {
    return value.get();
  }

  public static void main(String[] args) {
    Object res = MyStream.of(1).map(x -> x + 1).flatMap(x -> MyStream.of(x + 1)).run();
    System.out.println(res);
  }
}
