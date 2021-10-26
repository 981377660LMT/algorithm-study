https://laishuxin.github.io/blog_fe/cate/framework/react/7-react-redux/#%E5%BC%95%E4%BE%8B
https://zhuanlan.zhihu.com/p/358551518

1. reducer 规范 state 创建
   ```JS
   function reducer(state,{type,payload})
   ```
2. dispatch 规范 setState 流程
   ```JS
   function dispatch(action)
   ```
3. 将组件包裹写一个 HOC 使得 dispath 能够访问 state (react-redux) **此时每一个组件都有自己的 dispatch**
4. 提取函数 connect (接收一个组件) 一个经过包裹后的组件 连接组件与全局状态
5. 利用 connect 减少 render
   每次修改 state 导致没用到 state 的组件也重新渲染
   只要 setState 传的新的对象(非原来引用)浅比较 就会重新渲染
   解决:不能用 setState 干掉它
   换成自己的 store 里面有 state 和方法
   通过 Provider 传下去
   注意每次调用的自己 setState 后都要自己手动触发视图更新 即**随便用一个变量触发 setState 钩子** 保证 connect 组件才会被刷新
   问题是此时 update 只能实现一个组件的更新
   **其他需要订阅数据的变化**
   使用 useEffect 订阅
   subscribe 返回一个 subscription 可以取消订阅
   触发状态变化后调用回调
6. 此时的 redux:store,reducer,connect,appConext
7. connect 函数的 selector 选取 state
   selector 精准渲染:组件只在自己的数据变化时渲染
   如果 selector 后的 state 更新了(shallowDiff)则更新视图
   selector 变化了则取消订阅
8. **mapDispatcherToProps**
   connect(mapStateToProps,mapDispatcherToProps)
9. 为什么要这样设计 Connect?
   mapDispatcherToProps 封装写 mapStateToProps 封装读 易于维护
   **这个读写接口可以传给任何组件**
10. createStore(reducer,initState)
    state 和 reducer 不能写死 要用户传进来
11. redux 与 react-redux 思想总结
    react-redux 为我们实现的是组件和 store 连接（connect）以及两个读和写的 API。
    其中，connect 做的是对组件进行封装。

    1. 从 store 里面获取读和写的 API。（我们的代码直接从 store 上拿）
    2. 对拿到的接口进行封装。根据 mapStateToProps 和 mapDispatchToProps 进行封装。
    3. 在恰当的时候进行更新。只有 store 发生变化，才对页面进行更新。
    4. 渲染组件。

12. react-redux 概念
    1. store
    2. state
    3. dispatch
       - reducer
       - initState
       - action
         - type；action 的类型。
         - payload：action 携带的信息。
    4. connect
    5. Provider
    6. middlewares
13. 异步 action:让 action 支持函数和让 payload 支持 promise (递归做法)
    Redux 是不支持异步 action 的，这是为了确保 reducer 是纯函数。 为了让 Redux 支持异步 action 就出现了许多有名的中间件，下面就 介绍一下异步 action 的由来以及实现异步 action 的原理。
    不写 fetchUser(dispatch) 而是写 dispatch(fetchUser)
    包裹一下 dispatch 判断 action 是函数还是对象

    ```JS
    // react-redux.js
    let dispatch = action => setState(reducer(state, action))
    // 支持异步 action
    let prevDispatch = dispatch
    dispatch = action => {
      if (typeof action === 'function') {
        return action(dispatch)
      }
      return prevDispatch(action)
    }
    ```

    注意：这里不能使用 prevDispatch。这是因为我们不知道异步 action 内部是否还嵌套另一个 异步 action，所以采用递归 的方式确保最终拿到的 action 是我们想要的。
    事实上，上面的代码是 redux-thunk 的简化版。我们来看下 react-thunk 是实现（也就几行代码）：

    ```JS
    function createThunkMiddleware(extraArgument) {
      return ({ dispatch, getState }) => next => action => {
        if (typeof action === 'function') {
          // 关键是 action(dispatch) 其他参数无关紧要
          return action(dispatch, getState, extraArgument)
        }

        // next是preDispatch
        return next(action)
      }
    }

    const thunk = createThunkMiddleware()
    thunk.withExtraArgument = createThunkMiddleware

    export default thunk
    ```

14. redux 中间件 middlewares
    createStore(reducer,initState,**applyMiddlewares(...middlewares)**)
    如果每次要拓展 dispatch 都需要在源码层面进行修改，那显然是违反了开闭原则。 我们先把前面的实现异步 action 的代码删除，我们先实现中间件功能，等下再通过中间件 实现异步 action。
15. redux-thunk 和 redux-promise (递归 dispatch)
    createStore(reducer,initState,**applyMiddlewares(redux-thunk,redux-promise)**)
    **redux-thunk**:
    action 是个函数就调用 action 传入 dispatch
    否则直接 preDispatch(action)

    ```JS
    const fetchUser=dispatch=>...dispatch({type：'..',payload:..})
    ```

    不是函数就直接 preDispatch(action)

    ```JS
    dispatch({type：'..',payload:..})
    ```

    **redux-promise**
    payload 是个 promise 就调用 payload.then(dispatch({...action,payload:promise 的结果}))
    否则 preDispatch(action)

    ```JS

    export default function promiseMiddleware({ dispatch }) {
      return next => action => {
        if (!isFSA(action)) {
          return isPromise(action) ? action.then(dispatch) : next(action);
        }

        return isPromise(action.payload)
          ? action.payload
              .then(result => dispatch({ ...action, payload: result }))
              .catch(error => {
                dispatch({ ...action, payload: error, error: true });
                return Promise.reject(error);
              })
          : next(action);
      };
    }
    ```

16. 中间件写法
    ({dispatch})=>next(preDispatch)=>action=>{...}

    ```JS
    function createMiddleware() {
      return ({ dispatch }) => next => action => {
        if (typeof action === 'function') {
          // 关键是 action(dispatch) 其他参数无关紧要
          return action(dispatch)
        }

        // next是preDispatch
        return next(action)
      }
    }
    ```

**文档看了几天，不及这一时的视频**
