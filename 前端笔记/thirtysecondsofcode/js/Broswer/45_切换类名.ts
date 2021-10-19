const toggleClass = (el: HTMLElement, className: string) => el.classList.toggle(className)
toggleClass(document.querySelector('p.special')!, 'special')
// The paragraph will not have the 'special' class anymore
