// https://juejin.cn/post/6844903570496733192
// 迷你版本的 redux

// const store = {
//   state: {}, // 全局唯一的state，内部变量，通过getState()获取
//   listeners: [], // listeners，用来诸如视图更新的操作
//   dispatch: () => {}, // 分发action
//   subscribe: () => {}, // 用来订阅state变化
//   getState: () => {}, // 获取state
// }

const createStore = (reducer, initialState) => {
  const store = {}
  store.state = initialState
  store.listeners = []

  store.subscribe = listener => store.listeners.push(listener)
  store.dispatch = action => {
    store.state = reducer(store.state, action)
    store.listeners.forEach(listener => listener())
  }
  store.getState = () => store.state

  return store
}

// reducer
function counter(state = 0, action) {
  switch (action.type) {
    case 'INCREMENT':
      return state + 1
    case 'DECREMENT':
      return state - 1
    default:
      return state
  }
}

const store = createStore(counter)
store.subscribe(() => console.log(store.getState())) // 如果需要更新view，就根据我们暴露的subscribe去更新就好了
store.dispatch({ type: 'INCREMENT' })
store.dispatch({ type: 'INCREMENT' })
store.dispatch({ type: 'DECREMENT' })

// 思路
// 1.
// redux中核心就是一个单一的state
// state通过闭包的形式存放在redux store中，保证其是只读的。如果你想要更改state，只能通过发送action进行
// action本质上就是一个普通的对象
// 2.
// 你的应用可以通过redux暴露的subscribe方法，订阅state变化。
// 如果你在react应用中使用redux，则表现为react订阅store变化，并re-render视图。
// 3.
// 如何根据action来更新视图，这部分是业务相关的。 redux通过reducer来更新state

// reducer可以说是redux的精髓所在
// reducer应该是一个纯函数，这样state才可预测
////////////////////////////////////////////////////////////////////////////////////////
// middlewares

// store.dispatch = function dispatchAndLog(action) {
//   console.log('dispatching', action)
//   let result = next(action)
//   console.log('next state', store.getState())
//   return result
// }

// 加入中间件
// applyMiddleware函数开始，它主要的目的就是为了处理store的dispatch函数。
// return applyMiddleware(thunkMiddleware)(createStore)(reducer, preloadedState)
function applyMiddleware(...middlewares) {
  return createStore =>
    (reducer, ...args) => {
      const store = createStore(reducer, ...args)
      const middlewareAPI = {
        getState: store.getState,
        dispatch: (action, ...args) => dispatch(action, ...args),
      }

      const chain = middlewares.map(middleware => middleware(middlewareAPI))
      // 将middlewares组成一个函数
      // 也就是说就从前到后依次执行middlewares
      const dispatch = compose(...chain)(store.dispatch)
      return {
        ...store,
        dispatch,
      }
    }
}

// 使用
const s = createStore(
  todoApp,
  // applyMiddleware() tells createStore() how to handle middleware
  applyMiddleware(logger, dispatchAndLog)
)

// chain大概长这个样子：
// chain = [
//   function middleware1(next) {
//     // 内部可以通过闭包访问到getState和dispath

//   },
//   function middleware2(next) {
//     // 内部可以通过闭包访问到getState和dispath

//   },
//   ...
// ]
// 第二个中间件开始，next其实就是上一个中间件返回的 action => retureValue
// 这个函数签名就是dispatch的函数签名。

function compose(...funcs) {
  if (funcs.length === 0) {
    return arg => arg
  }

  if (funcs.length === 1) {
    return funcs[0]
  }

  return funcs.reduce(
    (a, b) =>
      (...args) =>
        a(b(...args))
  )
}
