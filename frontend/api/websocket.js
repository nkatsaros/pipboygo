import store from '../getStore'
import { updateLocalMap, updateDatabase } from '../actions'

let ws_uri

if (location.protocol === "https:") {
    ws_uri = "wss:"
} else {
    ws_uri = "ws:"
}

ws_uri += "//" + location.host
ws_uri += "/ws"

let socket = new WebSocket(ws_uri)

socket.onmessage = (event) => {
  if (event.data instanceof Blob) {
    store.dispatch(updateLocalMap(URL.createObjectURL(event.data)))
  } else {
    let data = JSON.parse(event.data)
    if (data.type === "db") {
      store.dispatch(updateDatabase(data.database))
    }
  }
}

export default socket
