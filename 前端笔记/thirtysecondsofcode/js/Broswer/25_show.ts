show(...document.querySelectorAll('img'))

function show(...eles: HTMLElement[]) {
  eles.forEach(ele => (ele.style.display = ''))
  // eles.forEach(ele => (ele.style.display = 'none'))
}
// Shows all <img> elements on the page
