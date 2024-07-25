<script>

  let address = '';
  let validMacAddress = false;

  function handleInput(e) {
    const value = e.target.value.replace(/[^\dA-Fa-f]/g, '').substring(0, 12).toUpperCase();

    const match = value.match(/.{1,2}/g);
    address = match ? match.join(':') : '';

    e.target.value = address;
  }

  async function handleSubmit() {

    if (!validMacAddress) {
      return;
    }

    const params = new URLSearchParams();
    params.append('a', address);

    fetch(`/run?${params.toString()}`);
  }

  $: {
    validMacAddress = address.length === 17;
  }
</script>

<div class="relative">
  <div class="input input-bordered text-xl flex items-center p-0 w-full">
    <input
      on:input={handleInput}
      type="text"
      class="grow text-center"
      placeholder="MAC Address"
      />
  </div>
  <!-- <button class="btn btn-ghost absolute top-0 right-0">Ghost</button> -->
</div>

<div class="card-actions justify-center">
  <button
    on:click={handleSubmit}
    class="btn btn-primary"
    disabled={!validMacAddress}
  >
    Send Packet
  </button>
</div>
