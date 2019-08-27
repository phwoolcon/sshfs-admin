!function (w, d, l) {
    let loggedInUser = null;

    const pageEvents = {
        '/login': function () {
            const form = getElementById('login_form');
            form.addEventListener('submit', function (e) {
                postFormUrlEncoded('auth/login', form).then(function (data) {
                    console.debug(data);
                    if (data.username) {
                        navigateTo('/index');
                    }
                }).catch(catchAlert);
                e.preventDefault();
                return false;
            })
        },
        '/account': function () {
            const passwordForm = getElementById('change_password_form'),
                settingsForm = getElementById('account_settings_form'),
                showDataMessage = function (data) {
                    console.debug(data);
                    if (!data.message) {
                        throw new Error('Something went wrong')
                    }
                    alert(data.message);
                };
            passwordForm.addEventListener('submit', function (e) {
                postFormUrlEncoded('auth/change-pass', passwordForm).then(showDataMessage)
                    .then(passwordForm.reset).catch(catchAlert);
                e.preventDefault();
                return false;
            });

            settingsForm.addEventListener('submit', function (e) {
                postFormUrlEncoded('auth/settings', settingsForm).then(showDataMessage).catch(catchAlert);
                e.preventDefault();
                return false;
            });

            request('auth/settings').then(function (data) {
                data.session_ttl && (getElementById('session_ttl').value = data.session_ttl)
            }).catch(catchAlert);
        },
    }, pagePrefix = '/admin', pageSuffix = '.html', publicPages = ['/login', '/logout'], console = w.console;

    function init() {
        const pageRoute = getPageRoute(), loginToAccess = getElementsBySelector('.login-to-access');
        console.debug(loggedInUser, pageRoute);
        if (needLogin(pageRoute)) {
            navigateTo('/login');
            return;
        }
        showWelcome();
        listenPageEvents(pageRoute);
        loginToAccess.forEach(function (e) {
            e.classList.remove('login-to-access');
        });
    }

    function catchAlert(error) {
        console.error(error);
        alert(error.message);
    }

    function checkLoginStatus() {
        return request('auth/status').then(function (data) {
            loggedInUser = data.username;
        });
    }

    function getElementById(id) {
        return d.getElementById(id);
    }

    function getElementsBySelector(selector) {
        return d.querySelectorAll(selector);
    }

    function getPageRoute() {
        return l.pathname.slice(pagePrefix.length, -pageSuffix.length) || '/';
    }

    function listenPageEvents(pageRoute) {
        if (pageEvents.hasOwnProperty(pageRoute)) {
            pageEvents[pageRoute]();
        }
    }

    function navigateTo(route) {
        l.href = pageUrl(route);
    }

    function needLogin(route) {
        if (loggedInUser) {
            return false;
        }
        return publicPages.indexOf(route) < 0;
    }

    function pageUrl(route) {
        return pagePrefix + (route.charAt(0) === '/' ? '' : '/') + route + pageSuffix;
    }

    function postFormUrlEncoded(url, form) {
        const data = new URLSearchParams(new FormData(form));
        return request(url, {
            method: 'post',
            body: data,
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
        });
    }

    /**
     * @returns {Promise}
     */
    function request(url, options) {
        options = Object.assign({credentials: 'same-origin', cache: 'no-store'}, options || {});
        return fetch('/api/' + url, options).then(function (response) {
            return response.json();
        }).then(function (data) {
            if (data.error) {
                throw new Error(data.error)
            }
            return data;
        });
    }

    function showWelcome() {
        let template;
        const container = getElementById('welcome');
        if (!container) {
            return;
        }
        if (loggedInUser) {
            template = 'Welcome, {user} | ' +
                '<a href="{account_url}">Account</a> | <a id="logout" href="javascript:">Logout</a>';
            container.innerHTML = template.replace('{user}', loggedInUser)
                .replace('{account_url}', pageUrl('/account'));
            getElementById('logout').addEventListener('click', function () {
                request('auth/logout').then(function () {
                    navigateTo('/login');
                });
            });
        }
    }

    checkLoginStatus().finally(init);
}(window, document, location);
