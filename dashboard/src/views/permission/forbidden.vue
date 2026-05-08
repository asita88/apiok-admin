<template>
  <div class="perm-forbidden">
    <a-result status="403" title="403" sub-title="当前账号无权访问该页面。">
      <template #extra>
        <a-button type="primary" @click="$router.replace('/')">返回首页</a-button>
      </template>
    </a-result>
  </div>
</template>

<script>
export default {
  name: 'Forbidden'
}
</script>

<style scoped>
.perm-forbidden {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}
</style>
