export const LOCAL_MAP_UPDATE = 'LOCAL_MAP_UPDATE'
export const MEMORY_UPDATE = 'MEMORY_UPDATE'

export function updateLocalMap(url) {
  return {
    type: LOCAL_MAP_UPDATE,
    url
  }
}

export function updateMemory(memory) {
  return {
    type: MEMORY_UPDATE,
    memory
  }
}
