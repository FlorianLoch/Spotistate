import axios from "axios"

const CSRF_HEADER_NAME = "cassette_csrf_token"
const API_PATH = "/api"
const URL_DATA = API_PATH + "/you"
const URL_CSRF_TOKEN = API_PATH + "/csrfToken"
const URL_PLAYER_STATES = API_PATH + "/playerStates"
const URL_ACTIVE_DEVICES = API_PATH + "/activeDevices"


const API = function () {
  const client = axios.create()

  this.fetchCSRFToken = () => {
    return client.head(URL_CSRF_TOKEN).then((res) => {
      return res.headers[CSRF_HEADER_NAME]
    })
  }

  this.setCSRFToken = (csrfToken) => {
    client.defaults.headers.common[CSRF_HEADER_NAME] = csrfToken
  }

  this.fetchActiveDevices = () => {
    return client.get(URL_ACTIVE_DEVICES).then((res) => {
      return res.data
    })
  }

  this.fetchPlayerStates = () => {
    return client.get(URL_PLAYER_STATES).then((res) => {
      return res.data.states
    })
  }

  this.updatePlayerState = (slotNumber) => {
    return client.put(`${URL_PLAYER_STATES}/${slotNumber}`)
  }

  this.storePlayerState = () => {
    return client.post(URL_PLAYER_STATES)
  }

  this.deletePlayerState = (slotNumber) => {
    return client.delete(`${URL_PLAYER_STATES}/${slotNumber}`)
  }

  this.restoreFromPlayerState = (slotNumber, deviceID) => {
    const url = `${URL_PLAYER_STATES}/${slotNumber}/restore${(deviceID) ? `?deviceID=${deviceID}` : ""}`
    return client.post(url)
  }

  this.deleteYourData = () => {
    return client.delete(URL_DATA)
  }

  this.URL_DATA = URL_DATA
}

API.install = function (Vue) {
  Vue.prototype.$api = new API();
};

export default API