import type { Tunnel, TunnelCreate, TunnelUpdate, Proxy, ProxyCreate } from './types';

export class API {
	baseUrl: string;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
	}

	tunnels = {
		list: async (): Promise<Tunnel[] | null> => {
			const response = await fetch(`${this.baseUrl}/tunnels`);
			if (!response.ok) {
				return null;
			}
			return response.json();
		},
		get: async (id: number): Promise<Tunnel | null> => {
			const response = await fetch(`${this.baseUrl}/tunnels/${id}`);
			if (!response.ok) {
				return null;
			}
			return response.json();
		},
		create: async (tunnel: TunnelCreate): Promise<Tunnel | null> => {
			const response = await fetch(`${this.baseUrl}/tunnels`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(tunnel)
			});
			if (!response.ok) {
				return null;
			}
			return response.json();
		},
		delete: async (id: number): Promise<boolean> => {
			const response = await fetch(`${this.baseUrl}/tunnels/${id}`, {
				method: 'DELETE'
			});
			if (!response.ok) {
				return false;
			}
			return true;
		},
		update: async (id: number, tunnel: TunnelUpdate): Promise<boolean> => {
			const response = await fetch(`${this.baseUrl}/tunnels/${id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(tunnel)
			});
			if (!response.ok) {
				return false;
			}
			return true;
		},
		start: async (id: number): Promise<boolean> => {
			const response = await fetch(`${this.baseUrl}/tunnels/${id}/start`, {
				method: 'POST'
			});
			if (!response.ok) {
				return false;
			}
			return true;
		},
		stop: async (id: number): Promise<boolean> => {
			const response = await fetch(`${this.baseUrl}/tunnels/${id}/stop`, {
				method: 'POST'
			});
			if (!response.ok) {
				return false;
			}
			return true;
		}
	};

	proxies = {
		list: async (): Promise<Proxy[] | null> => {
			const response = await fetch(`${this.baseUrl}/proxies`);
			if (!response.ok) {
				return null;
			}
			return response.json();
		},
		get: async (id: number): Promise<Proxy | null> => {
			const response = await fetch(`${this.baseUrl}/proxies/${id}`);
			if (!response.ok) {
				return null;
			}
			return response.json();
		},
		create: async (proxy: ProxyCreate): Promise<Proxy | null> => {
			const response = await fetch(`${this.baseUrl}/proxies`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(proxy)
			});
			if (!response.ok) {
				return null;
			}
			return response.json();
		},
		delete: async (id: number): Promise<boolean> => {
			const response = await fetch(`${this.baseUrl}/proxies/${id}`, {
				method: 'DELETE'
			});
			if (!response.ok) {
				return false;
			}
			return true;
		}
	};
}
