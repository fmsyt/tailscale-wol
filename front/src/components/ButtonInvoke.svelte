<script lang="ts">
    import { wolHandler } from '../stores/wol.ts';

    export let text = 'Send Packet';
    export let variant = 'primary';
    export let className = '';

    let buttonClassName = `btn btn-${variant} ${className}`;
</script>

<button
    class={buttonClassName}
    disabled={!$wolHandler}
    on:click={() => $wolHandler?.()}
>
    {text}
</button>
