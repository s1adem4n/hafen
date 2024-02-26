import { API } from './api';

export const BASE_URL = 'http://192.168.1.100:8080';
export const api = new API(BASE_URL);

export * from './types';
export * from './api';
