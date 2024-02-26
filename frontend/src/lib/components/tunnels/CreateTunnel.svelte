<script lang="ts">
	import { Button, Dialog, NumberInput, Input } from '$lib/components/ui';
	import { api, type TunnelCreate } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import { createEventDispatcher } from 'svelte';

	let tunnel: TunnelCreate = {
		name: '',
		localHost: '',
		localPort: 0,
		remotePort: 0
	};

	export let open = false;
	const eventDispatcher = createEventDispatcher<{ submit: null }>();
</script>

<Dialog title="Neuen Tunnel erstellen" bind:open>
	<div class="flex flex-col gap-2">
		<Input label="Name" placeholder="Plex" bind:value={tunnel.name} />
		<Input label="Lokaler Host" placeholder="localhost" bind:value={tunnel.localHost} />
		<NumberInput label="Lokaler Port" placeholder="32400" bind:value={tunnel.localPort} />
		<NumberInput label="Entfernter Port" placeholder="32400" bind:value={tunnel.remotePort} />
		<Button
			disabled={!tunnel.name ||
				!tunnel.localHost ||
				tunnel.localPort <= 0 ||
				tunnel.remotePort <= 0}
			on:click={() => {
				api.tunnels
					.create(tunnel)
					.then(() => {
						toast.success('Tunnel erstellt');
						open = false;
					})
					.catch(() => {
						toast.error('Tunnel konnte nicht erstellt werden');
					})
					.finally(() => {
						eventDispatcher('submit');
					});
			}}>Erstellen</Button
		>
	</div>
</Dialog>
