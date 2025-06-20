<!-- src/routes/+page.svelte -->
<script>
    let activeTab = 'discover';
    let searchQuery = '';
    
    // Mock data
    const posts = [
      {
        id: 1,
        author: 'The Tennessee Holler',
        handle: '@theholler.bsky.social',
        time: '1d',
        content: 'Another SPACE X explosion',
        image: '/api/placeholder/600/300',
        replies: 1100,
        reposts: 2600,
        likes: 10100,
        hasVideo: true
      },
      {
        id: 2,
        author: 'Megan (she/her)',
        handle: '@americanmegalo.bsky.social',
        time: '1d',
        content: 'With everything LA has gone through in the last few months (and couple of years), let\'s celebrate the hatching of:\n\nTEN CALIFORNIA CONDORS!!!!!!!\n\nCongrats to the team at the Los Angeles Zoo!!!',
        image: '/api/placeholder/400/400',
        replies: 89,
        reposts: 234,
        likes: 1200
      }
    ];
  
    const sidebarItems = [
      { icon: '🏠', label: 'Home', active: false },
      { icon: '🔍', label: 'Explore', active: false },
      { icon: '🔔', label: 'Notifications', active: false },
      { icon: '💬', label: 'Chat', active: false },
      { icon: '📋', label: 'Feeds', active: false },
      { icon: '📝', label: 'Lists', active: false },
      { icon: '👤', label: 'Profile', active: false },
      { icon: '⚙️', label: 'Settings', active: false }
    ];
  
    const trendingTopics = [
      'NBA Finals',
      'National Guard', 
      'Beyoncé and Miley',
      'Superman',
      'Juneteenth',
      'SCOTUS'
    ];
  
    function setActiveTab(tab) {
      activeTab = tab;
    }
  
    function formatNumber(num) {
      if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
      if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
      return num.toString();
    }
  </script>
  
  <div class="container">
    <!-- Left Sidebar -->
    <aside class="sidebar">
      <div class="logo">
        <div class="logo-icon">🦋</div>
      </div>
      
      <nav class="nav-menu">
        {#each sidebarItems as item}
          <a href="#" class="nav-item" class:active={item.active}>
            <span class="nav-icon">{item.icon}</span>
            <span class="nav-label">{item.label}</span>
          </a>
        {/each}
      </nav>
      
      <button class="new-post-btn">New Post</button>
      
      <div class="user-menu">
        <div class="user-avatar"></div>
      </div>
    </aside>
  
    <!-- Main Content -->
    <main class="main-content">
      <!-- Header Tabs -->
      <header class="header-tabs">
        <button 
          class="tab-btn" 
          class:active={activeTab === 'discover'}
          on:click={() => setActiveTab('discover')}
        >
          Discover
        </button>
        <button 
          class="tab-btn" 
          class:active={activeTab === 'following'}
          on:click={() => setActiveTab('following')}
        >
          Following
        </button>
        <button 
          class="tab-btn" 
          class:active={activeTab === 'video'}
          on:click={() => setActiveTab('video')}
        >
          Video
        </button>
      </header>
  
      <!-- Posts Feed -->
      <div class="feed">
        {#each posts as post}
          <article class="post">
            <div class="post-header">
              <div class="author-avatar"></div>
              <div class="author-info">
                <h3 class="author-name">{post.author}</h3>
                <span class="author-handle">{post.handle}</span>
                <span class="post-time">· {post.time}</span>
              </div>
            </div>
            
            <div class="post-content">
              <p>{post.content}</p>
              {#if post.image}
                <div class="post-media">
                  {#if post.hasVideo}
                    <div class="video-overlay">
                      <div class="play-button">▶</div>
                    </div>
                  {/if}
                  <img src={post.image} alt="Post media" />
                </div>
              {/if}
            </div>
            
            <div class="post-actions">
              <button class="action-btn">
                <span class="action-icon">💬</span>
                <span class="action-count">{formatNumber(post.replies)}</span>
              </button>
              <button class="action-btn">
                <span class="action-icon">🔄</span>
                <span class="action-count">{formatNumber(post.reposts)}</span>
              </button>
              <button class="action-btn">
                <span class="action-icon">❤️</span>
                <span class="action-count">{formatNumber(post.likes)}</span>
              </button>
              <button class="action-btn">
                <span class="action-icon">📤</span>
              </button>
            </div>
          </article>
        {/each}
      </div>
    </main>
  
    <!-- Right Sidebar -->
    <aside class="right-sidebar">
      <div class="search-box">
        <input 
          type="text" 
          placeholder="Search" 
          bind:value={searchQuery}
          class="search-input"
        />
      </div>
      
      <div class="sidebar-section">
        <div class="section-tabs">
          <button class="section-tab active">Discover</button>
          <button class="section-tab">Following</button>
          <button class="section-tab">Video</button>
          <button class="section-tab">More feeds</button>
        </div>
      </div>
      
      <div class="sidebar-section">
        <h3 class="section-title">📈 Trending</h3>
        <div class="trending-list">
          {#each trendingTopics as topic}
            <button class="trending-item">{topic}</button>
          {/each}
        </div>
      </div>
      
      <div class="sidebar-footer">
        <a href="#">Feedback</a>
        <a href="#">Privacy</a>
        <a href="#">Terms</a>
        <a href="#">Help</a>
      </div>
    </aside>
  </div>
  
  <style>
    .container {
      display: grid;
      grid-template-columns: 275px 1fr 350px;
      min-height: 100vh;
      max-width: 1200px;
      margin: 0 auto;
    }
  
    /* Left Sidebar */
    .sidebar {
      padding: 20px;
      border-right: 1px solid var(--border-color);
      position: sticky;
      top: 0;
      height: 100vh;
      display: flex;
      flex-direction: column;
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
    }
  
    .nav-item:hover {
      background-color: var(--bg-secondary);
    }
  
    .nav-item.active {
      font-weight: 700;
    }
  
    .nav-icon {
      font-size: 24px;
      width: 26px;
      text-align: center;
    }
  
    .new-post-btn {
      background: var(--primary-blue);
      color: white;
      padding: 16px 32px;
      font-size: 17px;
      margin: 20px 0;
      width: 100%;
    }
  
    .new-post-btn:hover {
      background: var(--primary-blue-hover);
    }
  
    .user-menu {
      margin-top: auto;
    }
  
    .user-avatar {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background: linear-gradient(135deg, #ff6b6b, #4ecdc4);
    }
  
    /* Main Content */
    .main-content {
      border-right: 1px solid var(--border-color);
      min-height: 100vh;
    }
  
    .header-tabs {
      display: flex;
      border-bottom: 1px solid var(--border-color);
      position: sticky;
      top: 0;
      background: var(--bg-primary);
      z-index: 10;
    }
  
    .tab-btn {
      flex: 1;
      padding: 16px;
      background: none;
      color: var(--text-secondary);
      font-size: 15px;
      font-weight: 500;
      border-radius: 0;
      position: relative;
    }
  
    .tab-btn:hover {
      background: var(--bg-secondary);
    }
  
    .tab-btn.active {
      color: var(--text-primary);
      font-weight: 700;
    }
  
    .tab-btn.active::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 50%;
      transform: translateX(-50%);
      width: 60px;
      height: 4px;
      background: var(--primary-blue);
      border-radius: 2px;
    }
  
    .feed {
      max-width: 100%;
    }
  
    .post {
      padding: 16px 20px;
      border-bottom: 1px solid var(--border-color);
      transition: background-color 0.2s ease;
    }
  
    .post:hover {
      background: rgba(0, 0, 0, 0.02);
    }
  
    .post-header {
      display: flex;
      gap: 12px;
      margin-bottom: 12px;
    }
  
    .author-avatar {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background: linear-gradient(135deg, #667eea, #764ba2);
      flex-shrink: 0;
    }
  
    .author-info {
      display: flex;
      align-items: center;
      gap: 8px;
      flex-wrap: wrap;
    }
  
    .author-name {
      font-weight: 700;
      font-size: 15px;
      margin: 0;
      color: var(--text-primary);
    }
  
    .author-handle, .post-time {
      color: var(--text-secondary);
      font-size: 15px;
    }
  
    .post-content p {
      margin: 0 0 12px 0;
      font-size: 15px;
      line-height: 1.5;
      white-space: pre-line;
    }
  
    .post-media {
      position: relative;
      border-radius: 16px;
      overflow: hidden;
      margin-top: 12px;
    }
  
    .post-media img {
      width: 100%;
      height: auto;
      display: block;
    }
  
    .video-overlay {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      z-index: 2;
    }
  
    .play-button {
      width: 60px;
      height: 60px;
      border-radius: 50%;
      background: rgba(0, 0, 0, 0.8);
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 20px;
      cursor: pointer;
    }
  
    .post-actions {
      display: flex;
      gap: 60px;
      margin-top: 12px;
      padding-left: 52px;
    }
  
    .action-btn {
      display: flex;
      align-items: center;
      gap: 8px;
      background: none;
      color: var(--text-secondary);
      font-size: 13px;
      padding: 8px;
      border-radius: 20px;
      transition: all 0.2s ease;
    }
  
    .action-btn:hover {
      background: rgba(29, 155, 240, 0.1);
      color: var(--primary-blue);
    }
  
    .action-icon {
      font-size: 18px;
    }
  
    /* Right Sidebar */
    .right-sidebar {
      padding: 20px 16px;
      position: sticky;
      top: 0;
      height: 100vh;
      overflow-y: auto;
    }
  
    .search-box {
      margin-bottom: 20px;
    }
  
    .search-input {
      width: 100%;
      background: var(--bg-secondary);
      border: 1px solid transparent;
      border-radius: 24px;
      padding: 12px 16px;
      font-size: 15px;
    }
  
    .search-input:focus {
      background: var(--bg-primary);
      border-color: var(--primary-blue);
    }
  
    .sidebar-section {
      background: var(--bg-secondary);
      border-radius: 16px;
      margin-bottom: 20px;
      overflow: hidden;
    }
  
    .section-tabs {
      display: flex;
      flex-wrap: wrap;
      padding: 12px;
      gap: 8px;
    }
  
    .section-tab {
      background: none;
      border: 1px solid var(--border-color);
      color: var(--text-secondary);
      padding: 8px 16px;
      border-radius: 20px;
      font-size: 13px;
      transition: all 0.2s ease;
    }
  
    .section-tab:hover {
      background: var(--bg-primary);
    }
  
    .section-tab.active {
      background: var(--primary-blue);
      color: white;
      border-color: var(--primary-blue);
    }
  
    .section-title {
      padding: 16px 16px 0 16px;
      margin: 0;
      font-size: 20px;
      font-weight: 800;
      color: var(--text-primary);
    }
  
    .trending-list {
      padding: 4px 16px 16px 16px;
      display: flex;
      flex-direction: column;
      gap: 4px;
    }
  
    .trending-item {
      background: none;
      border: none;
      color: var(--text-primary);
      padding: 8px 0;
      text-align: left;
      font-size: 15px;
      border-radius: 8px;
      transition: background-color 0.2s ease;
    }
  
    .trending-item:hover {
      background: rgba(0, 0, 0, 0.05);
    }
  
    .sidebar-footer {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      margin-top: auto;
      padding-top: 20px;
    }
  
    .sidebar-footer a {
      color: var(--text-secondary);
      text-decoration: none;
      font-size: 13px;
    }
  
    .sidebar-footer a:hover {
      text-decoration: underline;
    }
  
    /* Responsive Design */
    @media (max-width: 1024px) {
      .container {
        grid-template-columns: 68px 1fr 300px;
      }
      
      .nav-label {
        display: none;
      }
      
      .new-post-btn {
        padding: 12px;
        font-size: 24px;
      }
    }
  
    @media (max-width: 768px) {
      .container {
        grid-template-columns: 1fr;
      }
      
      .sidebar,
      .right-sidebar {
        display: none;
      }
      
      .main-content {
        border: none;
      }
    }
  </style>