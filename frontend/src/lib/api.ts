const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

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
