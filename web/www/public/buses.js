const urlParams = new URLSearchParams(location.search)

const res = await fetch(`/api/buses?stop=${urlParams.get('value')}&mode=text`)

const json = await res.json()

const grid = document.getElementById('grid')

for (const bus of json) {
  const btn = document.createElement('button')
  btn.textContent = `${bus.route_short_name} (${bus.trip_long_name})`
  btn.classList.add('bg-green-500', 'hover:bg-green-700' ,'text-white' ,'font-bold' ,'py-1' ,'px-2','text-sm' ,'rounded')

  grid.appendChild(btn)
}