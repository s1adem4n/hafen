<script lang="ts">
	import Trash from 'virtual:icons/lucide/trash';
	import { Button, Dialog, NumberInput, Input } from '$lib/components/ui';
	import { api, type Tunnel } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import { createEventDispatcher } from 'svelte';

	export let tunnel: Tunnel;
	export let open = false;
	const eventDispatcher = createEventDispatcher<{ submit: null }>();
</script>

<Dialog title="Tunnel bearbeiten" bind:open>
	<div class="flex flex-col gap-2">
		<Input label="Name" placeholder="Plex" bind:value={tunnel.name} />
		<Input label="Lokaler Host" placeholder="localhost" bind:value={tunnel.localHost} />
		<NumberInput label="Lokaler Port" placeholder="32400" bind:value={tunnel.localPort} />
		<NumberInput label="Entfernter Port" placeholder="32400" bind:value={tunnel.remotePort} />
		<div class="flex w-full gap-2">
			<button
				class="flex items-center justify-center"
				on:click={() => {
					api.tunnels
						.delete(tunnel.id)
						.then((res) => {
							if (res) {
								toast.success('Tunnel gelöscht');
								open = false;
							} else {
								toast.error('Tunnel konnte nicht gelöscht werden');
							}
						})
						.catch(() => {
							toast.error('Tunnel konnte nicht gelöscht werden');
						})
						.finally(() => {
							eventDispatcher('submit');
						});
				}}
			>
				<Trash class="h-6 w-6 text-red-600" />
			</button>
			<Button
				disabled={!tunnel.name ||
					!tunnel.localHost ||
					tunnel.localPort <= 0 ||
					tunnel.remotePort <= 0}
				on:click={() => {
					api.tunnels
						.update(tunnel.id, tunnel)
						.then(() => {
							toast.success('Tunnel aktualisiert');
							open = false;
						})
						.catch(() => {
							toast.error('Tunnel konnte nicht aktualisiert werden');
						})
						.finally(() => {
							eventDispatcher('submit');
						});
				}}>Bearbeiten</Button
			>
		</div>
	</div>
</Dialog>
