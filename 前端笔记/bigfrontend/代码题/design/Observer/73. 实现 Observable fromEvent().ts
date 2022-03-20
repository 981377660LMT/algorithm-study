import { Observable } from './57. 实现Observable'

/**
 * @param {HTMLElement} element
 * @param {string} eventName
 * @param {boolean} capture
 * @return {Observable}
 * 创建一个Observable，并传递DOM事件。
 */
function fromEvent(element: HTMLElement, eventName: string, capture = false): Observable {
  // your code here
  return new Observable(subscriber => {
    element.addEventListener(eventName, e => subscriber.next(e), capture)
  })
}

const source = fromEvent(document.createElement('button'), 'click')
source.subscribe(e => console.log(e))
