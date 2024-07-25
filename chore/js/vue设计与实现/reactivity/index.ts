import {
  /** scope */
  effectScope,
  getCurrentScope,
  onScopeDispose,

  /** effect */
  effect,
  stop,
  enableTracking,
  pauseTracking,
  resetTracking,
  pauseScheduling,
  resetScheduling,
  ReactiveEffectOptions,
  DebuggerOptions,
  DebuggerEvent,
  DebuggerEventExtraInfo,
  ReactiveEffectRunner,
  ReactiveEffect,
  EffectScheduler,

  /** computed */
  computed,
  ComputedGetter,
  ComputedRef,
  WritableComputedOptions,
  ComputedSetter,

  /** reactive */
  reactive,
  shallowReactive,
  readonly,
  shallowReadonly,
  toRaw,
  isReactive,
  isReadonly,
  isShallow,
  isProxy,
  markRaw,
  DeepReadonly,

  /** ref */
  ref,
  shallowRef,
  customRef,
  unref,
  proxyRefs,
  toValue,
  toRef,
  toRefs,
  isRef,
  Ref,
  UnwrapRef,

  /** trigger/track */
  track,
  trigger,
  triggerRef,
  TrackOpTypes,
  TriggerOpTypes
} from '@vue/reactivity'

type Simplify<M extends object> = {
  [K in keyof M]: M[K]
}
type Foo = Simplify<typeof import('@vue/reactivity')>

function testCustomRef(): void {
  // 创建一个防抖ref
  function useDebouncedRef<T>(value: T, delay = 200) {
    let timeoutId: ReturnType<typeof setTimeout>
    return customRef((track, trigger) => {
      return {
        // 当ref的值被读取时，调用track来追踪依赖
        get() {
          track()
          return value
        },
        // 当ref的值被设置时，启动一个防抖逻辑
        set(newValue) {
          clearTimeout(timeoutId)
          timeoutId = setTimeout(() => {
            value = newValue
            // 通知依赖更新
            trigger()
          }, delay)
        }
      }
    })
  }

  const count = useDebouncedRef(0, 200)
  setInterval(() => {
    count.value++
  }, 100)
  effect(() => {
    console.log(count.value)
  })
}

// !pause -> reset -> enable
// **resetTracking**用于在暂停追踪后清理响应式系统的追踪状态，确保没有不正确的依赖被留下。它通常在你计划恢复追踪之前使用，作为一种清理操作。
// **enableTracking**用于在之前暂停追踪后恢复正常的依赖追踪。它确保之后的响应式状态读取操作能够被正常追踪。
// 在实际使用中，如果你暂停了依赖追踪并在暂停期间进行了响应式状态的读取，
// 最佳实践是在调用enableTracking恢复追踪之前，先使用resetTracking来清理追踪状态。这样可以避免潜在的问题，确保响应式系统的行为是正确和预期的。
function testEffect(): void {
  // 创建一个响应式状态
  const count = ref(0)

  // 定义一个调试器事件处理函数
  const debuggerEvent = (event: DebuggerEvent) => {
    console.log('Debugger Event:')
  }

  // 定义一个副作用调度器
  const scheduler: EffectScheduler = effect => {
    console.log(effect, 987)
  }

  // 定义副作用函数选项
  const options: ReactiveEffectOptions = {
    /** 当设置为true时，副作用不会在创建时立即执行。这允许副作用在需要时手动触发，常用于实现惰性计算属性. */
    lazy: true,
    /** 当设置为true时，允许副作用函数在其自身执行过程中被递归调用。默认情况下，为了避免无限循环，副作用函数在执行时不允许再次被触发. */
    allowRecurse: false,
    scheduler,
    onTrack: debuggerEvent,
    onTrigger: debuggerEvent
  }

  // 创建并运行一个副作用
  const runner = effect(() => {
    console.log('Effect run:', count.value)
  }, options)
  runner()

  // 更新响应式状态
  count.value++

  // 暂停追踪
  pauseTracking()
  count.value++
  // 重置追踪
  resetTracking()

  // 暂停调度
  pauseScheduling()
  count.value++
  // 重置调度
  resetScheduling()

  // 停止副作用
  stop(runner)

  // 重新启用追踪
  enableTracking()
  count.value++
}

function testComputed(): void {
  const count = ref(0)
  const readonlyCount = computed(() => count.value * 2, {
    onTrack(event) {
      // console.log('onTrack', event)
    },
    onTrigger(event) {
      // console.log('onTrigger', event)
    }
  })
  console.log(readonlyCount.value)
  count.value++

  const count2 = ref(0)
  const mutableCount = computed({
    get: () => count2.value * 2,
    set: val => {
      count2.value = val / 2
    }
  })
  console.log(mutableCount.value)
  mutableCount.value = 10
  console.log(count2.value)
}

function testEffectScope(): void {
  // effectScope (类似 golang Context, 管理父子关系)
  // !用于管理 父子嵌套 effect 的取消(生命周期).
  //
  // !当detached为true时，表示创建的EffectScope是一个独立的作用域。
  // !这意味着，即使它在另一个活跃的EffectScope内部被创建，它也不会自动成为父作用域的一部分。
  // !因此，当父作用域被停止时，脱离的（detached）作用域不会被自动停止，需要手动管理其生命周期
  const parentValue = reactive({ parent: 1 })
  const childValue = reactive({ child: 2 })
  const parentScope = effectScope()
  console.log({ parentScope })

  parentScope.run(() => {
    effect(() => {
      console.log('parentValue.parent', parentValue.parent)
    })

    onScopeDispose(() => {
      console.log('父作用域已停止')
    })

    const childScope = effectScope() // 创建子作用域，自动成为父作用域的一部分
    childScope.run(() => {
      onScopeDispose(() => {
        console.log('子作用域已停止')
      })

      console.log(getCurrentScope())
      effect(() => {
        console.log('childValue.child', childValue.child)
      })
    })
  })

  setInterval(() => {
    childValue.child++
  }, 100)

  // 假设在某个时刻，我们决定停止父作用域
  // 这将同时停止子作用域中的所有副作用(如果子作用域不是detached的话)
  setTimeout(() => {
    // fromParent参数是EffectScope内部使用的一个标志，用于区分作用域停止操作的触发来源，帮助Vue的响应式系统更精确地管理作用域及其副作用的生命周期
    parentScope.stop()
    console.log('父作用域及其子作用域已停止')
  }, 1000)
}

function testTrackAndTrigger(): void {
  // track
  const obj = reactive({ foo: 1 })
  let dummy
  effect(() => {
    dummy = obj.foo
  })
  console.log(dummy)
  track(obj, TrackOpTypes.GET, 'foo')

  // trigger/
  const rawState = { count: 0 }
  const state = reactive(rawState)
  effect(() => {
    console.log(state.count)
  })
  trigger(rawState, TriggerOpTypes.SET, 'count')

  // triggerRef
  // Force trigger effects that depends on a shallow ref.
  // !当你使用ref来存储一个对象或数组，并且你在不替换整个对象或数组的情况下修改了它的内部状态
  // !在这种情况下，你可能需要手动调用triggerRef来确保视图或其他依赖于这个ref的副作用函数能够得到更新
  const shallow = shallowRef({ foo: 1 })
  effect(() => {
    console.log(shallow.value)
  })
  shallow.value.foo = 2
  triggerRef(shallow)
}

if (require.main === module) {
  // testCustomRef()
  // testEffect()
  // testComputed()
  // testEffectScope()
  testTrackAndTrigger()
}
