var loginCB = (auth, token) => {
    $("#navbar").removeClass('hide');

    setAxiosAuthorizationToken(token);

    if (location.pathname == '' || location.pathname == "/" || location.pathname == '/#!' || location.pathname == '/#!/' || location.pathname == '/#!/login')
        router.go('users');
}

$(function () {
    initAfterLoad($('body'), null);

    axios.defaults.responseType = 'json';
    axios.defaults.headers.common['Content-Type'] = 'application/json';
    axios.defaults.headers.common.post = {};
    axios.defaults.headers.common.put = {};
    axios.defaults.xsrfCookieName = '_CSRF';
    axios.defaults.xsrfHeaderName = 'CSRF';
    axios.defaults.withCredentials = true;

    if (isAuthenTokenOK())
        setAxiosAuthorizationToken(getToken());

    router.setHashMode('#!');
    router.setDefaultCallback(initAfterLoad);
    router.config(state => {
        state
            .when('login', {
                templateUrl: 'www/dist/login.html',
                enableCache: false
            })
            .when('users', {
                templateUrl: 'www/dist/components/user/index.html',
                enableCache: false
            })
            .when('group', {
                templateUrl: 'www/dist/components/group/index.html',
                enableCache: false
            })
            .when('resource', {
                templateUrl: 'www/dist/components/resource/index.html',
                enableCache: false
            })
            .when('action', {
                templateUrl: 'www/dist/components/action/index.html',
                enableCache: false
            })
            .when('department', {
                templateUrl: 'www/dist/components/department/index.html',
                enableCache: false
            })
            .when('planOnGoing', {
                templateUrl: 'www/dist/components/planOnGoing/index.html',
                enableCache: false
            })
            .when('planDone', {
                templateUrl: 'www/dist/components/planDone/index.html',
                enableCache: false
            });
    });

    checkLogin(loginCB);
});