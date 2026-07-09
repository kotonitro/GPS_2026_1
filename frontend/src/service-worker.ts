/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />

import { build, files, prerendered, version } from '$service-worker';

const CACHE_NAME = `gps-pwa-cache-${version}`;
const ASSETS = [...build, ...files, ...prerendered];

self.addEventListener('install', (event) => {
	event.waitUntil(
		caches.open(CACHE_NAME).then(async (cache) => {
			const results = await Promise.allSettled(ASSETS.map((asset) => cache.add(asset)));
			results.forEach((result, index) => {
				if (result.status === 'rejected') {
					console.warn(`[SW] No se pudo cachear ${ASSETS[index]}:`, result.reason);
				}
			});
			await self.skipWaiting();
		})
	);
});

self.addEventListener('activate', (event) => {
	event.waitUntil(
		caches
			.keys()
			.then((keys) =>
				Promise.all(keys.map((key) => (key !== CACHE_NAME ? caches.delete(key) : undefined)))
			)
			.then(() => self.clients.claim())
	);
});

function isPrecachedAsset(request: Request): boolean {
	const pathname = new URL(request.url).pathname;
	return ASSETS.some((asset) => pathname.endsWith(asset) || pathname === asset);
}

function shouldHandle(request: Request): boolean {
	if (request.method !== 'GET') return false;
	if (request.url.includes('/api/')) return false;
	if (!request.url.startsWith(self.location.origin)) return false;
	return true;
}

async function cacheFirst(request: Request): Promise<Response> {
	const cached = await caches.match(request);
	if (cached) return cached;

	const networkResponse = await fetch(request);
	if (networkResponse.ok) {
		const cache = await caches.open(CACHE_NAME);
		await cache.put(request, networkResponse.clone());
	}
	return networkResponse;
}

async function networkFirstWithCacheFallback(request: Request): Promise<Response> {
	try {
		const networkResponse = await fetch(request);
		if (networkResponse.ok) {
			const cache = await caches.open(CACHE_NAME);
			await cache.put(request, networkResponse.clone());
		}
		return networkResponse;
	} catch (error) {
		const cached = await caches.match(request);
		if (cached) return cached;
		throw error;
	}
}

self.addEventListener('fetch', (event) => {
	const { request } = event;
	const url = new URL(request.url);

	// Ignorar completamente las peticiones a la API para no interferir con la autenticación
	if (url.pathname.startsWith('/api/')) {
		return;
	}

	if (!shouldHandle(request)) return;

	if (isPrecachedAsset(request)) {
		event.respondWith(cacheFirst(request));
	} else {
		event.respondWith(networkFirstWithCacheFallback(request));
	}
});
