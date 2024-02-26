export interface TunnelCreate {
	name: string;
	remotePort: number;
	localHost: string;
	localPort: number;
}

export interface TunnelUpdate extends TunnelCreate {}

export interface Tunnel extends TunnelCreate {
	id: number;
	pid?: number;
}

export interface ProxyCreate {
	name: string;
	match: string;
	upstream: string;
}

export interface Proxy extends ProxyCreate {
	id: number;
}
