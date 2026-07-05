const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export interface Empleado {
    id: string;
	nombre: string,
    usuario: string;
    rol: string;
}

export async function apiFetch(endpoint: string, options: RequestInit = {}) {
	const headers = {
		'Content-Type': 'application/json',
		...options.headers
	};

	const response = await fetch(`${API_URL}${endpoint}`, {
		...options,
		headers,
		credentials: 'include'
	});

	if (!response.ok) {
		const errorData = await response.json().catch(() => ({}));
		throw new Error(errorData.error || 'Error de conexión con el servidor');
	}

	return response.json();
}

export async function login(usuario: string, contrasena: string): Promise<Empleado> {
    return apiFetch('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ usuario, contrasena })
    });
}

export async function logout(): Promise<void> {
    return apiFetch('/auth/logout', {
        method: 'POST'
    });
}

export async function checkSession(): Promise<Empleado> {
    return apiFetch('/auth/me');
}