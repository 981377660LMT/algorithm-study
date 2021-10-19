const insertAfter = (el: HTMLElement, htmlString: string) =>
  el.insertAdjacentHTML('afterend', htmlString)
insertAfter(document.getElementById('myId')!, '<p>after</p>')
// <div id="myId">...</div> <p>after</p>
