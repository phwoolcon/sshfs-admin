((self, caches, l) => {
    const CACHE_NAME = '19.9.18', cacheList = ['/admin/', 'assets/onload.js', 'assets/style.css'];

    self.addEventListener('install', e => {
        e.waitUntil(
            caches.open(CACHE_NAME)
                .then(cache => cache.addAll(cacheList))
                .then(() => self.skipWaiting())
        );
    });
    self.addEventListener('fetch', e => {
        const request = e.request;
        if (request.url.startsWith(l.origin + '/api/')) {
            e.respondWith(fetch(request).catch(err => {
                return new Response(JSON.stringify({offline: true}))
            }));
            return
        }
        e.respondWith(caches.open(CACHE_NAME).then(cache => {
                return cache.match(request).then(response => {
                    return response || fetch(request)
                        .then(response => {
                            cache.put(request, response.clone());
                            return response
                        })
                })
            }
        ));
    });
    self.addEventListener('activate', event => {
        event.waitUntil(
            caches.keys().then(cacheNames => {
                console.log(cacheNames);
                return Promise.all(
                    cacheNames.map(cacheName => {
                        cacheName === CACHE_NAME && caches.delete(cacheName);
                    })
                );
            })
        );
    });
})(self, caches, location);
