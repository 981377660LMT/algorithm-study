/**
 * 维护双向映射.
 */
export class BidirectionalMap<K, V> {
  private readonly keyToValue: Map<K, V> = new Map();
  private readonly valueToKey: Map<V, K> = new Map();

  constructor(entries?: Iterable<[K, V]>) {
    if (entries) {
      for (const [key, value] of entries) {
        this.set(key, value);
      }
    }
  }

  set(key: K, value: V): this {
    if (this.keyToValue.has(key)) {
      const oldValue = this.keyToValue.get(key)!;
      this.valueToKey.delete(oldValue);
    }
    if (this.valueToKey.has(value)) {
      const oldKey = this.valueToKey.get(value)!;
      this.keyToValue.delete(oldKey);
    }

    this.keyToValue.set(key, value);
    this.valueToKey.set(value, key);
    return this;
  }

  getValue(key: K): V | undefined {
    return this.keyToValue.get(key);
  }

  getKey(value: V): K | undefined {
    return this.valueToKey.get(value);
  }

  deleteByKey(key: K): boolean {
    const value = this.keyToValue.get(key);
    if (value !== undefined) {
      this.keyToValue.delete(key);
      this.valueToKey.delete(value);
      return true;
    }
    return false;
  }

  deleteByValue(value: V): boolean {
    const key = this.valueToKey.get(value);
    if (key !== undefined) {
      this.valueToKey.delete(value);
      this.keyToValue.delete(key);
      return true;
    }
    return false;
  }

  /**
   * 更新 key（保持 value 不变）
   */
  updateKey(oldKey: K, newKey: K): boolean {
    const value = this.keyToValue.get(oldKey);
    if (value !== undefined) {
      this.keyToValue.delete(oldKey);
      this.keyToValue.set(newKey, value);
      this.valueToKey.set(value, newKey);
      return true;
    }
    return false;
  }

  hasKey(key: K): boolean {
    return this.keyToValue.has(key);
  }

  hasValue(value: V): boolean {
    return this.valueToKey.has(value);
  }

  get size(): number {
    return this.keyToValue.size;
  }

  clear(): void {
    this.keyToValue.clear();
    this.valueToKey.clear();
  }

  toRecord(): Record<string, V> {
    const result: Record<string, V> = {};
    this.keyToValue.forEach((value, key) => {
      result[String(key)] = value;
    });
    return result;
  }
}
