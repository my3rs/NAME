import axios from 'axios';

axios.interceptors.use(
    config => {
        config.baseURL = '/api/v1/'
        config.withCredentials = true
        config.timeout = 6000

        let token = localStorage.getItem('jwt')
        let csrf = store.getters.csrf

        if (token) {
            console.log('Found token in localStorage: ', token)
            config.headers = {
                'Authorization': token,
            }
        }

        return config
    },
    error => {
        return Promise.reject(error)
    }
)

