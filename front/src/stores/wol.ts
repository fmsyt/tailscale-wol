import { derived, writable, type Writable } from 'svelte/store';

const macAddressMatcher = /^([\dA-Fa-f]{2}[:-]){5}([\dA-Fa-f]{2})$/;
const ipMatcher = /^(\d{1,3}\.){3}\d{1,3}$/;

export interface WolParams {
  /** MAC Address */
  a: string;

  /** Port */
  p: number | null;

  /** Broadcast IP */
  b: string | null;

  /** Command */
  command: "wol" | "netcat";
}

export const params = writable<WolParams>({
  a: '',
  p: null,
  b: null,
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

export const isValidMacAddress = derived(params, $params => validateMacAddress($params.a));
export const isValidPort = derived(params, $params => !$params.p || validatePort($params.p));
export const isValidBroadcast = derived(params, $params => !$params.b || validateBroadcastIp($params.b));

export type WolHandler = (params?: WolParams) => Promise<string>;

export const wolHandler = derived<Writable<WolParams>, WolHandler | null>(params, ($params) => {

  if (!validateMacAddress($params.a)) {
    return null;
  }

  if ($params.p && !validatePort($params.p)) {
    return null;
  }

  if ($params.b && !validateBroadcastIp($params.b)) {
    return null;
  }

  return async () => {

    const query = new URLSearchParams();
    Object.entries($params || {}).forEach(([key, value]) => {
      value && query.append(key, value);
    });

    const response = await fetch(`/run?${query.toString()}`);
    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }

    const text = await response.text();
    return text;
  }
})
