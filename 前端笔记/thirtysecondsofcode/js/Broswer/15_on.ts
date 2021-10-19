export {}
const fn = () => console.log('!')
on(document.body, 'click', fn) // logs '!' upon clicking the body
on(document.body, 'click', fn, { target: 'p' })
// logs '!' upon clicking a `p` element child of the body
on(document.body, 'click', fn, { capture: true })
// use capturing instead of bubbling

interface Options {
  target?: string
  capture?: boolean
}

function on(el: HTMLElement, event: string, fn: EventListener, options: Options = {}) {
  const delegationFunction = (e: Event) =>
    options.target && (e.target as Element).matches(options.target) && fn.call(e.target, e) // matches确保事件目标与指定的目标匹配
  el.addEventListener(event, options.target ? delegationFunction : fn, options.capture)
}
