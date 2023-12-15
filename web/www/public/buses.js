const urlParams = new URLSearchParams(location.search)

const res = await fetch(`/api/buses?stop=${urlParams.get('value')}&mode=text`)

const json = await res.json()

const grid = document.getElementById('grid')

/**
 * @type {GeolocationCoordinates}
 */
let coords = {}

navigator.geolocation.getCurrentPosition((pos) => {
  coords = pos.coords
}, (err) => {
  console.error(err)
})

for (const bus of json) {
  const btn = document.createElement('button')
  btn.textContent = `${bus.route_short_name} (${bus.trip_long_name})`
  btn.classList.add(
    'bg-green-500',
    'hover:bg-green-700',
    'text-white',
    'font-bold',
    'py-1',
    'px-2',
    'text-sm',
    'rounded',
  )

  btn.onclick = () => {
    fetch(
      `/api/arrivals?user_lat=${coords.latitude}&user_lon=${coords.longitude}&route_id=${bus.route_id}`,
    )
      .then((res) => res.json()).then((json) => {
        console.log(`Closest stop: `, json)

        console.log(bus)
      })
  }

  grid.appendChild(btn)
}
