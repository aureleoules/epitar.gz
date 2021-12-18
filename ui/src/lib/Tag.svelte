<script lang="ts">
    export let tag: string;
</script>

<div class="tag">
    <slot/>
</div>

<style lang="scss">
    .tag {
        display: inline-block;
        background-color: #c0392b;
        padding: 8px;
        font-size: 16px;
        border-radius: 12px;
    }
</style>