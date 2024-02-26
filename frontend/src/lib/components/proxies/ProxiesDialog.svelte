<script lang="ts">
	import MoveRight from 'virtual:icons/lucide/move-right';
	import Plus from 'virtual:icons/lucide/plus';
	import Trash from 'virtual:icons/lucide/trash';
	import { api, type Proxy, type ProxyCreate } from '$lib/api';
	import { onMount } from 'svelte';
	import { Dialog } from '$lib/components/ui';
	import Input from '../ui/Input.svelte';

	const getProxies = async () => {
		const res = await api.proxies.list();
		proxies = res || [];
	};

	let proxyCreate: ProxyCreate = {
		name: '',
		match: '',
		upstream: ''
	};

	let proxies: Proxy[] = [];
	export let open = false;
	onMount(getProxies);
</script>

<Dialog bind:open title="Proxies">
	<div class="flex flex-col gap-2">
		<p class="-mt-2 text-sm text-gray-500">
			Proxies leiten Anfragen an den Upstream weiter. Match kann auch eine Domain sein.
		</p>
		<h2 class="w-full border-b border-gray-200 font-bold">Neuer Proxy</h2>
		<div class="flex gap-2">
			<Input placeholder="Name" bind:value={proxyCreate.name} />
			<Input placeholder="Match" bind:value={proxyCreate.match} />
			<Input placeholder="Upstream" bind:value={proxyCreate.upstream} />
			<button
				class="flex items-center justify-center"
				on:click={() => {
					api.proxies
						.create(proxyCreate)
						.then(() => {
							getProxies();
							proxyCreate = {
								name: '',
								match: '',
								upstream: ''
							};
						})
						.catch(() => {
							console.error('Proxy konnte nicht erstellt werden');
						});
				}}
			>
				<Plus class="h-6 w-6" />
			</button>
		</div>
		{#if proxies.length > 0}
			<h2 class="w-full border-b border-gray-200 font-bold">Proxies</h2>
		{/if}
		<div class="grid grid-cols-5 items-center gap-y-2">
			{#each proxies as proxy}
				<span class="underline">
					{proxy.name}:
				</span>
				<span class="justify-self-center">
					<a class="hover:underline" href="https://{proxy.match}">
						{proxy.match}
					</a>
				</span>
				<span class="justify-self-center">
					<MoveRight class="h-5 w-5" />
				</span>
				<span class="justify-self-center">
					{proxy.upstream}
				</span>

				<button
					class="mr-0.5 flex items-center justify-center justify-self-end"
					on:click={() => {
						api.proxies
							.delete(proxy.id)
							.then(() => {
								getProxies();
							})
							.catch(() => {
								console.error('Proxy konnte nicht gelöscht werden');
							});
					}}
				>
					<Trash class="h-5 w-5 text-red-600" />
				</button>
			{/each}
		</div>
	</div>
</Dialog>
