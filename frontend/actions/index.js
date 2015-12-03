export const LOCAL_MAP_UPDATE = 'LOCAL_MAP_UPDATE'
export const MEMORY_UPDATE = 'MEMORY_UPDATE'

export function updateLocalMap(buffer) {
  return {
    type: LOCAL_MAP_UPDATE,
    buffer
  }
}

export function updateMemory(memory) {
  return {
    type: MEMORY_UPDATE,
    memory
  }
}
