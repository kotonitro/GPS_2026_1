<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Search,
		Plus,
		Edit2,
		Trash2,
		Shield,
		User,
		X,
		CheckCircle2,
		XCircle,
		Eye,
		EyeOff
	} from '@lucide/svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { apiEmpleados, apiRoles } from '$lib/api';

	// Estados
	let empleados = $state<any[]>([]);
	let roles = $state<any[]>([]);
	let searchQuery = $state('');
	let isModalOpen = $state(false);
	let isEditing = $state(false);
	let isLoading = $state(true);

	// Nuevos estados para los filtros
	let filtroEstado = $state('Todos'); // 'Todos', 'Activos', 'Inactivos'
	let filtroRol = $state('Todos'); // 'Todos', 'Admin', 'Cajero', etc.

	// Formulario
	let formData = $state({
		id_empleado: '',
		rut: '',
		nombre: '',
		usuario: '',
		contrasena: '',
		telefono: '',
		id_rol: '',
		activo: true
	});

	let errRut = $state('');
	let errNombre = $state('');
	let errUsuario = $state('');
	let errContrasena = $state('');
	let errTelefono = $state('');
	let formGeneralError = $state('');
	let submitLoading = $state(false);
	let showPassword = $state(false);
	let showDeleteModal = $state(false);
	let empleadoToDelete = $state<any>(null);

	function clearErrors() {
		errRut = '';
		errNombre = '';
		errUsuario = '';
		errContrasena = '';
		errTelefono = '';
		formGeneralError = '';
	}

	function validarContrasenaFormato(pass: string): boolean {
		const min8 = pass.length >= 8;
		const numero = /[0-9]/.test(pass);
		const especial = /[!@#~$%^&*(),.?":{}|<>]/.test(pass);
		return min8 && numero && especial;
	}

	function validarRut(rutOriginal: string): boolean {
		const rutLimpio = rutOriginal.replace(/\./g, '').replace(/-/g, '').trim().toUpperCase();
		if (rutLimpio.length < 8) return false;
		const cuerpo = rutLimpio.slice(0, -1);
		const dvIngresado = rutLimpio.slice(-1);

		let suma = 0;
		let multiplicador = 2;
		for (let i = cuerpo.length - 1; i >= 0; i--) {
			const digito = parseInt(cuerpo[i], 10);
			if (isNaN(digito)) return false;
			suma += digito * multiplicador;
			multiplicador++;
			if (multiplicador > 7) multiplicador = 2;
		}

		const resto = suma % 11;
		const resultado = 11 - resto;
		let dvCalculado = '';
		if (resultado === 11) dvCalculado = '0';
		else if (resultado === 10) dvCalculado = 'K';
		else dvCalculado = resultado.toString();

		return dvCalculado === dvIngresado;
	}

	function validarTelefono(tel: string): boolean {
		const patron = /^(\+56)?[\s.-]?9[\s.-]?\d{4}[\s.-]?\d{4}$/;
		return patron.test(tel);
	}

	function formatRutInput(val: string) {
		let clean = val.replace(/[^0-9kK]/g, '');
		if (clean.length === 0) return '';
		const dv = clean.slice(-1);
		let cuerpo = clean.slice(0, -1);
		if (cuerpo.length > 0) {
			let formattedCuerpo = '';
			let j = 0;
			for (let i = cuerpo.length - 1; i >= 0; i--) {
				formattedCuerpo = cuerpo[i] + formattedCuerpo;
				j++;
				if (j === 3 && i > 0) {
					formattedCuerpo = '.' + formattedCuerpo;
					j = 0;
				}
			}
			return `${formattedCuerpo}-${dv}`;
		}
		return dv;
	}

	function handleRutInput(e: Event) {
		const target = e.target as HTMLInputElement;
		const formatted = formatRutInput(target.value);
		formData.rut = formatted;
		target.value = formatted;
	}

	function formatNombreInput(val: string) {
		return val
			.toLowerCase()
			.replace(/\s+/g, ' ')
			.split(' ')
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	}

	function handleNombreBlur() {
		formData.nombre = formatNombreInput(formData.nombre).trim();
	}

	// Filtro reactivo actualizado (Texto + Estado + Rol)
	let empleadosFiltrados = $derived(
		empleados.filter((emp) => {
			// 1. Filtro de Búsqueda de texto (Nombre o RUT)
			const coincideTexto =
				emp.nombre.toLowerCase().includes(searchQuery.toLowerCase()) ||
				emp.rut.toLowerCase().includes(searchQuery.toLowerCase());

			// 2. Filtro de Estado
			const esActivo = emp.activo !== undefined ? emp.activo : emp.Activo;
			const coincideEstado =
				filtroEstado === 'Todos' ||
				(filtroEstado === 'Activos' && esActivo) ||
				(filtroEstado === 'Inactivos' && !esActivo);

			// 3. Filtro de Rol
			const objRol = emp.Rol || emp.rol || {};
			const nombreRol = objRol.nombre || objRol.Nombre || '';
			const coincideRol = filtroRol === 'Todos' || nombreRol === filtroRol;

			return coincideTexto && coincideEstado && coincideRol;
		})
	);

	onMount(async () => {
		await cargarDatos();
	});

	async function cargarDatos() {
		isLoading = true;
		try {
			const [resEmpleados, resRoles] = await Promise.all([
				apiEmpleados.getAll(),
				apiRoles.getAll()
			]);

			empleados = Array.isArray(resEmpleados) ? resEmpleados : [];
			roles = Array.isArray(resRoles) ? resRoles : [];
		} catch (error) {
			toast.show('Error al cargar los datos del sistema.', 'error'); // <- Corregido
		} finally {
			isLoading = false;
		}
	}

	function abrirModalNuevo() {
		isEditing = false;
		showPassword = false;
		formData = {
			id_empleado: '',
			rut: '',
			nombre: '',
			usuario: '',
			contrasena: '',
			telefono: '',
			id_rol: roles.length > 0 ? roles[0].id_rol || roles[0].id || roles[0].ID : '',
			activo: true
		};
		isModalOpen = true;
	}

	function abrirModalEditar(empleado: any) {
		isEditing = true;
		showPassword = false;
		formData = {
			...empleado,
			id_rol:
				empleado.id_rol ||
				empleado.RolID ||
				(empleado.Rol ? empleado.Rol.id || empleado.Rol.ID : ''),
			contrasena: ''
		};
		isModalOpen = true;
	}

	function cerrarModal() {
		isModalOpen = false;
	}

	async function guardarEmpleado(e: Event) {
		e.preventDefault();
		clearErrors();
		let isValid = true;

		// 1. Validar Nombre
		if (!formData.nombre.trim()) {
			errNombre = 'El nombre es obligatorio.';
			isValid = false;
		} else {
			const soloLetras = /^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$/.test(formData.nombre);
			const partes = formData.nombre.trim().split(/\s+/);
			if (!soloLetras) {
				errNombre = 'Solo debe contener letras y espacios.';
				isValid = false;
			} else if (partes.length < 2) {
				errNombre = 'Ingrese nombre y apellido.';
				isValid = false;
			} else {
				formData.nombre = formatNombreInput(formData.nombre).trim();
			}
		}

		// 2. Validar RUT
		const rawRut = formData.rut.replace(/\./g, '');
		if (!formData.rut) {
			errRut = 'El RUT es obligatorio.';
			isValid = false;
		} else if (!validarRut(rawRut)) {
			errRut = 'El RUT ingresado no es válido.';
			isValid = false;
		}

		// 3. Validar Usuario
		if (!formData.usuario.trim()) {
			errUsuario = 'El usuario es obligatorio.';
			isValid = false;
		} else if (formData.usuario.length < 3) {
			errUsuario = 'Debe tener al menos 3 caracteres.';
			isValid = false;
		}

		// 4. Validar Teléfono
		if (!formData.telefono.trim()) {
			errTelefono = 'El teléfono es obligatorio.';
			isValid = false;
		} else if (!validarTelefono(formData.telefono.trim())) {
			errTelefono = 'Formato inválido (ej: +56912345678).';
			isValid = false;
		}

		// 5. Validar Contraseña
		if (!isEditing && !formData.contrasena) {
			errContrasena = 'La contraseña es obligatoria.';
			isValid = false;
		} else if (formData.contrasena && !validarContrasenaFormato(formData.contrasena)) {
			errContrasena = 'Debe tener mín. 8 caracteres, un número y un carácter especial.';
			isValid = false;
		}

		if (!isValid) return;

		submitLoading = true;
		try {
			const payload = { ...formData, rut: rawRut };

			if (isEditing) {
				const dataToUpdate = { ...payload };
				if (!dataToUpdate.contrasena) delete dataToUpdate.contrasena;

				await apiEmpleados.update(payload.id_empleado || payload.ID, dataToUpdate);
				toast.show('Empleado actualizado correctamente.', 'success'); // <- Corregido
			} else {
				await apiEmpleados.create(payload);
				toast.show('Empleado creado exitosamente.', 'success'); // <- Corregido
			}

			cerrarModal();
			await cargarDatos();
		} catch (error: any) {
			if (error.errors) {
				if (error.errors.Rut || error.errors.rut) errRut = error.errors.Rut || error.errors.rut;
				if (error.errors.Usuario || error.errors.usuario)
					errUsuario = error.errors.Usuario || error.errors.usuario;
			} else {
				formGeneralError = error.message || 'Ocurrió un error al guardar el empleado.';
			}
			toast.show('No se pudo guardar el empleado.', 'error'); // <- Corregido
		} finally {
			submitLoading = false;
		}
	}

	function abrirModalEliminar(empleado: any) {
		empleadoToDelete = empleado;
		showDeleteModal = true;
	}

	async function confirmarEliminacion() {
		if (!empleadoToDelete) return;
		try {
			const id = empleadoToDelete.id_empleado || empleadoToDelete.ID;
			await apiEmpleados.delete(id);
			toast.show('Empleado eliminado con éxito.', 'success');
			showDeleteModal = false;
			empleadoToDelete = null;
			await cargarDatos();
		} catch (error: any) {
			toast.show(error.message || 'No se pudo eliminar el empleado.', 'error');
		}
	}
</script>

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
					placeholder="Buscar por nombre o RUT..."
					bind:value={searchQuery}
					class="w-full rounded-xl border border-border-color bg-bg-card py-2.5 pl-10 pr-4 text-sm text-text-primary transition-colors focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario"
				/>
			</div>

			<select
				bind:value={filtroRol}
				class="w-full sm:w-auto rounded-xl border border-border-color bg-bg-card py-2.5 pl-4 pr-10 text-sm text-text-primary transition-colors focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario cursor-pointer"
			>
				<option value="Todos">Todos los roles</option>
				{#each roles as r}
					<option value={r.nombre || r.Nombre}>{r.nombre || r.Nombre}</option>
				{/each}
			</select>

			<!-- Filtro por Estado -->
			<select
				bind:value={filtroEstado}
				class="w-full sm:w-auto rounded-xl border border-border-color bg-bg-card py-2.5 pl-4 pr-10 text-sm text-text-primary transition-colors focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario cursor-pointer"
			>
				<option value="Todos">Todos los estados</option>
				<option value="Activos">Activo</option>
				<option value="Inactivos">Inactivo</option>
			</select>

			<button
				type="button"
				title="Limpiar filtros"
				class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-muted transition-colors hover:bg-border-color hover:text-primario disabled:cursor-not-allowed disabled:opacity-50"
				onclick={() => { searchQuery = ''; filtroRol = 'Todos'; filtroEstado = 'Todos'; }}
				disabled={!searchQuery && filtroRol === 'Todos' && filtroEstado === 'Todos'}
			>
				<X size={14} strokeWidth={2.5} />
			</button>
		</div>

		<button
			onclick={abrirModalNuevo}
			class="flex w-full items-center justify-center gap-2 rounded-xl bg-primario px-5 py-2.5 text-sm font-semibold text-white transition-transform hover:scale-[1.02] active:scale-95 sm:w-auto"
		>
			<Plus size={18} />
			Nuevo Empleado
		</button>
	</div>

	<!-- Tabla -->
	<div class="overflow-x-auto rounded-xl border border-border-color bg-bg-card shadow-sm">
		<table class="w-full whitespace-nowrap text-left text-sm text-text-primary">
			<thead class="border-b border-border-color bg-bg-primary/50 text-text-muted">
				<tr>
					<th class="px-6 py-4 font-semibold">Empleado</th>
					<th class="px-6 py-4 font-semibold">Usuario</th>
					<th class="px-6 py-4 font-semibold">Teléfono</th>
					<th class="px-6 py-4 font-semibold">Rol</th>
					<th class="px-6 py-4 font-semibold">Estado</th>
					<th class="px-6 py-4 text-right font-semibold">Acciones</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-border-color">
				{#if isLoading}
					<tr>
						<td colspan="6" class="py-8 text-center text-text-muted">Cargando datos...</td>
					</tr>
				{:else if empleadosFiltrados.length === 0}
					<tr>
						<td colspan="6" class="py-8 text-center text-text-muted"
							>No se encontraron empleados con los filtros aplicados.</td
						>
					</tr>
				{:else}
					{#each empleadosFiltrados as emp (emp.id_empleado || emp.ID)}
						{@const objRol = emp.Rol || emp.rol || {}}
						<tr class="transition-colors hover:bg-bg-primary/30">
							<!-- 1. Empleado (Nombre y RUT) -->
							<td class="px-6 py-4">
								<div class="flex items-center gap-3">
									<div
										class="flex h-9 w-9 items-center justify-center rounded-full bg-border-color text-xs font-bold text-text-primary"
									>
										{emp.nombre.substring(0, 2).toUpperCase()}
									</div>
									<div class="flex flex-col">
										<span class="font-semibold">{emp.nombre}</span>
										<span class="text-xs text-text-muted">{emp.rut}</span>
									</div>
								</div>
							</td>

							<!-- 2. Usuario -->
							<td class="px-6 py-4 font-medium">
								{emp.usuario}
							</td>

							<!-- 3. Teléfono -->
							<td class="px-6 py-4 text-text-muted">
								{emp.telefono || '-'}
							</td>

							<!-- 4. Rol -->
							<td class="px-6 py-4">
								<div class="flex items-center gap-1.5">
									{#if objRol.es_admin || objRol.EsAdmin}
										<Shield size={16} class="text-primario" />
									{:else}
										<User size={16} class="text-text-muted" />
									{/if}
									<span class="font-medium">{objRol.nombre || objRol.Nombre || 'Sin Rol'}</span>
								</div>
							</td>

							<!-- 5. Estado -->
							<td class="px-6 py-4">
								{#if emp.activo || emp.Activo}
									<span
										class="inline-flex items-center gap-1.5 rounded-full bg-green-500/10 px-2.5 py-1 text-xs font-medium text-green-500"
									>
										<CheckCircle2 size={14} /> Activo
									</span>
								{:else}
									<span
										class="inline-flex items-center gap-1.5 rounded-full bg-red-500/10 px-2.5 py-1 text-xs font-medium text-red-500"
									>
										<XCircle size={14} /> Inactivo
									</span>
								{/if}
							</td>

							<!-- 6. Acciones -->
							<td class="px-6 py-4 text-right">
								<div class="flex items-center justify-end gap-2">
									<button
										onclick={() => abrirModalEditar(emp)}
										class="rounded-lg p-2 text-text-muted transition-colors hover:bg-border-color hover:text-primario"
										title="Editar"
									>
										<Edit2 size={18} />
									</button>
									<button
										onclick={() => abrirModalEliminar(emp)}
										class="rounded-lg p-2 text-text-muted transition-colors hover:bg-danger-bg hover:text-danger-color"
										title="Eliminar"
									>
										<Trash2 size={18} />
									</button>
								</div>
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>

<!-- Modal Formulario -->
{#if isModalOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-4 backdrop-blur-sm animate-modal-enter"
	>
		<div class="w-full max-w-lg rounded-2xl border border-border-color bg-bg-card shadow-2xl">
			<div class="flex items-center justify-between border-b border-border-color px-6 py-4">
				<h3 class="text-lg font-bold text-text-primary">
					{isEditing ? 'Editar Empleado' : 'Añadir Nuevo Empleado'}
				</h3>
				<button
					onclick={cerrarModal}
					class="rounded-lg p-1 text-text-muted hover:bg-border-color hover:text-text-primary transition-colors"
				>
					<X size={20} />
				</button>
			</div>

			<form onsubmit={guardarEmpleado} class="p-6">
				{#if formGeneralError}
					<div
						class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color"
						role="alert"
					>
						<span>{formGeneralError}</span>
					</div>
				{/if}

				<div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
					<!-- RUT -->
					<div class="flex flex-col gap-1.5">
						<label for="rut" class="text-sm font-semibold text-text-primary">RUT</label>
						<input
							id="rut"
							type="text"
							value={formData.rut}
							oninput={handleRutInput}
							disabled={submitLoading || isEditing}
							required
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50"
							placeholder="12.345.678-9"
						/>
						{#if errRut}<span class="text-xs font-medium text-danger-color">{errRut}</span>{/if}
					</div>

					<!-- Nombre -->
					<div class="flex flex-col gap-1.5">
						<label for="nombre" class="text-sm font-semibold text-text-primary"
							>Nombre Completo</label
						>
						<input
							id="nombre"
							type="text"
							bind:value={formData.nombre}
							onblur={handleNombreBlur}
							disabled={submitLoading}
							required
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50"
						/>
						{#if errNombre}<span class="text-xs font-medium text-danger-color">{errNombre}</span
							>{/if}
					</div>

					<!-- Usuario -->
					<div class="flex flex-col gap-1.5">
						<label for="usuario" class="text-sm font-semibold text-text-primary"
							>Alias de Usuario</label
						>
						<input
							id="usuario"
							type="text"
							bind:value={formData.usuario}
							disabled={submitLoading}
							required
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50"
						/>
						{#if errUsuario}<span class="text-xs font-medium text-danger-color">{errUsuario}</span
							>{/if}
					</div>

					<!-- Contraseña -->
					<div class="flex flex-col gap-1.5">
						<label for="contrasena" class="text-sm font-semibold text-text-primary">
							Contraseña {isEditing ? '(Dejar vacía para mantener)' : ''}
						</label>
						<div class="relative flex items-center">
							<input
								id="contrasena"
								type={showPassword ? 'text' : 'password'}
								bind:value={formData.contrasena}
								disabled={submitLoading}
								required={!isEditing}
								class="w-full rounded-xl border border-border-color bg-bg-primary py-2.5 pl-4 pr-11 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50"
								placeholder={isEditing ? '••••••••' : 'TuC0ntr4s3ña!'}
							/>
							<button
								type="button"
								onclick={() => (showPassword = !showPassword)}
								tabindex="-1"
								class="absolute right-3 text-text-muted hover:text-primario transition-colors focus:outline-none"
								title={showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
							>
								{#if showPassword}
									<EyeOff size={18} />
								{:else}
									<Eye size={18} />
								{/if}
							</button>
						</div>

						{#if errContrasena}<span class="text-xs font-medium text-danger-color"
								>{errContrasena}</span
							>{/if}
					</div>

					<!-- Teléfono -->
					<div class="flex flex-col gap-1.5">
						<label for="telefono" class="text-sm font-semibold text-text-primary">Teléfono</label>
						<input
							id="telefono"
							type="tel"
							bind:value={formData.telefono}
							disabled={submitLoading}
							required
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50"
							placeholder="+56 9 1234 5678"
						/>
						{#if errTelefono}<span class="text-xs font-medium text-danger-color">{errTelefono}</span
							>{/if}
					</div>

					<!-- Rol Dinámico -->
					<div class="flex flex-col gap-1.5">
						<label for="rol" class="text-sm font-semibold text-text-primary">Rol del Sistema</label>
						<select
							id="rol"
							bind:value={formData.id_rol}
							disabled={submitLoading}
							required
							class="rounded-xl border border-border-color bg-bg-primary pl-4 pr-10 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none cursor-pointer disabled:opacity-50"
						>
							{#if roles.length === 0}
								<option value="" disabled>Cargando roles...</option>
							{:else}
								{#each roles as r}
									<option value={r.id_rol || r.id || r.ID}>{r.nombre || r.Nombre}</option>
								{/each}
							{/if}
						</select>
					</div>

					<!-- Estado (Switch) -->
					<div class="col-span-1 sm:col-span-2 flex items-center gap-3 pt-2">
						<input
							id="activo"
							type="checkbox"
							bind:checked={formData.activo}
							disabled={submitLoading}
							class="h-4 w-4 accent-primario cursor-pointer disabled:opacity-50"
						/>
						<label
							for="activo"
							class="text-sm font-semibold text-text-primary cursor-pointer {submitLoading
								? 'opacity-50'
								: ''}"
						>
							Empleado Activo en el sistema
						</label>
					</div>
				</div>

				<div class="mt-8 flex justify-end gap-3">
					<button
						type="button"
						onclick={cerrarModal}
						disabled={submitLoading}
						class="rounded-xl px-5 py-2.5 text-sm font-semibold text-text-muted hover:bg-border-color hover:text-text-primary transition-colors disabled:opacity-50"
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
							{isEditing ? 'Guardar Cambios' : 'Crear Empleado'}
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

{#if showDeleteModal && empleadoToDelete}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm animate-modal-enter"
		onclick={() => (showDeleteModal = false)}
		role="presentation"
	>
		<div
			class="w-full max-w-[500px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg"
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
					¿Estás seguro de que deseas eliminar al empleado <strong>{empleadoToDelete.nombre}</strong
					> de forma permanente?
				</p>
				<p class="mt-3 text-xs text-text-muted">
					Esta acción no se puede deshacer y el empleado perderá el acceso al sistema de forma
					definitiva.
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
					onclick={confirmarEliminacion}>Eliminar Permanentemente</button
				>
			</footer>
		</div>
	</div>
{/if}
