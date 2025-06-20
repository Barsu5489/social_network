<!-- src/lib/components/Navbar.svelte -->
<script>
    import { createEventDispatcher } from 'svelte';
    
    const dispatch = createEventDispatcher();
    
    const sidebarItems = [
      { icon: '🏠', label: 'Home', active: false, path: '/' },
      { icon: '🔍', label: 'Explore', active: false, path: '/explore' },
      { icon: '🔔', label: 'Notifications', active: false, path: '/notifications' },
      { icon: '💬', label: 'Chat', active: false, path: '/chat' },
      { icon: '📋', label: 'Feeds', active: false, path: '/feeds' },
      { icon: '📝', label: 'Lists', active: false, path: '/lists' },
      { icon: '👤', label: 'Profile', active: false, path: '/profile' },
      { icon: '⚙️', label: 'Settings', active: false, path: '/settings' }
    ];
  
    function handleNewPost() {
      dispatch('newPost');
    }
  
    function handleNavClick(item) {
      dispatch('navigate', { path: item.path, label: item.label });
    }
  </script>
  
  <aside class="sidebar">
    <div class="logo">
      <div class="logo-icon">🦋</div>
    </div>
    
    <nav class="nav-menu">
      {#each sidebarItems as item}
        <button 
          class="nav-item" 
          class:active={item.active}
          on:click={() => handleNavClick(item)}
        >
          <span class="nav-icon">{item.icon}</span>
          <span class="nav-label">{item.label}</span>
        </button>
      {/each}
    </nav>
    
    <button class="new-post-btn" on:click={handleNewPost}>
      New Post
    </button>
    
    <div class="user-menu">
      <div class="user-avatar"></div>
    </div>
  </aside>
  
  <style>
    .sidebar {
      padding: 20px;
      border-right: 1px solid var(--border-color);
      position: sticky;
      top: 0;
      height: 100vh;
      display: flex;
      flex-direction: column;
      background: var(--bg-primary);
    }
  
    .logo {
      margin-bottom: 30px;
    }
  
    .logo-icon {
      font-size: 32px;
      width: 50px;
      height: 50px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 50%;
      background: linear-gradient(135deg, var(--primary-blue), #00d4ff);
      cursor: pointer;
      transition: transform 0.2s ease;
    }
  
    .logo-icon:hover {
      transform: scale(1.05);
    }
  
    .nav-menu {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 2px;
    }
  
    .nav-item {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 12px 16px;
      border-radius: 24px;
      text-decoration: none;
      color: var(--text-primary);
      font-size: 20px;
      font-weight: 400;
      transition: background-color 0.2s ease;
      background: none;
      border: none;
      cursor: pointer;
      text-align: left;
      width: 100%;
    }
  
    .nav-item:hover {
      background-color: var(--bg-secondary);
    }
  
    .nav-item.active {
      font-weight: 700;
      background-color: var(--bg-secondary);
    }
  
    .nav-icon {
      font-size: 24px;
      width: 26px;
      text-align: center;
      flex-shrink: 0;
    }
  
    .nav-label {
      flex: 1;
    }
  
    .new-post-btn {
      background: var(--primary-blue);
      color: white;
      padding: 16px 32px;
      font-size: 17px;
      margin: 20px 0;
      width: 100%;
      border-radius: 24px;
      border: none;
      cursor: pointer;
      font-weight: 600;
      transition: background-color 0.2s ease;
    }
  
    .new-post-btn:hover {
      background: var(--primary-blue-hover);
    }
  
    .user-menu {
      margin-top: auto;
      padding: 12px 0;
    }
  
    .user-avatar {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background: linear-gradient(135deg, #ff6b6b, #4ecdc4);
      cursor: pointer;
      transition: transform 0.2s ease;
    }
  
    .user-avatar:hover {
      transform: scale(1.05);
    }
  
    /* Responsive Design */
    @media (max-width: 1024px) {
      .nav-label {
        display: none;
      }
      
      .new-post-btn {
        padding: 12px;
        font-size: 24px;
      }
      
      .sidebar {
        width: 68px;
        align-items: center;
      }
      
      .nav-item {
        justify-content: center;
        width: 50px;
        height: 50px;
        padding: 0;
      }
    }
  
    @media (max-width: 768px) {
      .sidebar {
        display: none;
      }
    }
  </style>