export interface Toast {
	id: string;
	message: string;
	type: 'success' | 'error';
}

class ToastStore {
	toasts = $state<Toast[]>([]);

	show(message: string, type: 'success' | 'error' = 'success', duration = 4000) {
		const id = Math.random().toString(36).substring(2, 9);
		this.toasts.push({ id, message, type });
		setTimeout(() => {
			this.toasts = this.toasts.filter((t) => t.id !== id);
		}, duration);
	}

	dismiss(id: string) {
		this.toasts = this.toasts.filter((t) => t.id !== id);
	}
}

export const toast = new ToastStore();
