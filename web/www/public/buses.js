const urlParams = new URLSearchParams(location.search)

const grid = document.getElementById('grid')

const date = new Date()

const userTime = Intl.DateTimeFormat('et-EE', {
  hour: '2-digit',
  minute: '2-digit',
}).format(date)

document.getElementById('current_time').textContent += userTime

grid.innerHTML = 'loading'

navigator.geolocation.getCurrentPosition(async ({ coords }) => {
  const res = await fetch(
    `/api/buses?stop_name=${urlParams.get('stop_name')}&stop_area=${
      urlParams.get('stop_area')
    }&user_lat=${coords.latitude}&user_lon=${coords.longitude}&user_time=${userTime}`,
  )
  const json = await res.json()

  grid.innerHTML = ''

  document.getElementById('desc').textContent =
    `Closest stop: ${json.closest_stop.name}`

  if (!json.buses) {
    grid.innerHTML = 'No buses :('
    return
  }

  for (const bus of json.buses) {
    const div = document.createElement('div')
    div.innerHTML =
      `${bus.trip_long_name}<br />Bus code: ${bus.route_short_name}<br />Departure: ${bus.departure}`
    div.classList.add(
      'bg-white',
      'text-black',
      'border-solid',
      'border-2',
      'border-black',
      'py-1',
      'px-2',
      'text-md',
      'rounded',
    )

    grid.appendChild(div)
  }
}, (err) => {
  console.error(err)
})
