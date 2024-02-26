<script lang="ts">
	export let value = '';
	export let placeholder = '';
	export let label = '';

	const id = Math.random().toString(36).substring(7);
</script>

<div class="flex w-full flex-col">
	{#if label}
		<label class="text-sm font-semibold" for={id}>{label}</label>
	{/if}
	<input
		{id}
		type="text"
		{placeholder}
		bind:value
		class="w-full rounded-md border px-2 py-1 placeholder:text-gray-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gray-900"
	/>
</div>
