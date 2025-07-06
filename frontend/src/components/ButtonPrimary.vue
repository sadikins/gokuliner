
<script setup>

    const props = defineProps( {
    text: {type: String, default: 'Kirim'}
 })

</script>

<template>
  <button class="w-full text-white bg-primary hover:bg-blue-700 font-medium rounded-lg text-sm px-5 py-2.5">
    {{ props.text }}
</button>
</template>