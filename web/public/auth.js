import axios from 'axios';

function getLocalToken() {
    const token = window.localStorage.getItem('jwt');
}



const instance = axios.create({
    baseURI: '/api/v1',
    timeout: 300000,
    headers: {
        'Authorization': getLocalToken(),
    }
})

instance.setToken = (token) => {
    instance.defaults.headers['authorization'] = token;
    window.localStorage.setItem('jwt', token);
}

instance.interceptors.request.use(function (config) {
    let token = localStorage.getItem("jwt");

    if (token) {
        config.headers.Authorization = `${token}`;
        console.log(`Token in localStorage: ${token}`);
    } else {
        console.error("No token in localStorage");
        window.location = "/user/login";
    }

    return config;
}, function (error) {
    return Promise.reject(error);
});


instance.interceptors.response.use(function (response) {
    console.table(response.headers);
    console.log(response.data);

    const code = response.data.Code;
    if (code === 2) {   // token 过期
        refreshToken().then( function (response) {
            const { token } = response.headers['authorization'];
            instance.setToken(token);

            const config = response.config;
            // 重置一下配置
            config.headers['authorization'] = token;
            config.baseURL = '';

            return instance(config);
        }).catch(res => {
            console.log('refresh token error =>', res)
            window.location = '/user/login';
        })
    }

    return response;
}, function(error) {

    return Promise.reject(error);
});

function refreshToken() {
    return instance.post('/auth/refresh_token')
}