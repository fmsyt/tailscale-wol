import { derived, writable, type Writable } from 'svelte/store';

const macAddressMatcher = /^([\dA-Fa-f]{2}[:-]){5}([\dA-Fa-f]{2})$/;
const ipMatcher = /^(\d{1,3}\.){3}\d{1,3}$/;

export interface WolParams {
  macAddress: string;
  port: number | null;
  broadcast: string | null;
  command: "wol" | "netcat";
}

export const params = writable<WolParams>({
  macAddress: '',
  port: null,
  broadcast: null,
  command: 'wol'
});


export function validateMacAddress(value: string) {
  return macAddressMatcher.test(value);
}

export function validatePort(value: number) {
  return value >= 1 && value <= 65535;
}

export function validateBroadcastIp(value: string) {

  if (!value) {
    return true;
  }

  if (!ipMatcher.test(value)) {
    return false;
  }

  const parts = value.split('.');
  if (parts.length !== 4) {
    return false;
  }

  return parts.every(part => parseInt(part) >= 0 && parseInt(part) <= 255);
}

export const isValidMacAddress = derived(params, $params => validateMacAddress($params.macAddress));
export const isValidPort = derived(params, $params => !$params.port || validatePort($params.port));
export const isValidBroadcast = derived(params, $params => !$params.broadcast || validateBroadcastIp($params.broadcast));

export type WolHandler = (params?: WolParams) => Promise<void>;

export const wolHandler = derived<Writable<WolParams>, WolHandler | null>(params, ($params) => {

  if (!validateMacAddress($params.macAddress)) {
    return null;
  }

  if ($params.port && !validatePort($params.port)) {
    return null;
  }

  if ($params.broadcast && !validateBroadcastIp($params.broadcast)) {
    return null;
  }

  return async () => {

    const query = new URLSearchParams();
    query.append('a', $params.macAddress);
    Object.entries(params || {}).forEach(([key, value]) => {
      value && query.append(key, value);
    });

    const response = await fetch(`/run?${query.toString()}`);
    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }
  }
})
