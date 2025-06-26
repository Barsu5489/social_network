<script>
  export let user = null;
  
  let searchQuery = '';
  
  function handleSearch() {
    console.log('Searching for:', searchQuery);
  }
  
  function handleNotifications() {
    console.log('Opening notifications');
  }
  
  function handleMessages() {
    console.log('Opening messages');
  }
  
  function handleProfile() {
    console.log('Opening profile');
  }
</script>

<nav class="navbar">
  <div class="navbar-content">
    <!-- Left section -->
    <div class="navbar-left">
      <div class="logo">
        <span class="logo-icon">cs</span>
      </div>
      <div class="search-container">
        <input 
          type="text" 
          placeholder="Search ConnectSphere" 
          bind:value={searchQuery}
          on:keydown={(e) => e.key === 'Enter' && handleSearch()}
          class="search-input"
        />
        <button class="search-btn" on:click={handleSearch}>
          🔍
        </button>
      </div>
    </div>
    
    <!-- Center section -->
    <div class="navbar-center">
      <button class="nav-tab active">
        🏠
      </button>
      <button class="nav-tab">
        👥
      </button>
      <button class="nav-tab">
        📺
      </button>
      <button class="nav-tab">
        🏪
      </button>
      <button class="nav-tab">
        👫
      </button>
    </div>
    
    <!-- Right section -->
    <div class="navbar-right">
      <button class="nav-icon-btn">
        ⊞
      </button>
      <button class="nav-icon-btn" on:click={handleMessages}>
        💬
      </button>
      <button class="nav-icon-btn notification-btn" on:click={handleNotifications}>
        🔔
        <span class="notification-badge">2</span>
      </button>
      <button class="profile-btn" on:click={handleProfile}>
        <div class="profile-avatar">
          {user?.name?.charAt(0) || 'U'}
        </div>
      </button>
    </div>
  </div>
</nav>

<style>
  .navbar {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    background: white;
    border-bottom: 1px solid #dadde1;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    z-index: 1000;
    height: 56px;
  }
  
  .navbar-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 16px;
    height: 100%;
  }
  
  .navbar-left {
    display: flex;
    align-items: center;
    gap: 10px;
    flex: 1;
  }
  
  .logo {
    width: 40px;
    height: 40px;
    background: #1877f2;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }
  
  .logo-icon {
    color: white;
    font-weight: bold;
    font-size: 24px;
  }
  
  .search-container {
    position: relative;
    margin-left: 10px;
  }
  
  .search-input {
    width: 240px;
    height: 40px;
    border: none;
    border-radius: 20px;
    background: #f0f2f5;
    padding: 0 16px 0 40px;
    font-size: 15px;
    outline: none;
  }
  
  .search-input:focus {
    background: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .search-btn {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    cursor: pointer;
    font-size: 16px;
  }
  
  .navbar-center {
    display: flex;
    gap: 8px;
    flex: 0 0 auto;
  }
  
  .nav-tab {
    width: 112px;
    height: 56px;
    background: none;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 20px;
    position: relative;
    transition: background-color 0.2s;
  }
  
  .nav-tab:hover {
    background: #f0f2f5;
  }
  
  .nav-tab.active {
    color: #1877f2;
  }
  
  .nav-tab.active::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: #1877f2;
    border-radius: 3px 3px 0 0;
  }
  
  .navbar-right {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    justify-content: flex-end;
  }
  
  .nav-icon-btn {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: #e4e6ea;
    border: none;
    cursor: pointer;
    font-size: 16px;
    position: relative;
    transition: background-color 0.2s;
  }
  
  .nav-icon-btn:hover {
    background: #d8dadf;
  }
  
  .notification-btn {
    position: relative;
  }
  
  .notification-badge {
    position: absolute;
    top: -2px;
    right: -2px;
    background: #ff3040;
    color: white;
    border-radius: 10px;
    width: 20px;
    height: 20px;
    font-size: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
  }
  
  .profile-btn {
    background: none;
    border: none;
    cursor: pointer;
  }
  
  .profile-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: #1877f2;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 14px;
  }
</style>