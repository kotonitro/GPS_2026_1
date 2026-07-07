<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Search,
		Plus,
		Edit2,
		Trash2,
		MonitorSmartphone, // Icono representativo para Cajas/Terminales
		CheckCircle2,
		XCircle
	} from '@lucide/svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { apiCajas } from '$lib/api';

	// Estados
	let cajas = $state<any[]>([]);
	let isLoading = $state(true);
	let searchQuery = $state('');
	let filtroEstado = $state('Todos'); // 'Todos', 'Activas', 'Inactivas'

	// Estados de Modales
	let showModal = $state(false);
	let showDeleteModal = $state(false);
	let isEditing = $state(false);
	let cajaToDelete = $state<any>(null);

	// Formulario
	let formData = $state({
		id_caja: '',
		nombre: '',
		ubicacion: '',
		saldo_inicial: 0,
		saldo_final: 0,
		activo: true
	});

	// Errores y Loading del form
	let errNombre = $state('');
	let errUbicacion = $state('');
	let formGeneralError = $state('');
	let submitLoading = $state(false);

	// Filtro reactivo
	let cajasFiltradas = $derived(
		cajas.filter((c) => {
			// 1. Filtro de Búsqueda de texto (Nombre o Ubicación)
			const coincideTexto =
				(c.nombre || c.Nombre || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
				(c.ubicacion || c.Ubicacion || '').toLowerCase().includes(searchQuery.toLowerCase());

			// 2. Filtro de Estado
			const esActiva = c.activo !== undefined ? c.activo : c.Activo;
			const coincideEstado =
				filtroEstado === 'Todos' ||
				(filtroEstado === 'Activas' && esActiva) ||
				(filtroEstado === 'Inactivas' && !esActiva);

			return coincideTexto && coincideEstado;
		})
	);

	onMount(async () => {
		await cargarDatos();
	});

	async function cargarDatos() {
		isLoading = true;
		try {
			const res = await apiCajas.getAll();
			cajas = Array.isArray(res) ? res : [];
		} catch (error) {
			toast.show('Error al cargar la lista de cajas registradoras.', 'error');
		} finally {
			isLoading = false;
		}
	}

	// Formateador de CLP
	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('es-CL', {
			style: 'currency',
			currency: 'CLP',
			maximumFractionDigits: 0
		}).format(amount);
	}

	function clearErrors() {
		errNombre = '';
		errUbicacion = '';
		formGeneralError = '';
	}

	// --- Controladores de Modales de Formulario ---

	function abrirModalNuevo() {
		isEditing = false;
		clearErrors();
		formData = {
			id_caja: '',
			nombre: '',
			ubicacion: '',
			saldo_inicial: 0,
			saldo_final: 0,
			activo: true
		};
		showModal = true;
	}

	function abrirModalEditar(caja: any) {
		isEditing = true;
		clearErrors();
		formData = {
			id_caja: caja.id_caja || caja.id || caja.ID,
			nombre: caja.nombre || caja.Nombre,
			ubicacion: caja.ubicacion || caja.Ubicacion,
			saldo_inicial: caja.saldo_inicial !== undefined ? caja.saldo_inicial : caja.SaldoInicial,
			saldo_final: caja.saldo_final !== undefined ? caja.saldo_final : caja.SaldoFinal,
			activo: caja.activo !== undefined ? caja.activo : caja.Activo
		};
		showModal = true;
	}

	async function guardarCaja(e: Event) {
		e.preventDefault();
		clearErrors();
		let isValid = true;

		if (!formData.nombre.trim()) {
			errNombre = 'El nombre de la caja es obligatorio.';
			isValid = false;
		}

		if (!formData.ubicacion.trim()) {
			errUbicacion = 'La ubicación de la caja es obligatoria.';
			isValid = false;
		}

		if (formData.saldo_inicial < 0 || formData.saldo_final < 0) {
			formGeneralError = 'Los saldos no pueden ser negativos.';
			isValid = false;
		}

		if (!isValid) return;

		submitLoading = true;
		try {
			// Asegurarse de que los valores numéricos se envíen correctamente
			const payload = {
				...formData,
				saldo_inicial: Number(formData.saldo_inicial),
				saldo_final: Number(formData.saldo_final)
			};

			if (isEditing) {
				await apiCajas.update(payload.id_caja, payload);
				toast.show('Caja actualizada correctamente.', 'success');
			} else {
				await apiCajas.create(payload);
				toast.show('Caja registrada exitosamente.', 'success');
			}

			showModal = false;
			await cargarDatos();
		} catch (error: any) {
			// Manejo de errores específicos del backend (por ejemplo, ubicación duplicada)
			if (error.detalle) {
				formGeneralError = error.detalle;
			} else if (error.errors && (error.errors.Ubicacion || error.errors.ubicacion)) {
				errUbicacion = error.errors.Ubicacion || error.errors.ubicacion;
			} else {
				formGeneralError = error.message || 'Ocurrió un error al guardar la caja.';
			}
			toast.show('No se pudo guardar el registro de la caja.', 'error');
		} finally {
			submitLoading = false;
		}
	}

	// --- Controladores del Modal de Eliminación ---

	function abrirModalEliminar(caja: any) {
		cajaToDelete = caja;
		showDeleteModal = true;
	}

	async function confirmarEliminacion() {
		if (!cajaToDelete) return;
		try {
			const id = cajaToDelete.id_caja || cajaToDelete.id || cajaToDelete.ID;
			await apiCajas.delete(id);
			toast.show('Caja eliminada con éxito.', 'success');
			showDeleteModal = false;
			cajaToDelete = null;
			await cargarDatos();
		} catch (error: any) {
			toast.show(
				error.message ||
					'No se pudo eliminar la caja. Comprueba si tiene registros de turnos asociados.',
				'error'
			);
		}
	}
</script>

<svelte:head>
	<title>Gestión de Cajas - MinimarketGo</title>
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
					placeholder="Buscar por nombre o ubicación..."
					bind:value={searchQuery}
					class="w-full rounded-xl border border-border-color bg-bg-card py-2.5 pl-10 pr-4 text-sm text-text-primary transition-colors focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario"
				/>
			</div>

			<!-- Filtro por Estado -->
			<select
				bind:value={filtroEstado}
				class="w-full sm:w-48 rounded-xl border border-border-color bg-bg-card py-2.5 pl-4 pr-10 text-sm text-text-primary transition-colors focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario cursor-pointer"
			>
				<option value="Todos">Todos los estados</option>
				<option value="Activas">Solo Activas</option>
				<option value="Inactivas">Solo Inactivas</option>
			</select>
		</div>

		<button
			onclick={abrirModalNuevo}
			class="flex w-full items-center justify-center gap-2 rounded-xl bg-primario px-5 py-2.5 text-sm font-semibold text-white transition-transform hover:scale-[1.02] active:scale-95 sm:w-auto"
		>
			<Plus size={18} />
			Nueva Caja
		</button>
	</div>

	<!-- Tabla -->
	<div class="overflow-x-auto rounded-xl border border-border-color bg-bg-card shadow-sm">
		<table class="w-full text-left text-sm text-text-primary border-collapse whitespace-nowrap">
			<thead class="border-b border-border-color bg-text-primary/4 text-text-muted">
				<tr>
					<th class="px-6 py-4 font-bold uppercase tracking-wider text-xs">Caja Registradora</th>
					<th class="px-6 py-4 font-bold uppercase tracking-wider text-xs">Ubicación</th>
					<th class="px-6 py-4 font-bold uppercase tracking-wider text-xs">Saldo Inicial</th>
					<th class="px-6 py-4 font-bold uppercase tracking-wider text-xs">Saldo Final</th>
					<th class="px-6 py-4 font-bold uppercase tracking-wider text-xs">Estado</th>
					<th class="px-6 py-4 text-right font-bold uppercase tracking-wider text-xs w-[120px]"
						>Acciones</th
					>
				</tr>
			</thead>
			<tbody class="divide-y divide-border-color">
				{#if isLoading}
					<tr>
						<td colspan="6" class="py-8 text-center text-text-muted">Cargando cajas...</td>
					</tr>
				{:else if cajasFiltradas.length === 0}
					<tr>
						<td colspan="6" class="py-16 text-center">
							<div class="flex flex-col items-center justify-center">
								<MonitorSmartphone size={48} class="text-text-muted opacity-30 mb-4" />
								<h3 class="mb-1 text-lg font-semibold text-text-primary">
									No se encontraron cajas
								</h3>
								<p class="text-sm text-text-secondary">
									Ajusta los filtros o registra una nueva caja en el sistema.
								</p>
							</div>
						</td>
					</tr>
				{:else}
					{#each cajasFiltradas as caja (caja.id_caja || caja.id || caja.ID)}
						{@const nombre = caja.nombre || caja.Nombre}
						{@const ubicacion = caja.ubicacion || caja.Ubicacion}
						{@const saldoInicial =
							caja.saldo_inicial !== undefined ? caja.saldo_inicial : caja.SaldoInicial}
						{@const saldoFinal =
							caja.saldo_final !== undefined ? caja.saldo_final : caja.SaldoFinal}
						{@const esActiva = caja.activo !== undefined ? caja.activo : caja.Activo}

						<tr class="transition-colors hover:bg-text-primary/[0.015]">
							<!-- Nombre -->
							<td class="px-6 py-4">
								<div class="flex items-center gap-3">
									<div
										class="flex h-9 w-9 items-center justify-center rounded-lg bg-primario/10 text-primario"
									>
										<MonitorSmartphone size={18} />
									</div>
									<span class="font-semibold">{nombre}</span>
								</div>
							</td>

							<!-- Ubicación -->
							<td class="px-6 py-4 text-text-muted">
								{ubicacion}
							</td>

							<!-- Saldo Inicial -->
							<td class="px-6 py-4 font-medium text-text-primary">
								{formatCurrency(saldoInicial)}
							</td>

							<!-- Saldo Final -->
							<td class="px-6 py-4 font-medium text-text-primary">
								{formatCurrency(saldoFinal)}
							</td>

							<!-- Estado -->
							<td class="px-6 py-4">
								{#if esActiva}
									<span
										class="inline-flex items-center gap-1.5 rounded-full bg-green-500/10 px-2.5 py-1 text-xs font-medium text-green-500"
									>
										<CheckCircle2 size={14} /> Activa
									</span>
								{:else}
									<span
										class="inline-flex items-center gap-1.5 rounded-full bg-red-500/10 px-2.5 py-1 text-xs font-medium text-red-500"
									>
										<XCircle size={14} /> Inactiva
									</span>
								{/if}
							</td>

							<!-- Acciones -->
							<td class="px-6 py-4 text-right">
								<div class="flex items-center justify-end gap-2">
									<button
										onclick={() => abrirModalEditar(caja)}
										class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-primario"
										title="Editar Caja"
									>
										<Edit2 size={16} />
									</button>
									<button
										onclick={() => abrirModalEliminar(caja)}
										class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-red-500/20 hover:bg-danger-bg hover:text-danger-color"
										title="Eliminar Caja"
									>
										<Trash2 size={16} />
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
					{isEditing ? 'Configurar Caja Registradora' : 'Registrar Nueva Caja'}
				</h2>
				<button
					onclick={() => (showModal = false)}
					class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary"
				>
					&times;
				</button>
			</header>

			<form onsubmit={guardarCaja} class="p-6">
				{#if formGeneralError}
					<div
						class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color"
						role="alert"
					>
						<span>{formGeneralError}</span>
					</div>
				{/if}

				<div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
					<!-- Nombre de la Caja -->
					<div class="flex flex-col gap-1.5 sm:col-span-2">
						<label for="nombre" class="text-sm font-semibold text-text-secondary"
							>Identificador / Nombre</label
						>
						<input
							id="nombre"
							type="text"
							bind:value={formData.nombre}
							disabled={submitLoading}
							required
							placeholder="Ej: Caja Principal 01"
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50"
						/>
						{#if errNombre}<span class="text-xs font-medium text-danger-color">{errNombre}</span
							>{/if}
					</div>

					<!-- Ubicación -->
					<div class="flex flex-col gap-1.5 sm:col-span-2">
						<label for="ubicacion" class="text-sm font-semibold text-text-secondary"
							>Ubicación Física</label
						>
						<input
							id="ubicacion"
							type="text"
							bind:value={formData.ubicacion}
							disabled={submitLoading}
							required
							placeholder="Ej: Entrada Sur - Pasillo 1"
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50"
						/>
						{#if errUbicacion}<span class="text-xs font-medium text-danger-color"
								>{errUbicacion}</span
							>{/if}
					</div>

					<!-- Saldo Inicial -->
					<div class="flex flex-col gap-1.5">
						<label for="saldo_inicial" class="text-sm font-semibold text-text-secondary"
							>Saldo Inicial base</label
						>
						<div class="relative flex items-center">
							<span class="absolute left-4 text-text-muted font-medium">$</span>
							<input
								id="saldo_inicial"
								type="number"
								min="0"
								step="1"
								bind:value={formData.saldo_inicial}
								disabled={submitLoading}
								required
								class="w-full rounded-xl border border-border-color bg-bg-primary py-2.5 pl-8 pr-4 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
							/>
						</div>
					</div>

					<!-- Saldo Final (Base/Apertura) -->
					<div class="flex flex-col gap-1.5">
						<label for="saldo_final" class="text-sm font-semibold text-text-secondary"
							>Saldo Final Actual</label
						>
						<div class="relative flex items-center">
							<span class="absolute left-4 text-text-muted font-medium">$</span>
							<input
								id="saldo_final"
								type="number"
								min="0"
								step="1"
								bind:value={formData.saldo_final}
								disabled={submitLoading}
								required
								class="w-full rounded-xl border border-border-color bg-bg-primary py-2.5 pl-8 pr-4 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
							/>
						</div>
					</div>

					<!-- Estado (Switch) -->
					<div class="col-span-1 sm:col-span-2 flex items-center gap-3 pt-2">
						<input
							id="activo"
							type="checkbox"
							bind:checked={formData.activo}
							disabled={submitLoading}
							class="h-5 w-5 accent-primario cursor-pointer disabled:opacity-50 shrink-0"
						/>
						<div class="flex flex-col">
							<label
								for="activo"
								class="text-sm font-bold text-text-primary cursor-pointer {submitLoading
									? 'opacity-50'
									: ''}"
							>
								Caja Operativa
							</label>
							<p class="text-xs text-text-muted mt-1">
								Permite que esta caja sea seleccionada para la apertura de turnos.
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
							{isEditing ? 'Guardar Cambios' : 'Registrar Caja'}
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Modal Confirmación de Eliminación -->
{#if showDeleteModal && cajaToDelete}
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
					¿Estás seguro de que deseas eliminar la caja <strong
						>{cajaToDelete.nombre || cajaToDelete.Nombre}</strong
					>
					ubicada en <strong>{cajaToDelete.ubicacion || cajaToDelete.Ubicacion}</strong>?
				</p>
				<p class="mt-3 text-sm text-text-muted">
					No podrás eliminar esta caja si existe un historial de turnos o registros vinculados a
					ella.
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
					onclick={confirmarEliminacion}>Eliminar Caja</button
				>
			</footer>
		</div>
	</div>
{/if}
