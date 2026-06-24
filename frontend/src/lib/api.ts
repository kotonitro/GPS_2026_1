import { auth } from './authStore.svelte';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export async function fetchAPI(endpoint: string, options: RequestInit = {}) {
	const url = `${API_URL}${endpoint}`;

	const headers = new Headers(options.headers || {});
	if (!(options.body instanceof FormData) && !headers.has('Content-Type')) {
		headers.set('Content-Type', 'application/json');
	}

	const config: RequestInit = {
		...options,
		headers,
		credentials: 'include' // CRITICAL: This sends/receives the HttpOnly auth cookie
	};

	try {
		const response = await fetch(url, config);

		if (response.status === 401) {
			if (typeof window !== 'undefined') {
				auth.logout();
			}
			const errorData = await response.json().catch(() => ({}));
			throw { status: 401, message: errorData.error || 'Acceso denegado, inicie sesión.' };
		}

		const data = await response.json().catch(() => ({}));

		if (!response.ok) {
			throw {
				status: response.status,
				message: data.error || data.detalle || 'Ha ocurrido un error inesperado.',
				errors: data.errores
			};
		}

		return data;
	} catch (error: any) {
		if (error.status) throw error;
		throw { status: 500, message: 'No se pudo conectar al servidor backend o base de datos.' };
	}
}

export const apiClientes = {
	getAll: () => fetchAPI('/clientes/'),
	getById: (id: string) => fetchAPI(`/clientes/${id}`),
	getByRut: (rut: string) => fetchAPI(`/clientes/rut/${rut}`),
	searchByNombre: (nombre: string) => fetchAPI(`/clientes/search/nombre?nombre=${encodeURIComponent(nombre)}`),
	create: (data: { nombre: string; rut: string; telefono: string }) => fetchAPI('/clientes', {
		method: 'POST',
		body: JSON.stringify(data)
	}),
	update: (id: string, data: { nombre?: string; rut?: string; telefono?: string }) => fetchAPI(`/clientes/${id}`, {
		method: 'PATCH',
		body: JSON.stringify(data)
	}),
	delete: (id: string) => fetchAPI(`/clientes/${id}`, {
		method: 'DELETE'
	})
};

export const apiPromociones = {
	getAll: () => fetchAPI('/promociones/'),
	
	getById: (id: string) => fetchAPI(`/promociones/${id}`),
	
	// Filtro por tipo que usamos en la tabla (NXM o DESCUENTO)
	getByTipo: (tipo: string) => fetchAPI(`/promociones/tipo/${tipo}`),
	
	create: (data: { tipo: string; lleva: number; paga: number; descuento: number }) => fetchAPI('/promociones', {
		method: 'POST',
		body: JSON.stringify(data)
	}),
	
	// Usamos Partial en TS o hacemos los campos opcionales (?) para permitir actualizaciones parciales
	update: (id: string, data: { tipo?: string; lleva?: number; paga?: number; descuento?: number }) => fetchAPI(`/promociones/${id}`, {
		method: 'PATCH',
		body: JSON.stringify(data)
	}),
	
	delete: (id: string) => fetchAPI(`/promociones/${id}`, {
		method: 'DELETE'
	})
};

export const apiAuth = {
	login: (usuario: string, contrasena: string) => fetchAPI('/auth/login', {
		method: 'POST',
		body: JSON.stringify({ usuario, contrasena })
	}),
	logout: () => fetchAPI('/auth/logout', {
		method: 'POST'
	})
};
