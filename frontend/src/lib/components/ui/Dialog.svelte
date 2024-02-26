<script lang="ts">
	import { fade } from 'svelte/transition';
	import X from 'virtual:icons/lucide/x';

	export let open = false;
	export let title = '';
</script>

{#if open}
	<dialog
		transition:fade={{ duration: 200 }}
		class="fixed inset-0 z-10 flex h-full w-full items-center justify-center bg-gray-50 bg-opacity-50 p-2 backdrop-blur-sm"
	>
		<div class="flex h-fit w-full max-w-xl flex-col rounded-md border bg-white">
			<div class="flex w-full items-center justify-between border-b p-2">
				<h2 class="ml-2 text-lg font-bold">{title}</h2>
				<button on:click={() => (open = false)}>
					<X class="h-6 w-6" />
				</button>
			</div>
			<div class="p-4">
				<slot />
			</div>
		</div>
	</dialog>
{/if}
