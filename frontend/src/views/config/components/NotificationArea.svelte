<script lang="ts">
  export let loading = false;
  export let error = '';
  export let success = '';
</script>

<div class="notification-area">
  {#if loading}
    <div class="loading-indicator">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>
  {/if}

  {#if error}
    <div class="notification error">
      <div class="icon">✕</div>
      <p>{error}</p>
    </div>
  {/if}

  {#if success}
    <div class="notification success">
      <div class="icon">✓</div>
      <p>{success}</p>
    </div>
  {/if}
</div>

<style>
  .notification-area {
    margin-bottom: 16px;
  }
  
  .loading-indicator {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    background-color: #e3f2fd;
    border-radius: 6px;
    margin-bottom: 12px;
    animation: fadeIn 0.3s;
  }
  
  .spinner {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(74, 111, 165, 0.2);
    border-left-color: #4a6fa5;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  .notification {
    display: flex;
    align-items: center;
    padding: 10px 12px;
    border-radius: 6px;
    margin-bottom: 12px;
    animation: slideIn 0.3s;
  }
  
  .notification .icon {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    margin-right: 10px;
    font-weight: bold;
    font-size: 12px;
  }
  
  .notification.error {
    background-color: #ffebee;
    color: #c62828;
    border-left: 3px solid #c62828;
  }
  
  .notification.error .icon {
    background-color: #c62828;
    color: white;
  }
  
  .notification.success {
    background-color: #e8f5e9;
    color: #2e7d32;
    border-left: 3px solid #2e7d32;
  }
  
  .notification.success .icon {
    background-color: #2e7d32;
    color: white;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  @keyframes slideIn {
    from { transform: translateY(-10px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }
</style>