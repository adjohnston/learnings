const STATES = {
  lit: 'lit',
  unlit: 'unlit',
  broken: 'broken',
}

const lightbulb = () => {
  let state = STATES.unlit

  return {
    state() {
      return state
    },
    toggle() {
      switch (state) {
        case STATES.lit:
          state = STATES.unlit
          break

        case STATES.unlit:
          state = STATES.lit
          break
      }
    },
    break() {
      state = STATES.broken
    },
  }
}

const bulb = lightbulb()
