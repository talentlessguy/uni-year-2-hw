const url = document.querySelector('meta[name="current-url"]').content
const res = await fetch(url)
const stops = await res.json()

const input = document.getElementById('input')
const results = document.getElementById('results')

input.addEventListener('input', () => {
  const inputValue = input.value
  results.innerHTML = ''
  if (inputValue) {
    const filteredSuggestions = stops.filter((suggestion) =>
      suggestion.toLowerCase().includes(inputValue.toLowerCase())
    )

    filteredSuggestions.forEach((suggestion) => {
      const div = document.createElement('div')
      div.textContent = suggestion
      div.classList.add('px-3', 'py-2', 'cursor-pointer', 'hover:bg-gray-100')
      div.addEventListener('click', () => {
        input.value = suggestion
        results.innerHTML = ''
        results.classList.add('hidden')
      })
      results.appendChild(div)
    })

    results.classList.toggle('hidden', !filteredSuggestions.length)
  } else {
    results.classList.add('hidden')
  }
})

input.addEventListener('click', () => {
  console.log('trigger')
  stops.forEach((suggestion) => {
    const div = document.createElement('div')
    div.textContent = suggestion
    div.classList.add('px-3', 'py-2', 'cursor-pointer', 'hover:bg-gray-100')
    div.addEventListener('click', () => {
      input.value = suggestion
      results.innerHTML = ''
      results.classList.add('hidden')
    })
    results.appendChild(div)
    results.classList.toggle('hidden', false)
  })
})

document.addEventListener('click', (event) => {
  if (event.target !== input) {
    results.classList.add('hidden')
  }
})