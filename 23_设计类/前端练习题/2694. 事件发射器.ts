type Callback = (...args: unknown[]) => unknown
type Subscription = {
  unsubscribe: () => void
}

class EventEmitter {
  private readonly _events: Record<string, Callback[]> = {}

  subscribe(eventName: string, callback: Callback): Subscription {
    if (!this._events[eventName]) this._events[eventName] = []
    this._events[eventName].push(callback)

    return {
      unsubscribe: () => {
        this._events[eventName] = this._events[eventName].filter(cb => cb !== callback)
      }
    }
  }

  emit(eventName: string, args: unknown[] = []): unknown {
    const callbacks = this._events[eventName] || []
    return callbacks.map(cb => cb(...args))
  }
}

/**
 * const emitter = new EventEmitter();
 *
 * // Subscribe to the onClick event with onClickCallback
 * function onClickCallback() { return 99 }
 * const sub = emitter.subscribe('onClick', onClickCallback);
 *
 * emitter.emit('onClick'); // [99]
 * sub.unsubscribe(); // undefined
 * emitter.emit('onClick'); // []
 */

export {}
