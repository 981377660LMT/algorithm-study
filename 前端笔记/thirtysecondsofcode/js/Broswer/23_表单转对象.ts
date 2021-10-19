console.log(formToObject(document.querySelector('#form')!))
// { email: 'test@email.com', name: 'Test Name' }

function formToObject(form: HTMLFormElement) {
  return (
    // @ts-ignore
    Array.from(new FormData(form)).reduce(
      (acc, [key, value]) => ({
        ...acc,
        [key]: value,
      }),
      {}
    )
  )
}
