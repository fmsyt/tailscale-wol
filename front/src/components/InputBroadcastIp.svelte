<script lang="ts">
  import { params } from '../stores/wol.ts';

  function handleInput(target: EventTarget & HTMLInputElement) {
    const segments = target.value
      .replace(/[^\d.]/g, '')
      .split('.')
      ;

    if (segments.length === 0) {
      params.update((prev) => ({ ...prev, b: "" }))
      return;
    }

    const lastIndex = segments.length - 1;
    const lastSegmentValue = Number(segments[lastIndex]);
    const digits = Math.floor(Math.log10(lastSegmentValue)) + 1;

    let lastSegment = Number(segments[lastIndex]);
    if (segments.length >= 4 ) {
      if (lastSegment > 255) {
        segments[3] = '255';
      }

      // remove segments after 4
      segments.splice(4);

    } else {
      if (digits === 3) {
        if (lastSegmentValue > 255) {
          segments[lastIndex] = '255';
        }
      } else if (digits > 3) {
        const divideBy = Math.pow(10, digits - 3);
        const replaceValue = Math.floor(lastSegmentValue / divideBy);
        segments[lastIndex] = replaceValue.toString();

        const carry = lastSegmentValue % divideBy;
        segments.push(carry.toString());
      }
    }

    const next = segments.join('.');
    params.update((prev) => ({ ...prev, b: next }));

    target.value = next;
  }
</script>

<input
  type="text"
  class="grow"
  placeholder="255.255.255.255"
  on:input={(e) => handleInput(e.currentTarget)}
  />
