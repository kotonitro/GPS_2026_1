<script lang="ts">
	import { page } from '$app/stores';

	const status = $page.status || 404;
	const mensaje =
		status === 404
			? 'Página no encontrada'
			: 'Error interno';
</script>

<div class="text-center mt-[20vh]">
	<h1 class="text-[8rem] m-0 text-primario font-bold">{status}</h1>

	<h2 class="text-xl mt-2">{mensaje}</h2>

	<br />

	<a href="/" class="text-primario underline hover:text-primario-hover">Volver al inicio</a>
</div>
