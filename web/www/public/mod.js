import {
  html,
  render,
  useEffect,
  useState,
} from 'https://unpkg.com/htm@3.1.1/preact/standalone.module.js'

const url = document.querySelector('meta[name="current-url"]').content
const res = await fetch(url)
const stops = await res.json()
const entries = normalizeEntries(stops)

function normalizeEntries(stops) {
  if (typeof stops[0] === 'string') {
    return stops.map((s) => ({ name: s, value: s }))
  } else return stops
}

const Dropdown = () => {
  const [inputValue, setInputValue] = useState('')
  const [suggestions, setSuggestions] = useState([])
  const [showResults, setShowResults] = useState(false)
  const [data, setData] = useState({})

  const handleInput = (e) => {
    const value = e.target.value
    setInputValue(value)

    const filteredSuggestions = entries.filter((suggestion) =>
      (typeof suggestion === 'string' ? suggestion : suggestion.name)
        .toLowerCase().includes(value.toLowerCase())
    )

    setSuggestions(filteredSuggestions)
    setShowResults(filteredSuggestions.length > 0)
  }

  const handleItemClick = (name) => {
    setInputValue(name)
    setData(entries.find((entry) => entry.name === name))
    setSuggestions([])
    setShowResults(false)
  }

  useEffect(() => {
    document.addEventListener('click', () => {
      setShowResults(false)
    })

    return () => {
      document.removeEventListener('click', () => {
        setShowResults(false)
      })
    }
  }, [])

  if (location.pathname === '/stops') {
    return html`
    <div>
      <input
        id="input"
        type="text"
        placeholder="Type to search..."
        name="stop_name"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:border-blue-300"
        value=${inputValue}
        onInput=${handleInput}
        onClick=${() => setShowResults(true)}
      />
      <input name="stop_area" value=${data.area} type="hidden" />
      <input name="stop_id" value=${data.id} type="hidden" />
      <div id="results" class=${`absolute w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg ${
      showResults ? '' : 'hidden'
    }`}>
        ${
      suggestions.map((suggestion) =>
        html`
          <div
            key=${suggestion.name}
            class="px-3 py-2 cursor-pointer hover:bg-gray-100"
            onClick=${() => handleItemClick(suggestion.name)}
          >
            ${suggestion.name}
          </div>
        `
      )
    }
      </div>
    </div>
  `
  } else {
    return html`
    <div>
      <input
        id="input"
        type="text"
        placeholder="Type to search..."
        name="value"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:border-blue-300"
        value=${inputValue}
        onInput=${handleInput}
        onClick=${() => setShowResults(true)}
      />
      <div id="results" class=${`absolute w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg ${
      showResults ? '' : 'hidden'
    }`}>
        ${
      suggestions.map((suggestion) =>
        html`
          <div
            key=${suggestion.name}
            class="px-3 py-2 cursor-pointer hover:bg-gray-100"
            onClick=${() => handleItemClick(suggestion.name)}
          >
            ${suggestion.name}
          </div>
        `
      )
    }
      </div>
    </div>
  `
  }
}

// Render Preact Dropdown component into the designated div
render(html`<${Dropdown} />`, document.getElementById('root'))
