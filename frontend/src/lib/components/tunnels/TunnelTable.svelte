<script lang="ts">
	import Pen from 'virtual:icons/lucide/pen';
	import Play from 'virtual:icons/lucide/play';
	import Pause from 'virtual:icons/lucide/Pause';
	import Plus from 'virtual:icons/lucide/plus';
	import Refresh from 'virtual:icons/lucide/refresh-cw';
	import Waypoints from 'virtual:icons/lucide/waypoints';

	import { api, type Tunnel } from '$lib/api';
	import { ProxiesDialog } from '$lib/components/proxies';
	import UpdateTunnel from './UpdateTunnel.svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import CreateTunnel from './CreateTunnel.svelte';

	const getTunnels = async () => {
		api.tunnels.list().then((res) => {
			tunnels = res || [];
		});
	};

	let tunnels: Tunnel[] = [];
	let updatingTunnel: Tunnel;
	let updateOpen = false;
	let createOpen = false;
	let proxiesOpen = false;

	let refreshRotate = false;
	onMount(getTunnels);
</script>

<UpdateTunnel bind:tunnel={updatingTunnel} on:submit={getTunnels} bind:open={updateOpen} />
<CreateTunnel bind:open={createOpen} on:submit={getTunnels} />
<ProxiesDialog bind:open={proxiesOpen} />

<div class="flex w-full flex-col p-2">
	<div class="flex w-fit items-center gap-2 rounded-md border bg-white p-1 shadow-sm">
		<button
			class="flex rounded-md p-1 text-sm transition-colors hover:bg-gray-100"
			on:click={() => {
				createOpen = true;
			}}
		>
			<Plus class="h-5 w-5" />
			<span class="ml-2">Neuer Tunnel</span>
		</button>
		<span class="h-6 w-px border-r" />
		<button
			class="flex rounded-md p-1 text-sm transition-colors hover:bg-gray-100"
			on:click={() => {
				refreshRotate = true;
				setTimeout(() => {
					refreshRotate = false;
				}, 500);
				getTunnels();
			}}
		>
			<Refresh
				class="h-5 w-5 {refreshRotate ? 'rotate-180 transition-transform duration-500' : ''}"
			/>
			<span class="ml-2">Aktualisieren</span>
		</button>
		<span class="h-6 w-px border-r" />
		<button
			class="flex rounded-md p-1 text-sm transition-colors hover:bg-gray-100"
			on:click={() => (proxiesOpen = true)}
		>
			<Waypoints class="h-5 w-5" />
			<span class="ml-2">Proxies</span>
		</button>
	</div>
	<div class="w-full overflow-auto">
		<table class="w-full text-sm">
			<thead>
				<tr class="border-b border-gray-200">
					<th class="h-10 px-2 text-left">Name</th>
					<th class="h-10 px-2 text-left">Lokaler Host</th>
					<th class="h-10 px-2 text-left">Lokaler Port</th>
					<th class="h-10 px-2 text-left">Entfernter Port</th>
					<th class="h-10 px-2 text-left">Status</th>
					<th class="h-10 w-[1%] px-2 text-left">Aktionen</th>
				</tr>
			</thead>
			<tbody>
				{#each tunnels as tunnel (tunnel.id)}
					<tr class="border-b border-gray-100 transition-colors hover:bg-gray-100">
						<td class="p-2 align-middle">{tunnel.name}</td>
						<td class="p-2 align-middle">{tunnel.localHost}</td>
						<td class="p-2 align-middle">{tunnel.localPort}</td>
						<td class="p-2 align-middle">{tunnel.remotePort}</td>
						<td class="p-2 align-middle">
							{#if tunnel.pid}
								<span class="text-green-600">Aktiv</span>
							{:else}
								<span class="text-red-600">Inaktiv</span>
							{/if}
						</td>
						<td class="flex justify-around p-2 align-middle">
							<button
								class="flex items-center"
								on:click={() => {
									updatingTunnel = tunnel;
									updateOpen = true;
								}}
							>
								<Pen class="h-5 w-5" />
							</button>
							<button
								class="flex items-center"
								on:click={() => {
									if (tunnel.pid) {
										api.tunnels
											.stop(tunnel.id)
											.then(() => {
												getTunnels();
											})
											.catch(() => {
												toast.error('Tunnel konnte nicht gestoppt werden');
											});
									} else {
										api.tunnels
											.start(tunnel.id)
											.then(() => {
												getTunnels();
											})
											.catch(() => {
												toast.error('Tunnel konnte nicht gestartet werden');
											});
									}
								}}
							>
								{#if tunnel.pid}
									<Pause class="h-5 w-5 text-red-600" />
								{:else}
									<Play class="h-5 w-5 text-green-600" />
								{/if}
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
