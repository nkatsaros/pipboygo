export const LOCAL_MAP_UPDATE = 'LOCAL_MAP_UPDATE'
export const DATABASE_UPDATE = 'DATABASE_UPDATE'

export function updateLocalMap(url) {
  return {
    type: LOCAL_MAP_UPDATE,
    url
  }
}

export function updateDatabase(database) {
  return {
    type: DATABASE_UPDATE,
    database
  }
}
