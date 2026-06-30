<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { apiFetch } from '$lib/api';

    // Sintaxis de Runes (Svelte 5) para recibir las páginas hijas
    let { children } = $props();
    let verificando = $state(true);

    onMount(async () => {
        try {
            // Intentamos llamar a cualquier ruta protegida de tu API (ej. empleados, cajas o un endpoint de verificación)
            // Si el usuario no está autenticado, el backend responderá con un error y saltará al catch
            await apiFetch('/empleados'); 
            
            // Si la petición fue exitosa, cambiamos el estado para dar acceso
            verificando = false;
        } catch (error) {
            // Si no hay sesión válida o la cookie expiró, lo expulsamos al login
            goto('/login');
        }
    });
</script>

{#if verificando}
    <div class="min-h-screen flex items-center justify-center">
        <p class="text-xl font-bold text-primario animate-pulse">
            Verificando credenciales...
        </p>
    </div>
{:else}
    {@render children()}
{/if}