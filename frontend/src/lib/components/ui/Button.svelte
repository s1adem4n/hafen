<script lang="ts">
	export let disabled = false;
</script>

<button
	class="flex w-full items-center justify-center rounded-md bg-gray-900 px-2 py-1 text-gray-200 transition-opacity hover:opacity-90 disabled:opacity-75"
	on:click
	{disabled}
>
	<slot />
</button>
