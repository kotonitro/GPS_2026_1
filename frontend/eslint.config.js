import prettier from 'eslint-config-prettier';
import path from 'node:path';
import js from '@eslint/js';
import svelte from 'eslint-plugin-svelte';
import { defineConfig, includeIgnoreFile } from 'eslint/config';
import globals from 'globals';
import ts from 'typescript-eslint';

const gitignorePath = path.resolve(import.meta.dirname, '.gitignore');

export default defineConfig(
	includeIgnoreFile(gitignorePath),
	js.configs.recommended,
	ts.configs.recommended,
	svelte.configs.recommended,
	prettier,
	svelte.configs.prettier,
	{
		languageOptions: { globals: { ...globals.browser, ...globals.node } },
		rules: {
			// typescript-eslint strongly recommend that you do not use the no-undef lint rule on TypeScript projects.
			// see: https://typescript-eslint.io/troubleshooting/faqs/eslint/#i-get-errors-from-the-no-undef-rule-about-global-variables-not-being-defined-even-though-there-are-no-typescript-errors
			'no-undef': 'off',
			// El codebase actual usa `any` extensivamente en helpers genéricos y payloads de API.
			// Deshabilitado para evitar ruido mientras se hace una migración progresiva a tipos estrictos.
			'@typescript-eslint/no-explicit-any': 'off',
			// Muchas variables e imports declarados están reservados para uso futuro o usados en plantillas.
			'@typescript-eslint/no-unused-vars': 'off',
			// Los bloques {#each} del proyecto no usan key en su mayoría.
			'svelte/require-each-key': 'off',
			// El proyecto usa goto() directamente en callbacks de auth; migrar a resolve() requiere refactor mayor.
			'svelte/no-navigation-without-resolve': 'off',
			// fix_focus.js y scripts similares usan require() por compatibilidad.
			'@typescript-eslint/no-require-imports': 'off',
			// Asignaciones intermedias como dvCalculado se usan para legibilidad antes de retornar.
			'no-useless-assignment': 'off',
			// Uso de Map/Date nativos en lugar de SvelteMap/SvelteDate; migración pendiente.
			'svelte/prefer-svelte-reactivity': 'off'
		}
	},
	{
		files: ['**/*.svelte', '**/*.svelte.ts', '**/*.svelte.js'],
		languageOptions: {
			parserOptions: {
				projectService: true,
				extraFileExtensions: ['.svelte'],
				parser: ts.parser
			}
		}
	},
	{
		// Override or add rule settings here, such as:
		// 'svelte/button-has-type': 'error'
		rules: {}
	}
);
