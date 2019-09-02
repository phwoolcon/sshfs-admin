((w, d, l) => {
    let loggedInUser = null;

    const pageEvents = {
        '/': () => {
            apiRequest('depts/count').then(data => {
                data.count && (getElementById('dept_count').innerHTML = '(' + data.count + ')');
            });
            apiRequest('users/count').then(data => {
                data.count && (getElementById('user_count').innerHTML = '(' + data.count + ')');
            });
        },
        '/login': () => {
            const form = getElementById('login_form');
            form.addEventListener('submit', e => {
                postFormUrlEncoded('auth/login', form).then(data => {
                    console.debug(data);
                    if (data.username) {
                        navigateTo('/index');
                    }
                }).catch(catchAlert);
                e.preventDefault();
                return false;
            })
        },
        '/account': () => {
            const passwordForm = getElementById('change_password_form'),
                settingsForm = getElementById('account_settings_form'),
                showDataMessage = data => {
                    console.debug(data);
                    if (!data.message) {
                        throw new Error('Something went wrong')
                    }
                    alert(data.message);
                };
            passwordForm.addEventListener('submit', e => {
                postFormUrlEncoded('auth/change-pass', passwordForm).then(showDataMessage)
                    .then(() => passwordForm.reset()).catch(catchAlert);
                e.preventDefault();
                return false;
            });

            settingsForm.addEventListener('submit', e => {
                postFormUrlEncoded('auth/settings', settingsForm).then(showDataMessage).catch(catchAlert);
                e.preventDefault();
                return false;
            });

            apiRequest('auth/settings').then(data => {
                data.session_ttl && (getElementById('session_ttl').value = data.session_ttl)
            }).catch(catchAlert);
        },
        '/depts': () => {
            const loadDepartmentList = () => {
                const deptList = getElementById('depts');
                apiRequest('depts').then(data => {
                    let child;
                    console.debug(data);
                    if (!data.depts) {
                        throw new Error('Something went wrong')
                    }
                    while (child = deptList.firstChild) {
                        deptList.removeChild(child)
                    }
                    data.depts.forEach(dept => {
                        const li = newElement('li'), listHtml = [];
                        listHtml.push('<a href="javascript:" data-dept="{dept}">{dept}</a>'.replace(/{dept}/g, dept.name));
                        listHtml.push('({usage})'.replace(/{usage}/g, dept.usage));
                        // TODO Add user list
                        li.innerHTML = listHtml.join("\n");
                        deptList.appendChild(li)
                    });
                });
            }, newDepartmentForm = getElementById('new_dept_form');
            loadDepartmentList();
            newDepartmentForm.addEventListener('submit', e => {
                postFormUrlEncoded('depts/create', newDepartmentForm).then(loadDepartmentList)
                    .then(() => newDepartmentForm.reset()).catch(catchAlert);
                e.preventDefault();
                return false;
            });
        },
        '/users': () => {
            const userList = getElementById('users');
            apiRequest('users').then(data => {
                let child;
                console.debug(data);
                if (!data.users) {
                    throw new Error('Something went wrong')
                }
                while (child = userList.firstChild) {
                    userList.removeChild(child)
                }
                data.users.forEach(user => {
                    const li = newElement('li');
                    li.innerHTML = '<a href="">{user}</a>'.replace('{user}', user);
                    userList.appendChild(li)
                });
            })
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
        initTogglers();
        listenPageEvents(pageRoute);
        loginToAccess.forEach(e => {
            e.classList.remove('login-to-access');
        });
    }

    /**
     * @returns {Promise}
     */
    function apiRequest(url, options) {
        options = Object.assign({credentials: 'same-origin', cache: 'no-store'}, options || {});
        return fetch('/api/' + url, options).then(response => {
            return response.json();
        }).then(data => {
            if (data.hasOwnProperty('error')) {
                throw new Error(data.error)
            }
            return data;
        });
    }

    function catchAlert(error) {
        console.error(error);
        alert(error.message);
    }

    function checkLoginStatus() {
        return apiRequest('auth/status').then(data => {
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

    function initTogglers() {
        d.on('click', '[data-toggle-target]', e => {
            const toggler = e.target, target = getElementById(toggler.dataset.toggleTarget),
                hidden = 'hidden', nextAction = target.classList.contains(hidden) ? 'remove' : 'add';
            if (!target) {
                return;
            }
            target.classList[nextAction](hidden);
        });
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

    function newElement(tagName, options) {
        return d.createElement(tagName, options);
    }

    function pageUrl(route) {
        return pagePrefix + (route.charAt(0) === '/' ? '' : '/') + (route ? route + pageSuffix : '');
    }

    function postFormUrlEncoded(url, form) {
        const data = new URLSearchParams(new FormData(form));
        return apiRequest(url, {
            method: 'post',
            body: data,
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
        });
    }

    function putElementAfter(element, afterMe) {
        afterMe.parentNode.insertBefore(element, afterMe.nextSibling);
    }

    function showBreadcrumbAfterWelcome(welcome) {
        const breadcrumb = newElement('ol'), routes = ('home' + getPageRoute()).split('/').filter(route => {
            return route.length > 0;
        }), current = routes.pop(), currentLi = newElement('li'), linkStack = [], routeLabel = route => {
            return route.charAt(0).toUpperCase() + route.slice(1);
        };

        routes.forEach(route => {
            if (!route.length) {
                return;
            }
            linkStack.push(route === 'home' ? '' : route);
            const li = newElement('li');
            li.classList.add('link');
            li.innerHTML = '<a href="' + pageUrl(linkStack.join('/')) + '">' + routeLabel(route) + '</a>';
            breadcrumb.appendChild(li);
        });
        currentLi.innerHTML = routeLabel(current);
        breadcrumb.appendChild(currentLi);
        breadcrumb.id = 'breadcrumb';
        putElementAfter(breadcrumb, welcome);
    }

    function showWelcome() {
        let template;
        const container = getElementById('welcome');
        if (!container) {
            return;
        }
        showBreadcrumbAfterWelcome(container);
        if (loggedInUser) {
            template = 'Welcome, {user} | ' +
                '<a href="{account_url}">Account</a> | <a id="logout" href="javascript:">Logout</a>';
            container.innerHTML = template.replace('{user}', loggedInUser)
                .replace('{account_url}', pageUrl('/account'));
            d.on('click', '#logout', () => {
                apiRequest('auth/logout').then(() => {
                    navigateTo('/login');
                });
            });
        }
    }

    d.on = (eventName, selector, handler) => {
        const click = 'ontouchstart' in d.documentElement ? 'touchend' : 'click';
        eventName === 'click' && (eventName = click);
        d.addEventListener(eventName, event => {
            for (let target = event.target; target && target !== d; target = target.parentNode) {
                // loop parent nodes from the target to the delegation node
                if (target.matches(selector)) {
                    event.target = target;
                    handler.call(target, event);
                    break;
                }
            }
        }, false);
    };

    checkLoginStatus().finally(init);
})(window, document, location);
