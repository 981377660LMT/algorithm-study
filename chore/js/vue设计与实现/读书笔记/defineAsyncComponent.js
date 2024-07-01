function defineAsyncComponent(options) {
  if (isFunction(options)) {
    options = { loader: options }
  }

  const {
    loader,
    loadingComponent,
    errorComponent,
    delay = 200,
    timeout, // undefined
    suspensible = true,
    onError: userOnError
  } = options

  let pendingRequest = null
  let resolvedComp = null
  let retries = 0

  const retry = () => {
    retries++
    pendingRequest = null
    return load()
  }

  const load = () => {
    // already resolved
    if (resolvedComp) {
      return resolvedComp
    }

    // in flight
    if (pendingRequest) {
      return pendingRequest
    }

    // use the loader directly
    const promise = (pendingRequest = loader())

    pendingRequest
      .then(comp => {
        if (comp.__esModule || comp[Symbol.toStringTag] === 'Module') {
          comp = comp.default
        }
        resolvedComp = comp
      })
      .catch(err => {
        error = err
      })
      .then(() => {
        pendingRequest = null
      })

    return promise
  }

  return defineComponent({
    name: 'AsyncComponentWrapper',

    __asyncLoader: load,

    get __asyncResolved() {
      return resolvedComp
    },

    data() {
      if (suspensible && resolvedComp) {
        return { ResolvedComponent: resolvedComp }
      } else {
        return { ResolvedComponent: null }
      }
    },

    mounted() {
      this.$nextTick(() => {
        this.__load()
      })
    },

    updated() {
      this.__load()
    },

    methods: {
      __load() {
        if (!resolvedComp) {
          const p = this.__load.__promise || (this.__load.__promise = load())
          if (p) {
            p.then(() => {
              this.ResolvedComponent = resolvedComp
            })
          }
        }
      }
    },

    render() {
      const { ResolvedComponent } = this
      if (ResolvedComponent) {
        return h(ResolvedComponent)
      } else {
        return h('div', null, 'loading...')
      }
    }
  })
}
