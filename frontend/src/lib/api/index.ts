import { API } from './api';

export const BASE_URL = window.location.origin;
export const api = new API(BASE_URL);

export * from './types';
export * from './api';
