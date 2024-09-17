<script lang="ts">
  import { wolHandler } from '../stores/wol.ts';

  export let text = 'Send Packet';
  export let variant = 'primary';
  export let className = '';

  let buttonClassName = `btn btn-${variant} ${className}`;
  let invokeResult: string|null = null;

  let disabled = true;

  async function handleClick() {
    if (invokeResult) {
      return;
    }

    if (!$wolHandler) {
      return;
    }

    invokeResult = null;

    try {
      invokeResult = await $wolHandler();
    } catch (e) {
      invokeResult = e instanceof Error ? e.message : 'An error occurred';
    }
  }

  $: disabled = !$wolHandler || Boolean(invokeResult);
</script>

<button
  class={buttonClassName}
  disabled={disabled}
  on:click={handleClick}
>
  {text}
</button>
