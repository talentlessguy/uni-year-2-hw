const urlParams = new URLSearchParams(location.search)

const res = await fetch(
  `/api/buses?stop_name=${urlParams.get('name')}&stop_area=${
    urlParams.get('region')
  }`,
)

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
  btn.textContent = `${bus.route_short_name} (${bus.route_long_name})`
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
    if (coords.latitude) {
      fetch(
        `/api/nearest_stop?user_lat=${coords.latitude}&user_lon=${coords.longitude}&route_id=${bus.route_id}`,
      )
        .then((res) => res.json()).then((json) => {
          console.log(`Closest stop: `, json)

          console.log(`Bus: `, bus)

          // fetch(`/api/schedule?trip_id=${bus.trip_id}&stop_id=${json.id}`).then((
          //   res,
          // ) => res.json()).then((json) => {
          //   console.log(json)
          // })
        })
    }
  }

  grid.appendChild(btn)
}
