const getDefaultState = () => {
  return {
    darkTheme: false,
  };
};

const state = getDefaultState();

const getters = {
  getDarkTheme(state) {
    return state.darkTheme;
  },
};

const mutations = {
  setDarkTheme(state, bool) {
    state.darkTheme = bool == "true";
  },
};

const actions = {
  toggleDarkTheme({ commit, state }) {
    const bool = (!state.darkTheme).toString();
    window.go.surge.MiddlewareFunctions.WriteSetting("DarkMode", bool).then(() => {
      commit("setDarkTheme", bool);
    });
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
