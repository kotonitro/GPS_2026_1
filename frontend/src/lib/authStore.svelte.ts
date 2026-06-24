import { goto } from '$app/navigation';

export interface User {
	id_empleado: string;
	usuario: string;
	rol: string;
}

class AuthStore {
	user = $state<User | null>(null);
	loading = $state(true);

	init() {
		if (typeof window !== 'undefined') {
			const stored = localStorage.getItem('user_session');
			if (stored) {
				try {
					this.user = JSON.parse(stored);
				} catch {
					this.user = null;
				}
			}
			this.loading = false;
		}
	}

	login(userData: User) {
		this.user = userData;
		if (typeof window !== 'undefined') {
			localStorage.setItem('user_session', JSON.stringify(userData));
		}
	}

	logout() {
		this.user = null;
		if (typeof window !== 'undefined') {
			localStorage.removeItem('user_session');
			goto('/login');
		}
	}
}

export const auth = new AuthStore();
