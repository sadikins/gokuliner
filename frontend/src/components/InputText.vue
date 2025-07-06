<script setup>

    const props = defineProps( {
        inputType: {
            type: String,
        default: "text"
        }


    } )

    const modelValue = defineModel()

</script>

<template>
    <input
        :type="props.inputType"
        v-model="modelValue"
        class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg block w-full p-2.5"

    />
</template>