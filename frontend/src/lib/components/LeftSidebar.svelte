<script>
  export let user = null;
  
  const shortcuts = [
    { icon: '🧑‍🤝‍🧑', label: 'Friends', href: '/friends' },
    { icon: '📊', label: 'Meta AI', href: '/ai' },
    { icon: '📝', label: 'Memories', href: '/memories' },
    { icon: '📱', label: 'Saved', href: '/saved' },
    { icon: '👥', label: 'Groups', href: '/groups' },
    { icon: '📹', label: 'Reels', href: '/reels' },
    { icon: '🛍️', label: 'Marketplace', href: '/marketplace' }
  ];
  
  let showMore = false;
  
  function toggleShowMore() {
    showMore = !showMore;
  }
</script>

<aside class="left-sidebar">
  <div class="sidebar-content">
    <!-- User Profile -->
    <a href="/profile" class="sidebar-item user-item">
      <div class="item-icon profile-pic">
        {user?.name?.charAt(0) || 'U'}
      </div>
      <span class="item-label">{user?.name || 'User Name'}</span>
    </a>
    
    <!-- Main Navigation -->
    {#each shortcuts.slice(0, showMore ? shortcuts.length : 6) as shortcut}
      <a href={shortcut.href} class="sidebar-item">
        <div class="item-icon">
          {shortcut.icon}
        </div>
        <span class="item-label">{shortcut.label}</span>
      </a>
    {/each}
    
    <!-- Show More/Less Button -->
    <button class="sidebar-item show-more-btn" on:click={toggleShowMore}>
      <div class="item-icon">
        {showMore ? '🔼' : '🔽'}
      </div>
      <span class="item-label">{showMore ? 'See less' : 'See more'}</span>
    </button>
    
    <div class="sidebar-divider"></div>
    
    <!-- Your shortcuts section -->
    <div class="shortcuts-section">
      <h3 class="section-title">Your shortcuts</h3>
      <a href="/games/angry-birds" class="sidebar-item">
        <div class="item-icon game-icon">
          🐦
        </div>
        <span class="item-label">Angry Birds Friends</span>
      </a>
    </div>
  </div>
</aside>

<style>
  .left-sidebar {
    position: fixed;
    left: 0;
    top: 56px;
    width: 280px;
    height: calc(100vh - 56px);
    background: white;
    overflow-y: auto;
    padding: 16px 0;
    border-right: 1px solid #dadde1;
  }
  
  .sidebar-content {
    padding: 0 8px;
  }
  
  .sidebar-item {
    display: flex;
    align-items: center;
    padding: 8px;
    border-radius: 8px;
    text-decoration: none;
    color: #1c1e21;
    transition: background-color 0.2s;
    margin-bottom: 4px;
    cursor: pointer;
    border: none;
    background: none;
    width: 100%;
    text-align: left;
    font-size: 15px;
  }
  
  .sidebar-item:hover {
    background: #f0f2f5;
  }
  
  .user-item .item-icon {
    background: #1877f2;
    color: white;
  }
  
  .item-icon {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
    font-size: 18px;
    flex-shrink: 0;
  }
  
  .profile-pic {
    font-weight: bold;
  }
  
  .game-icon {
    background: #f0f2f5;
    border-radius: 8px;
  }
  
  .item-label {
    font-weight: 500;
    color: #1c1e21;
  }
  
  .show-more-btn .item-label {
    color: #65676b;
  }
  
  .sidebar-divider {
    height: 1px;
    background: #dadde1;
    margin: 16px 8px;
  }
  
  .shortcuts-section {
    padding: 0 8px;
  }
  
  .section-title {
    font-size: 17px;
    font-weight: 600;
    color: #65676b;
    margin: 16px 0 8px 0;
    padding-left: 8px;
  }
  
  /* Scrollbar styling */
  .left-sidebar::-webkit-scrollbar {
    width: 8px;
  }
  
  .left-sidebar::-webkit-scrollbar-track {
    background: transparent;
  }
  
  .left-sidebar::-webkit-scrollbar-thumb {
    background: #bcc0c4;
    border-radius: 4px;
  }
  
  .left-sidebar::-webkit-scrollbar-thumb:hover {
    background: #8a8d91;
  }
</style>