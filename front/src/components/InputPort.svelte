<script lang="ts">
  import { params } from '../stores/wol.ts';

  function handleInput(target: EventTarget & HTMLInputElement) {
    let next = +target.value;
    if (isNaN(next)) {
      next = 9;
    }

    if (next < 1) {
      next = 1;
    } else if (next > 65535) {
      next = 65535;
    }

    params.update((prev) => ({
      ...prev,
      port: next,
    }));

    target.value = next.toString();
  }
</script>

<input
  type="number"
  class="grow"
  min="1"
  max="65535"
  placeholder="9"
  on:input={(e) => handleInput(e.currentTarget)}
  />
