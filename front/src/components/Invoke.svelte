<script lang="ts">
    import { wolHandler } from '../stores/wol.ts';

    export let text = 'Send Packet';

</script>

<button
    class="btn btn-primary"
    on:click={() => $wolHandler?.()}
    disabled={!$wolHandler}
>
    {text}
</button>
