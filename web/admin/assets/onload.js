((w, d, l) => {
    let loggedInUser = null, version = "dev", createAdmin = false;

    const pagePrefix = '/admin', pageSuffix = '.html', console = w.console, pageEvents = {
        '/': () => {
            apiRequest('depts/count').then(data => {
                data.count && (getElementById('dept_count').innerHTML = '(' + data.count + ')');
            });
            apiRequest('users/count').then(data => {
                data.count && (getElementById('user_count').innerHTML = '(' + data.count + ')');
            });
        },
        '/init': () => {
            const form = getElementById('init_form');
            form.addEventListener('submit', e => {
                postFormUrlEncoded('auth/init', form).then(data => {
                    console.debug(data);
                    if (data.message) {
                        alert(data.message);
                        navigateTo('/login');
                    }
                });
                e.preventDefault();
                return false;
            })
        },
        '/login': () => {
            const form = getElementById('login_form');
            form.addEventListener('submit', e => {
                postFormUrlEncoded('auth/login', form).then(data => {
                    console.debug(data);
                    if (data.username) {
                        navigateTo('/index');
                    }
                });
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
                    .then(() => passwordForm.reset());
                e.preventDefault();
                return false;
            });

            settingsForm.addEventListener('submit', e => {
                postFormUrlEncoded('auth/settings', settingsForm).then(showDataMessage);
                e.preventDefault();
                return false;
            });

            apiRequest('auth/settings').then(data => {
                data.session_ttl && (getElementById('session_ttl').value = data.session_ttl)
            });
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
                    .then(() => newDepartmentForm.reset());
                e.preventDefault();
                return false;
            });
        },
        '/users': () => {
            const newUserForm = getElementById('new_user_form'), departmentsSelect = getElementById('dept'),
                submitButton = getElementsBySelector('#new_user_form [type=submit]')[0],
                loadUserList = () => {
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
                            li.innerHTML = '<a href="/admin/users/details.html?user={user}">{user}</a> ({usage})'
                                .replace(/{user}/g, user.name).replace(/{usage}/g, user.usage);
                            userList.appendChild(li)
                        });
                    });
                };
            submitButton.disabled = true;
            loadUserList();
            loadDepartmentOptions(departmentsSelect, 1).finally(() => submitButton.disabled = false);
            newUserForm.addEventListener('submit', e => {
                postFormUrlEncoded('users/create', newUserForm).then(loadUserList)
                    .then(() => newUserForm.reset());
                e.preventDefault();
                return false;
            });
        },
        '/users/details': () => {
            const editUserForm = getElementById('edit_user_form'),
                originalNameInput = getElementById('orig_name'),
                nameInput = getElementById('name'),
                departmentsSelect = getElementById('dept'),
                submitButton = getElementsBySelector('#edit_user_form [type=submit]')[0],
                downloadHeader = getElementById('download_header'),
                downloadLink = getElementById('download_link'),
                user = getUrlParam('user');
            originalNameInput.value = nameInput.value = user;
            submitButton.disabled = true;
            downloadHeader.innerHTML = downloadHeader.innerHTML.replace(/{user}/g, user);
            downloadLink.addEventListener('focus', e => downloadLink.select());
            apiRequest('users/details', {queryParams: {user: user}})
                .then(data => {
                    console.debug(data);
                    if (!data.token) {
                        return;
                    }
                    downloadLink.value = l.origin + '/download/' + data.token;
                    loadDepartmentOptions(departmentsSelect, 0).finally(() => {
                        departmentsSelect.value = data.dept;
                        submitButton.disabled = false;
                    })
                });
            editUserForm.addEventListener('submit', e => {
                postFormUrlEncoded('users/edit', editUserForm).then(data => {
                    if (!data.user) {
                        return;
                    }
                    data.message && alert(data.message);
                    l.replace(l.origin + l.pathname + '?user=' + data.user);
                });
                e.preventDefault();
                return false;
            });
        },
        '/front-download': () => {
            const configDownload = getElementById('user_config'), keyDownload = getElementById('ssh_key'),
                regenerateButton = getElementById('regenerate'), token = l.pathname.split('/').pop();
            configDownload.href = l.origin + l.pathname + '/config';
            keyDownload.href = l.origin + l.pathname + '/key';
            apiRequest('users/has-key/' + token).then(data => {
                data.result && keyDownload.parentNode.classList.remove('hidden');
            });
            d.on('click', regenerateButton, e => {
                apiRequest('users/regenerate-key/' + token, {method: 'post'}).then(data => {
                    if (!data.new_token) {
                        return;
                    }
                    const newUrl = l.origin + '/download/' + data.new_token;
                    alert('! NOTICE !\nYour download URL has been changed to:\n' + newUrl);
                    l.replace(newUrl);
                });
            });
        },
        '/system': () => {
            const sshfsForm = getElementById('sshfs_form');
            apiRequest('system').then(data => {
                const config = data.config;
                config && Object.keys(config).map(k => {
                    getElementById(k) && (getElementById(k).value = config[k]);
                });
            });
            sshfsForm.addEventListener('submit', e => {
                postFormUrlEncoded('system/sshfs', sshfsForm).then(data => {
                    data.message && alert(data.message);
                });
                e.preventDefault();
                return false;
            });
        },
    }, publicPages = ['/init', '/login', '/logout', '/front-download'];

    function init() {
        const pageRoute = getPageRoute(), loginToAccess = getElementsBySelector('.login-to-access'),
            header = getElementsBySelector('h1')[0], versionSpan = newElement('span');
        console.debug(loggedInUser, pageRoute);
        if (!createAdmin && pageRoute === '/init') {
            navigateTo('/index');
            return;
        }
        if (createAdmin && pageRoute !== '/init') {
            navigateTo('/init');
            return;
        }
        if (needLogin(pageRoute)) {
            navigateTo('/login');
            return;
        }
        versionSpan.classList.add('version');
        versionSpan.innerHTML = version;
        header.appendChild(versionSpan);
        showWelcome();
        initTogglers();
        initRepeatPassword();
        listenPageEvents(pageRoute);
        loginToAccess.forEach(e => {
            e.classList.remove('login-to-access');
        });
    }

    /**
     * @returns {Promise}
     */
    function apiRequest(url, options) {
        if (options && options.queryParams) {
            url += (url.indexOf('?') === -1 ? '?' : '&') + buildQueryParams(options.queryParams);
            delete options.queryParams;
        }
        options = Object.assign({credentials: 'same-origin', cache: 'no-store'}, options || {});
        return fetch('/api/' + url, options).then(response => {
            return response.json();
        }).then(data => {
            if (data.hasOwnProperty('error')) {
                throw new Error(data.error)
            }
            return data;
        }).catch(error => {
            console.error(error);
            alert(error.message);
            return false
        });
    }

    function buildQueryParams(params) {
        const esc = encodeURIComponent;
        return Object.keys(params).map(k => esc(k) + '=' + esc(params[k])).join('&');
    }

    function checkLoginStatus() {
        return apiRequest('auth/status').then(data => {
            loggedInUser = data.username;
            version = data.version;
            createAdmin = !!data.create_admin;
        });
    }

    function getElementById(id) {
        return d.getElementById(id);
    }

    function getElementsBySelector(selector) {
        return d.querySelectorAll(selector);
    }

    function getPageRoute() {
        if (l.pathname.startsWith('/download/')) {
            return '/front-download';
        }
        return l.pathname.slice(pagePrefix.length, -pageSuffix.length) || '/';
    }

    function getUrlParam(name) {
        return (new URL(l.href)).searchParams.get(name);
    }

    function initRepeatPassword() {
        let reverseCheck = false, setValidity = (repeat, password) => {
            repeat.setCustomValidity(repeat.value === password.value ? '' : 'Password must be repeated exactly');
        };
        d.on('change', '[data-repeat-password]', e => {
            const repeat = e.target, password = getElementById(repeat.dataset.repeatPassword);
            setValidity(repeat, password);
            if (reverseCheck) {
                return;
            }
            d.on('change', password, e => {
                setValidity(repeat, password);
            });
        });
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

    /**
     * @param departmentsSelect {HTMLSelectElement}
     * @param initialOptionCount {int}
     * @returns {Promise}
     */
    function loadDepartmentOptions(departmentsSelect, initialOptionCount) {
        return apiRequest('depts/options').then(data => {
            console.debug(data);
            const options = data.options;
            if (!options) {
                return;
            }
            while (departmentsSelect.options.length > initialOptionCount) {
                departmentsSelect.remove(initialOptionCount);
            }
            for (const key of Object.keys(options)) {
                const option = newElement('option');
                option.value = key;
                option.text = options[key];
                departmentsSelect.add(option);
            }
        });
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

    /**
     * @returns {Promise}
     */
    function postFormUrlEncoded(url, form) {
        const data = new URLSearchParams(new FormData(form)), submitButton = form.querySelectorAll('[type=submit]');
        submitButton.forEach(e => e.disabled = true);
        return apiRequest(url, {
            method: 'post',
            body: data,
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
        }).finally(() => submitButton.forEach(e => e.disabled = false));
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
                if (target === selector || (typeof selector === 'string' && target.matches(selector))) {
                    event.target = target;
                    handler.call(target, event);
                    break;
                }
            }
        }, false);
    };

    checkLoginStatus().finally(init);
})(window, document, location);
