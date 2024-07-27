<script lang="ts">
  import { params } from '../stores/wol.ts';

  function handleInput(target: EventTarget & HTMLInputElement) {
    const value = target.value.replace(/[^\d]/g, '').substring(0, 12);

    const match = value.match(/.{1,3}/g);
    const next = match ? match.join('.') : '';

    params.update((prev) => ({
      ...prev,
      broadcast: next,
    }));

    target.value = next;
  }
</script>

<input
  type="text"
  class="grow"
  placeholder="255.255.255.255"
  on:input={(e) => handleInput(e.currentTarget)}
  />
