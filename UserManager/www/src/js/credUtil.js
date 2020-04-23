// { UserID, Group, Role, exp }
var cacheAuth = null;

var saveCreds = (cred) => {
    return new Promise((resolve, reject) => {
        try {
            // Try to decode jwt token
            const auth = jwt_decode(cred);

            // log decoded
            storage.setItem('auth', JSON.stringify(auth));
            storage.setItem('token', cred);

            // fullfill now
            resolve(auth);
        } catch (err) {
            reject(err);
        }
    });
}

var removeCreds = () => {
    storage.removeItem('auth');
};

var getCreds = () => {
    var auth = storage.getItem('auth');
    if (auth)
        cacheAuth = JSON.parse(auth);
    else
        cacheAuth = null;
    return cacheAuth;
}

var getToken = () => {
    return storage.getItem('token');
}

var isAuthenTokenOK = () => {
    let auth = getCreds();
    return auth && auth.exp * 1000 > Date.now();
}

var checkLogin = (cb) => {
    let auth = getCreds();

    if (!isAuthenTokenOK()) {
        storage.removeItem('auth');
        router.go('login');
    }
    else
        cb(auth, getToken());
}