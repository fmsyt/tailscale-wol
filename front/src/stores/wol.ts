import { derived, writable, type Writable } from 'svelte/store';

const matcher = /^([\dA-Fa-f]{2}[:-]){5}([\dA-Fa-f]{2})$/;

export const macAddress = writable('');

export interface WolParams {
  command: "wol" | "netcat";
}

export type WolHandler = (params?: WolParams) => Promise<void>;

export const wolHandler = derived<Writable<string>, WolHandler | null>(macAddress, ($macAddress) => {

  if (!matcher.test($macAddress)) {
    return null;
  }

  return async (params?: WolParams) => {

    const query = new URLSearchParams();
    query.append('a', $macAddress);
    Object.entries(params || {}).forEach(([key, value]) => {
      query.append(key, value);
    });

    const response = await fetch(`/run?${query.toString()}`);
    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }
  }
})
