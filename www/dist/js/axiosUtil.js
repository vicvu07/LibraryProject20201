var setAxiosAuthorizationToken = (token) => {
    if (token) {
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
        delete axios.defaults.headers.common['Authorization'];
    }
}

var parseError = (error) => {
    if (error.response) {
        // The request was made, but the server responded with a status code
        // that falls out of the range of 2xx
        return {
            status: error.response.status,
            data: error.response.data
        }
    } else {
        // Something happened in setting up the request that triggered an Error
        return {
            status: -1,
            data: error
        }
    }
}

var post = (url, data) => {
    return new Promise((resolve, reject) => {
        axios.post(url, data).then(res => {
            if (res.status == 401 || res.data.Error.Code == 401 || res.status == 403 || res.data.Error.Code == 403) { // unauthorized
                if (confirm('Bạn không được phân quyền để thực hiện chức năng này')) {
                    if (res.status == 401 || res.data.Error.Code == 401)
                        router.go('login'); // go login
                }
                return
            }

            if (res.status != 200) {
                reject({ status: res.status, data: 'Server return with status: ' + res.status });
                return
            }

            if (res.data.Error.Code !== 200) {
                reject({ status: res.data.Error.Code, data: res.data.Error.Message });
            } else {
                resolve(res.data);
            }
        }).catch(err => {
            reject(parseError(err));
        });
    });
}

var get = (url, param) => {
    return new Promise((resolve, reject) => {
        axios.get(url, { params: param }).then(res => {
            if (res.status == 401 || res.data.Error.Code == 401 || res.status == 403 || res.data.Error.Code == 403) { // unauthorized
                if (confirm('Bạn không được phân quyền để thực hiện chức năng này')) {
                    if (res.status == 401 || res.data.Error.Code == 401)
                        router.go('login'); // go login
                }
                return
            }

            if (res.status != 200) {
                reject({ status: res.status, data: 'Server return with status: ' + res.status });
                return
            }

            if (res.data.Error.Code !== 200) {
                reject({ status: res.data.Error.Code, data: res.data.Error.Message });
            } else {
                resolve(res.data);
            }
        }).catch(err => {
            reject(parseError(err));
        });
    });
}