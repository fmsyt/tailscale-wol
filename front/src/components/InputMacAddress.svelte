<script lang="ts">

  import { params } from '../stores/wol.ts';

  function handleInput(target: EventTarget & HTMLInputElement) {
    const value = target.value.replace(/[^\dA-Fa-f]/g, '').substring(0, 12).toUpperCase();

    const match = value.match(/.{1,2}/g);
    const nextAddress = match ? match.join(':') : '';

    params.update((prev) => ({
      ...prev,
      macAddress: nextAddress,
    }));

    target.value = nextAddress;
  }

</script>

<input
  on:input={(e) => handleInput(e.currentTarget)}
  type="text"
  class="grow text-center"
  placeholder="MAC Address"
  />
