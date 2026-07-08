<script lang="ts">
	import { onMount } from 'svelte';
	import { Search, Plus, Edit2, Trash2, Shield, User, CheckCircle2, XCircle, X } from '@lucide/svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { apiRoles } from '$lib/api';

	// Estados
	let roles = $state<any[]>([]);
	let isLoading = $state(true);
	let searchQuery = $state('');
	let filtroPermisos = $state('Todos'); // 'Todos', 'Admin', 'Regular'

	// Estados de Modales
	let showModal = $state(false);
	let showDeleteModal = $state(false);
	let isEditing = $state(false);
	let rolToDelete = $state<any>(null);

	// Formulario
	let formData = $state({
		id_rol: '',
		nombre: '',
		descripcion: '',
		es_admin: false
	});

	// Errores y Loading del form
	let errNombre = $state('');
	let errDescripcion = $state('');
	let formGeneralError = $state('');
	let submitLoading = $state(false);

	// Filtro reactivo
	let rolesFiltrados = $derived(
		roles.filter((r) => {
			// 1. Filtro de Búsqueda de texto (Nombre o Descripción)
			const coincideTexto =
				(r.nombre || r.Nombre || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
				(r.descripcion || r.Descripcion || '').toLowerCase().includes(searchQuery.toLowerCase());

			// 2. Filtro de Permisos
			const esAdmin = r.es_admin !== undefined ? r.es_admin : r.EsAdmin;
			const coincidePermisos =
				filtroPermisos === 'Todos' ||
				(filtroPermisos === 'Admin' && esAdmin) ||
				(filtroPermisos === 'Regular' && !esAdmin);

			return coincideTexto && coincidePermisos;
		})
	);

	onMount(async () => {
		await cargarDatos();
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

	// --- Controladores de Modales de Formulario ---

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

			// Si tu backend espera los nombres capitalizados (ej: Nombre, Descripcion, EsAdmin), ajustalos aquí si es necesario.
			// Asumiremos que el backend de Go con GORM puede bindear el JSON en minúsculas sin problemas.

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

	// --- Controladores del Modal de Eliminación ---

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
	<title>Roles y Permisos - MinimarketGo</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<!-- Barra Superior de Herramientas -->
	<div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
		<!-- Bloque de Búsqueda y Filtros -->
		<div class="flex flex-col gap-3 sm:flex-row sm:items-center flex-1">
			<!-- Buscador -->
			<div class="relative w-full sm:max-w-xs">
				<Search class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-text-muted" />
				<input
					type="text"
					placeholder="Buscar rol..."
					bind:value={searchQuery}
					class="w-full rounded-xl border border-border-color bg-bg-card py-2.5 pl-10 pr-4 text-sm text-text-primary transition-colors focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario"
				/>
			</div>

			<!-- Filtro por Nivel de Permisos -->
			<select
				bind:value={filtroPermisos}
				class="w-full sm:w-48 rounded-xl border border-border-color bg-bg-card py-2.5 pl-4 pr-10 text-sm text-text-primary transition-colors focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario cursor-pointer"
			>
				<option value="Todos">Todos los permisos</option>
				<option value="Admin">Administradores</option>
				<option value="Regular">Regulares</option>
			</select>

			<button
				type="button"
				title="Limpiar filtros"
				class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-muted transition-colors hover:bg-border-color hover:text-primario disabled:cursor-not-allowed disabled:opacity-50"
				onclick={() => { searchQuery = ''; filtroPermisos = 'Todos'; }}
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

	<!-- Tabla -->
	<div class="overflow-x-auto rounded-xl border border-border-color bg-bg-card shadow-sm">
		<table class="w-full text-left text-sm text-text-primary border-collapse">
			<thead class="border-b border-border-color bg-text-primary/4 text-text-muted">
				<tr>
					<th class="px-6 py-4 font-bold uppercase tracking-wider text-xs">Rol</th>
					<th class="px-6 py-4 font-bold uppercase tracking-wider text-xs">Descripción</th>
					<th class="px-6 py-4 font-bold uppercase tracking-wider text-xs">Nivel de Acceso</th>
					<th class="px-6 py-4 text-right font-bold uppercase tracking-wider text-xs w-[120px]"
						>Acciones</th
					>
				</tr>
			</thead>
			<tbody class="divide-y divide-border-color">
				{#if isLoading}
					<tr>
						<td colspan="4" class="py-8 text-center text-text-muted">Cargando roles...</td>
					</tr>
				{:else if rolesFiltrados.length === 0}
					<tr>
						<td colspan="4" class="py-16 text-center">
							<div class="flex flex-col items-center justify-center">
								<Shield size={48} class="text-text-muted opacity-30 mb-4" />
								<h3 class="mb-1 text-lg font-semibold text-text-primary">
									No se encontraron roles
								</h3>
								<p class="text-sm text-text-secondary">Prueba con otro filtro o crea uno nuevo.</p>
							</div>
						</td>
					</tr>
				{:else}
					{#each rolesFiltrados as rol (rol.id_rol || rol.id || rol.ID)}
						{@const nombre = rol.nombre || rol.Nombre}
						{@const desc = rol.descripcion || rol.Descripcion}
						{@const esAdmin = rol.es_admin !== undefined ? rol.es_admin : rol.EsAdmin}
						<tr class="transition-colors hover:bg-text-primary/[0.015]">
							<!-- Nombre -->
							<td class="px-6 py-4">
								<span class="font-semibold">{nombre}</span>
							</td>

							<!-- Descripción -->
							<td class="px-6 py-4 text-text-muted max-w-md truncate" title={desc}>
								{desc || '-'}
							</td>

							<!-- Nivel de Acceso -->
							<td class="px-6 py-4">
								<div class="flex items-center gap-1.5">
									{#if esAdmin}
										<span
											class="inline-flex items-center gap-1.5 rounded-full bg-primario/10 px-2.5 py-1 text-xs font-medium text-primario"
										>
											<Shield size={14} /> Acceso Total (Admin)
										</span>
									{:else}
										<span
											class="inline-flex items-center gap-1.5 rounded-full bg-border-color px-2.5 py-1 text-xs font-medium text-text-secondary"
										>
											<User size={14} /> Acceso Limitado
										</span>
									{/if}
								</div>
							</td>

							<!-- Acciones -->
							<td class="px-6 py-4 text-right">
								<!-- Protección para no editar ni eliminar el super admin base (opcional, pero buena práctica) -->
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

<!-- Modal Formulario (Crear/Editar) -->
{#if showModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-4 backdrop-blur-sm animate-modal-enter"
		onclick={() => (showModal = false)}
		role="presentation"
	>
		<div
			class="w-full max-w-[500px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
		>
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-text-primary">
					{isEditing ? 'Editar Rol' : 'Añadir Nuevo Rol'}
				</h2>
				<button
					onclick={() => (showModal = false)}
					class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary"
				>
					&times;
				</button>
			</header>

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
					<!-- Nombre del Rol -->
					<div class="flex flex-col gap-1.5">
						<label for="nombre" class="text-sm font-semibold text-text-secondary"
							>Nombre del Rol</label
						>
						<input
							id="nombre"
							type="text"
							autocomplete="off"
							bind:value={formData.nombre}
							disabled={submitLoading}
							required
							placeholder="Ej: Supervisor"
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-3 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
						/>
						{#if errNombre}<span class="text-xs font-medium text-danger-color">{errNombre}</span
							>{/if}
					</div>

					<!-- Descripción -->
					<div class="flex flex-col gap-1.5">
						<label for="descripcion" class="text-sm font-semibold text-text-secondary"
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
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-3 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50 resize-none"
						></textarea>
						{#if errDescripcion}<span class="text-xs font-medium text-danger-color"
								>{errDescripcion}</span
							>{/if}
					</div>

					<!-- Permisos (Switch) -->
					<div
						class="flex items-start gap-4 rounded-xl border border-border-color bg-bg-primary/50 p-4 mt-2"
					>
						<input
							id="es_admin"
							type="checkbox"
							bind:checked={formData.es_admin}
							disabled={submitLoading}
							class="mt-1 h-5 w-5 accent-primario cursor-pointer disabled:opacity-50 shrink-0"
						/>
						<div class="flex flex-col">
							<label
								for="es_admin"
								class="text-sm font-bold text-text-primary cursor-pointer {submitLoading
									? 'opacity-50'
									: ''}"
							>
								Otorgar privilegios de Administrador
							</label>
							<p class="text-xs text-text-muted mt-1 leading-relaxed">
								Si está activado, los empleados con este rol tendrán acceso total al sistema
								(incluyendo configuraciones, dashboard completo y gestión de otros usuarios).
							</p>
						</div>
					</div>
				</div>

				<div class="mt-8 flex justify-end gap-3">
					<button
						type="button"
						onclick={() => (showModal = false)}
						disabled={submitLoading}
						class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5 disabled:opacity-50"
					>
						Cancelar
					</button>
					<button
						type="submit"
						disabled={submitLoading}
						class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-gradient-to-r from-accent-light to-accent px-5 py-2.5 text-sm font-semibold text-white transition-all duration-200 hover:-translate-y-[1px] hover:shadow-glow disabled:cursor-not-allowed disabled:opacity-50"
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

<!-- Modal Confirmación de Eliminación -->
{#if showDeleteModal && rolToDelete}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-4 backdrop-blur-[4px]"
		onclick={() => (showDeleteModal = false)}
		role="presentation"
	>
		<div
			class="w-full max-w-[500px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
		>
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-danger-color">Confirmar Eliminación</h2>
				<button
					class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary"
					onclick={() => (showDeleteModal = false)}>&times;</button
				>
			</header>
			<div class="p-6 text-text-primary">
				<p>
					¿Estás seguro de que deseas eliminar el rol <strong
						>{rolToDelete.nombre || rolToDelete.Nombre}</strong
					>?
				</p>
				<p class="mt-3 text-sm text-text-muted">
					Si hay empleados asignados a este rol, no podrás eliminarlo hasta que los reasignes a un
					rol diferente.
				</p>
			</div>
			<footer
				class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-4 px-6"
			>
				<button
					type="button"
					class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5"
					onclick={() => (showDeleteModal = false)}>Cancelar</button
				>
				<button
					type="button"
					class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-red-500/15 bg-danger-bg px-5 py-2.5 text-sm font-semibold text-danger-color transition-all duration-200 hover:bg-danger-color hover:text-white"
					onclick={confirmarEliminacion}>Eliminar Rol</button
				>
			</footer>
		</div>
	</div>
{/if}
