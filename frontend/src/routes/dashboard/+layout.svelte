<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { checkSession } from '$lib/api';

	let { children } = $props();
	let verificando = $state(true);

	onMount(async () => {
		try {
			const empleado = await checkSession();

			verificando = false;
		} catch (error) {
			goto('/login');
		}
	});
</script>

{#if verificando}
	<div class="min-h-screen flex items-center justify-center">
		<p class="text-xl font-bold text-primario animate-pulse">Verificando credenciales...</p>
	</div>
{:else}
	{@render children()}
{/if}
