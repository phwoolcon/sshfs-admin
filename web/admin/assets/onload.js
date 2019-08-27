!function (w, d, l) {
    let loggedInUser = null;

    const app = w.app = {
        init: function () {
            const pageRoute = getPageRoute(), loginToAccess = getElementsBySelector('.login-to-access');
            console.debug(pageRoute);
            if (loginRequired(pageRoute)) {
                navigateTo('/login');
                return;
            }
            showWelcome();
            listenPageEvents(pageRoute);
            loginToAccess.forEach(function (e) {
                e.classList.remove('login-to-access');
            });
        }
    }, pagePrefix = '/admin', pageSuffix = '.html', skipLoginForPages = ['/login', '/logout'], console = w.console;

    function checkLoginStatus() {
        return fetch('/api/auth/status')
            .then(parseJSON)
            .then(processError)
            .then(function (data) {
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

    function loginRequired(route) {
        console.debug(loggedInUser);
        if (loggedInUser) {
            return false;
        }
        return skipLoginForPages.indexOf(route) < 0;
    }

    function listenPageEvents(pageRoute) {
        const pageEvents = {
            '/login': function () {
                const form = getElementById('login_form');
                form.addEventListener('submit', function (e) {
                    postFormUrlEncoded('/api/auth/login', form)
                        .then(parseJSON)
                        .then(processError)
                        .then(function (data) {
                            console.debug(data);
                            if (data.username) {
                                navigateTo('/index');
                            }
                        })
                        .catch(function (error) {
                            console.error(error);
                            alert(error.message);
                        });
                    e.preventDefault();
                    return false;
                })
            },
            '/account': function () {
                const passwordForm = getElementById('change_password_form'),
                    settingsForm = getElementById('account_settings_form');
                passwordForm.addEventListener('submit', function (e) {
                    postFormUrlEncoded('/api/auth/change-pass', passwordForm)
                        .then(parseJSON)
                        .then(processError)
                        .then(function (data) {
                            console.debug(data);
                            if (!data.message) {
                                throw new Error('Something went wrong')
                            }
                            passwordForm.reset();
                            alert(data.message);
                        })
                        .catch(function (error) {
                            console.error(error);
                            alert(error.message);
                        });
                    e.preventDefault();
                    return false;
                });

                settingsForm.addEventListener('submit', function (e) {
                    postFormUrlEncoded('/api/auth/settings', settingsForm)
                        .then(parseJSON)
                        .then(processError)
                        .then(function (data) {
                            console.debug(data);
                            if (!data.message) {
                                throw new Error('Something went wrong')
                            }
                            alert(data.message);
                        })
                        .catch(function (error) {
                            console.error(error);
                            alert(error.message);
                        });
                    e.preventDefault();
                    return false;
                });

                fetch('/api/auth/settings')
                    .then(parseJSON)
                    .then(processError)
                    .then(function (data) {
                        data.session_ttl && (getElementById('session_ttl').value = data.session_ttl)
                    })
                    .catch(function (error) {
                        console.error(error);
                        alert(error.message);
                    });
            },
        };
        if (pageEvents.hasOwnProperty(pageRoute)) {
            pageEvents[pageRoute]();
        }
    }

    function navigateTo(route) {
        l.href = pageUrl(route);
    }

    function pageUrl(route) {
        return pagePrefix + (route.charAt(0) === '/' ? '' : '/') + route + pageSuffix;
    }

    function postFormUrlEncoded(url, form) {
        const data = new URLSearchParams(new FormData(form));
        return fetch(url, {
            method: 'post',
            body: data,
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
        });
    }

    /**
     *
     * @param {Response} response
     * @returns {any}
     */
    function parseJSON(response) {
        return response.json();
    }

    function processError(data) {
        if (data.error) {
            throw new Error(data.error)
        }
        return data;
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
                fetch('/api/auth/logout');
                navigateTo('/login');
            });
        }
    }

    checkLoginStatus().finally(function () {
        app.init();
    });
}(window, document, location);
