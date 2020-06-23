const state = () => ({
  user: null
});

const getters = {};

const actions = {};

// mutations
const mutations = {
  set(state, user) {
    state.user = user;
  },
  reset(state) {
    state.user = null;
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
};
