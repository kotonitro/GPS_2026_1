<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Search,
		Plus,
		Edit2,
		Trash2,
		Shield,
		User,
		CheckCircle2,
		XCircle,
		X
	} from '@lucide/svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { apiRoles, checkSession } from '$lib/api';
	import { goto } from '$app/navigation';

	let roles = $state<any[]>([]);
	let isLoading = $state(true);
	let searchQuery = $state('');
	let filtroPermisos = $state('Todos');

	let showModal = $state(false);
	let showDeleteModal = $state(false);
	let isEditing = $state(false);
	let rolToDelete = $state<any>(null);
	let ultimoModificadoNombre = $state<string | null>(null);

	let formData = $state({
		id_rol: '',
		nombre: '',
		descripcion: '',
		es_admin: false
	});

	let errNombre = $state('');
	let errDescripcion = $state('');
	let formGeneralError = $state('');
	let submitLoading = $state(false);

	let rolesFiltrados = $derived(
		roles
			.filter((r) => {
				const coincideTexto =
					(r.nombre || r.Nombre || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
					(r.descripcion || r.Descripcion || '').toLowerCase().includes(searchQuery.toLowerCase());

				const esAdmin = r.es_admin !== undefined ? r.es_admin : r.EsAdmin;
				const coincidePermisos =
					filtroPermisos === 'Todos' ||
					(filtroPermisos === 'Admin' && esAdmin) ||
					(filtroPermisos === 'Regular' && !esAdmin);

				return coincideTexto && coincidePermisos;
			})
			.sort((a, b) => {
				const nombreA = a.nombre || a.Nombre || '';
				const nombreB = b.nombre || b.Nombre || '';

				if (nombreA === 'Admin') return -1;
				if (nombreB === 'Admin') return 1;

				if (nombreA === ultimoModificadoNombre) return -1;
				if (nombreB === ultimoModificadoNombre) return 1;

				return nombreA.localeCompare(nombreB);
			})
	);

	onMount(async () => {
		try {
			const empleado = await checkSession();

			const objRol = empleado?.es_admin;
			const esAdmin = objRol === true;

			if (!esAdmin) {
				toast.show('Acceso denegado. Se requieren privilegios de administrador.', 'error');
				goto('/dashboard');
				return;
			}
			await cargarDatos();
		} catch (error) {
			goto('/login');
		}
	});

	async function cargarDatos() {
		isLoading = true;
		try {
			const res = await apiRoles.getAll();
			roles = Array.isArray(res) ? res : [];
		} catch (error) {
			toast.show('Error al cargar la lista de roles.', 'error');
		} finally {
			isLoading = false;
		}
	}

	function clearErrors() {
		errNombre = '';
		errDescripcion = '';
		formGeneralError = '';
	}

	function abrirModalNuevo() {
		isEditing = false;
		clearErrors();
		formData = {
			id_rol: '',
			nombre: '',
			descripcion: '',
			es_admin: false
		};
		showModal = true;
	}

	function abrirModalEditar(rol: any) {
		isEditing = true;
		clearErrors();
		formData = {
			id_rol: rol.id_rol || rol.id || rol.ID,
			nombre: rol.nombre || rol.Nombre,
			descripcion: rol.descripcion || rol.Descripcion,
			es_admin: rol.es_admin !== undefined ? rol.es_admin : rol.EsAdmin
		};
		showModal = true;
	}

	async function guardarRol(e: Event) {
		e.preventDefault();
		clearErrors();
		let isValid = true;

		if (!formData.nombre.trim()) {
			errNombre = 'El nombre del rol es obligatorio.';
			isValid = false;
		}

		if (!formData.descripcion.trim()) {
			errDescripcion = 'La descripción es obligatoria para identificar el rol.';
			isValid = false;
		}

		if (!isValid) return;

		submitLoading = true;
		try {
			const payload = { ...formData };

			ultimoModificadoNombre = payload.nombre;

			if (isEditing) {
				await apiRoles.update(payload.id_rol, payload);
				toast.show('Rol actualizado correctamente.', 'success');
			} else {
				await apiRoles.create(payload);
				toast.show('Rol creado exitosamente.', 'success');
			}

			showModal = false;
			await cargarDatos();
		} catch (error: any) {
			formGeneralError = error.message || 'Ocurrió un error al guardar el rol.';
			toast.show('No se pudo guardar el rol.', 'error');
		} finally {
			submitLoading = false;
		}
	}

	function abrirModalEliminar(rol: any) {
		rolToDelete = rol;
		showDeleteModal = true;
	}

	async function confirmarEliminacion() {
		if (!rolToDelete) return;
		try {
			const id = rolToDelete.id_rol || rolToDelete.id || rolToDelete.ID;
			await apiRoles.delete(id);
			toast.show('Rol eliminado con éxito.', 'success');
			showDeleteModal = false;
			rolToDelete = null;
			await cargarDatos();
		} catch (error: any) {
			toast.show(
				error.message || 'No se pudo eliminar el rol. Comprueba si hay empleados asignados a él.',
				'error'
			);
		}
	}
</script>

<svelte:head>
	<title>Gestión de Roles - MinimarketGo</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
		<div class="flex flex-col gap-3 sm:flex-row sm:items-center flex-1">
			<div class="relative w-full sm:max-w-xs">
				<Search class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-text-muted" />
				<input
					type="text"
					placeholder="Buscar rol..."
					bind:value={searchQuery}
					class="w-full rounded-xl border border-border-color bg-bg-card py-2.5 pl-10 pr-4 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario"
				/>
			</div>
			<select
				bind:value={filtroPermisos}
				class="w-full sm:w-48 rounded-xl border border-border-color bg-bg-card py-2.5 pl-4 pr-10 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario cursor-pointer"
			>
				<option value="Todos">Todos los permisos</option>
				<option value="Admin">Administrativo</option>
				<option value="Regular">Regular</option>
			</select>

			<button
				type="button"
				title="Limpiar filtros"
				class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-muted hover:bg-border-color hover:text-primario disabled:cursor-not-allowed disabled:opacity-50"
				onclick={() => {
					searchQuery = '';
					filtroPermisos = 'Todos';
				}}
				disabled={!searchQuery && filtroPermisos === 'Todos'}
			>
				<X size={14} strokeWidth={2.5} />
			</button>
		</div>

		<button
			onclick={abrirModalNuevo}
			class="flex w-full items-center justify-center gap-2 rounded-xl bg-primario px-5 py-2.5 text-sm font-semibold text-white transition-transform hover:scale-[1.02] active:scale-95 sm:w-auto"
		>
			<Plus size={18} />
			Nuevo Rol
		</button>
	</div>

	<div class="overflow-x-auto rounded-xl border border-border-color bg-bg-card shadow-sm">
		<table class="w-full whitespace-nowrap text-left text-sm text-text-primary">
			<thead class="border-b border-border-color bg-bg-primary/50 text-text-muted">
				<tr>
					<th class="px-6 py-4 font-semibold">Rol</th>
					<th class="px-6 py-4 font-semibold">Descripción</th>
					<th class="px-6 py-4 font-semibold">Permisos</th>
					<th class="px-6 py-4 text-right font-semibold">Acciones</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-border-color">
				{#if isLoading}
					<tr>
						<td colspan="4" class="py-8 text-center text-text-muted">Cargando roles...</td>
					</tr>
				{:else if rolesFiltrados.length === 0}
					<tr>
						<td colspan="4" class="py-8 text-center text-text-muted">
							No se encontraron roles con los filtros aplicados.
						</td>
					</tr>
				{:else}
					{#each rolesFiltrados as rol (rol.id_rol || rol.id || rol.ID)}
						{@const nombre = rol.nombre || rol.Nombre}
						{@const desc = rol.descripcion || rol.Descripcion}
						{@const esAdmin = rol.es_admin !== undefined ? rol.es_admin : rol.EsAdmin}
						<tr class="hover:bg-bg-primary/30">
							<td class="px-6 py-4">
								<div class="flex items-center gap-2">
									<span class="font-semibold">{nombre}</span>
									{#if nombre === 'Admin'}
										<span
											class="rounded-md bg-primario/10 px-1.5 py-0.5 text-[10px] font-bold uppercase tracking-wider text-primario"
										>
											Sistema
										</span>
									{/if}

									{#if nombre === ultimoModificadoNombre && nombre !== 'Admin'}
										<span
											class="rounded-md bg-amber-500/10 px-1.5 py-0.5 text-[10px] font-bold uppercase tracking-wider text-amber-600"
										>
											Reciente
										</span>
									{/if}
								</div>
							</td>

							<td class="px-6 py-4 text-text-muted max-w-md truncate" title={desc}>
								{desc || '-'}
							</td>

							<td class="px-6 py-4">
								<div class="flex items-center gap-1.5">
									{#if esAdmin}
										<Shield size={16} class="text-primario" />
										<span class="font-medium text-text-primary">Acceso administrativo</span>
									{:else}
										<User size={16} class="text-text-muted" />
										<span class="font-medium text-text-primary">Acceso regular</span>
									{/if}
								</div>
							</td>

							<td class="px-6 py-4 text-right">
								{#if nombre === 'Admin' && esAdmin}
									<span class="text-xs text-text-muted italic px-2">Sistema</span>
								{:else}
									<div class="flex items-center justify-end gap-2">
										<button
											onclick={() => abrirModalEditar(rol)}
											class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-primario"
											title="Editar Rol"
										>
											<Edit2 size={16} />
										</button>
										<button
											onclick={() => abrirModalEliminar(rol)}
											class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-red-500/20 hover:bg-danger-bg hover:text-danger-color"
											title="Eliminar Rol"
										>
											<Trash2 size={16} />
										</button>
									</div>
								{/if}
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>

{#if showModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-4 backdrop-blur-sm animate-modal-enter"
	>
		<div class="w-full max-w-lg rounded-2xl border border-border-color bg-bg-card shadow-2xl">
			<div class="flex items-center justify-between border-b border-border-color px-6 py-4">
				<h3 class="text-lg font-bold text-text-primary">
					{isEditing ? 'Editar Rol' : 'Añadir Nuevo Rol'}
				</h3>
				<button
					onclick={() => (showModal = false)}
					class="rounded-lg p-1 text-text-muted hover:bg-border-color hover:text-text-primary"
				>
					<X size={20} />
				</button>
			</div>

			<form onsubmit={guardarRol} class="p-6">
				{#if formGeneralError}
					<div
						class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color"
						role="alert"
					>
						<span>{formGeneralError}</span>
					</div>
				{/if}

				<div class="flex flex-col gap-5">
					<div class="flex flex-col gap-1.5">
						<label for="nombre" class="text-sm font-semibold text-text-primary"
							>Nombre del Rol</label
						>
						<input
							id="nombre"
							type="text"
							autocomplete="off"
							bind:value={formData.nombre}
							disabled={submitLoading}
							required
							placeholder="Ej: Cajero"
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
						/>
						{#if errNombre}<span class="text-xs font-medium text-danger-color">{errNombre}</span
							>{/if}
					</div>

					<div class="flex flex-col gap-1.5">
						<label for="descripcion" class="text-sm font-semibold text-text-primary"
							>Descripción</label
						>
						<textarea
							id="descripcion"
							autocomplete="off"
							bind:value={formData.descripcion}
							disabled={submitLoading}
							required
							rows="3"
							placeholder="Describe los permisos y funciones de este rol..."
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50 resize-none"
						></textarea>
						{#if errDescripcion}<span class="text-xs font-medium text-danger-color"
								>{errDescripcion}</span
							>{/if}
					</div>

					<div class="flex items-center gap-3 pt-2">
						<input
							id="es_admin"
							type="checkbox"
							bind:checked={formData.es_admin}
							disabled={submitLoading}
							class="h-4 w-4 accent-primario cursor-pointer disabled:opacity-50"
						/>
						<div class="flex flex-col">
							<label
								for="es_admin"
								class="text-sm font-semibold text-text-primary cursor-pointer {submitLoading
									? 'opacity-50'
									: ''}"
							>
								Otorgar permisos administrativos
							</label>
							<p class="text-xs text-text-muted mt-0.5">
								Permite que los empleados que posean este rol tengan acceso administrativo en el
								sistema.
							</p>
						</div>
					</div>
				</div>

				<div class="mt-8 flex justify-end gap-3">
					<button
						type="button"
						onclick={() => (showModal = false)}
						disabled={submitLoading}
						class="rounded-xl px-5 py-2.5 text-sm font-semibold text-text-muted hover:bg-border-color hover:text-text-primary disabled:opacity-50"
					>
						Cancelar
					</button>
					<button
						type="submit"
						disabled={submitLoading}
						class="rounded-xl bg-primario px-5 py-2.5 text-sm font-semibold text-white transition-transform hover:scale-[1.02] active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
					>
						{#if submitLoading}
							Procesando...
						{:else}
							{isEditing ? 'Guardar Cambios' : 'Crear Rol'}
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

{#if showDeleteModal && rolToDelete}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm animate-modal-enter"
		onclick={() => (showDeleteModal = false)}
		role="presentation"
	>
		<div
			class="w-full max-w-lg rounded-2xl border border-border-color bg-bg-card shadow-2xl"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
		>
			<div class="flex items-center justify-between border-b border-border-color px-6 py-4">
				<h3 class="text-lg font-bold text-danger-color">Confirmar Eliminación</h3>
				<button
					class="rounded-lg p-1 text-text-muted hover:bg-border-color hover:text-text-primary"
					onclick={() => (showDeleteModal = false)}
				>
					<X size={20} />
				</button>
			</div>

			<div class="p-6 text-text-primary">
				<p>
					¿Estás seguro de que deseas eliminar el rol <strong
						>{rolToDelete.nombre || rolToDelete.Nombre}</strong
					>?
				</p>
				<p class="mt-3 text-xs text-text-muted">
					Si hay empleados asignados a este rol, no podrás eliminarlo hasta que los reasignes a un
					rol diferente.
				</p>
			</div>

			<div class="flex justify-end gap-3 border-t border-border-color bg-bg-primary/30 px-6 py-4">
				<button
					type="button"
					class="rounded-xl px-5 py-2.5 text-sm font-semibold text-text-muted hover:bg-border-color hover:text-text-primary"
					onclick={() => (showDeleteModal = false)}
				>
					Cancelar
				</button>
				<button
					type="button"
					class="rounded-xl bg-danger-color px-5 py-2.5 text-sm font-semibold text-white transition-transform hover:scale-[1.02] active:scale-95"
					onclick={confirmarEliminacion}
				>
					Eliminar Rol
				</button>
			</div>
		</div>
	</div>
{/if}
